// Copyright 2009 Bart de Boer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cobrahooks

import (
	"github.com/spf13/cobra"
)

type Command struct {
	*cobra.Command
}

func NewCommand(c *cobra.Command) *Command {
	return &Command{c}
}

type commandHook struct {
	cmd  *cobra.Command
	hook func(cmd *cobra.Command, args []string) error
}

var (
	preRunHooks            []*commandHook
	persistentPreRunHooks  []*commandHook
	runHooks               []*commandHook
	persistentPostRunHooks []*commandHook
	postRunHooks           []*commandHook
	helpHooks              []*commandHook
)

func (c *Command) OnRun(h func(cmd *cobra.Command, args []string) error) {
	OnRun(c.Command, h)
}

// OnRun registers a Run hook onto the command.
func OnRun(c *cobra.Command, h func(cmd *cobra.Command, args []string) error) {
	// Register the hook
	runHooks = append(runHooks, &commandHook{
		cmd:  c,
		hook: h,
	})
	if c.RunE != nil {
		return
	}
	c.RunE = func(cmd *cobra.Command, args []string) error {
		// find and execute any registered Run hooks
		for _, ch := range runHooks {
			if ch.cmd == cmd {
				if err := ch.hook(cmd, args); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

func (c *Command) OnPreRun(h func(cmd *cobra.Command, args []string) error) {
	OnPreRun(c.Command, h)
}

// OnPreRun registers a PreRun hook on the command.
func OnPreRun(c *cobra.Command, h func(cmd *cobra.Command, args []string) error) {
	// Register the hook
	preRunHooks = append(preRunHooks, &commandHook{
		cmd:  c,
		hook: h,
	})
	if c.PreRunE != nil {
		return
	}
	c.PreRunE = func(cmd *cobra.Command, args []string) error {
		// find and execute any registered PreRun hooks
		for _, ch := range preRunHooks {
			if ch.cmd == cmd {
				if err := ch.hook(cmd, args); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

func (c *Command) OnPostRun(h func(cmd *cobra.Command, args []string) error) {
	OnPostRun(c.Command, h)
}

// OnPostRun registers a PostRun hook on the command.
func OnPostRun(c *cobra.Command, h func(cmd *cobra.Command, args []string) error) {
	// Register the hook
	postRunHooks = append(postRunHooks, &commandHook{
		cmd:  c,
		hook: h,
	})
	if c.PostRunE != nil {
		return
	}
	c.PostRunE = func(cmd *cobra.Command, args []string) error {
		// find and execute any registered PreRun hooks
		for _, ch := range postRunHooks {
			if ch.cmd == cmd {
				if err := ch.hook(cmd, args); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

func (c *Command) OnPersistentPreRun(h func(cmd *cobra.Command, args []string) error) {
	OnPersistentPreRun(c.Command, h)
}

// OnPersistentPostRun registers a PreRun hook on the command and all of its childs
func OnPersistentPreRun(c *cobra.Command, h func(cmd *cobra.Command, args []string) error) {
	// Register the hook
	// Prepend to the array to ensure hooks are executed
	// in the same order as how they are registered for the command
	// (since the runChain will be executed in reversed order from parent to child)
	persistentPreRunHooks = append([]*commandHook{
		&commandHook{
			cmd:  c,
			hook: h,
		},
	}, persistentPreRunHooks...)

	if c.PersistentPreRunE != nil {
		return
	}
	c.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		var runChain []*commandHook
		// Walk up the command chain
		for p := cmd; p != nil; p = p.Parent() {
			// find any registered PersistentPreRun hooks and build the run chain
			for _, ch := range persistentPreRunHooks {
				if ch.cmd == p {
					runChain = append(runChain, &commandHook{
						hook: ch.hook,
					})
				}
			}
		}
		// Run the command chain hooks from parent to child
		for i := len(runChain) - 1; i >= 0; i-- {
			if err := runChain[i].hook(cmd, args); err != nil {
				return err
			}
		}
		return nil
	}
}

func (c *Command) OnPersistentPostRun(h func(cmd *cobra.Command, args []string) error) {
	OnPersistentPostRun(c.Command, h)
}

// OnPersistentPostRun registers a PostRun hook on the command and all of its childs
func OnPersistentPostRun(c *cobra.Command, h func(cmd *cobra.Command, args []string) error) {
	// Register the hook
	persistentPostRunHooks = append(persistentPostRunHooks, &commandHook{
		cmd:  c,
		hook: h,
	})
	if c.PersistentPostRunE != nil {
		return
	}
	c.PersistentPostRunE = func(cmd *cobra.Command, args []string) error {
		// Walk up the command chain
		for p := cmd; p != nil; p = p.Parent() {
			// find and execute any registered PostRun hooks
			for _, ch := range persistentPostRunHooks {
				if ch.cmd == p {
					if err := ch.hook(cmd, args); err != nil {
						return err
					}
				}
			}
		}
		return nil
	}
}

func (c *Command) OnHelp(h func(cmd *cobra.Command, args []string) error) {
	OnHelp(c.Command, h)
}

// OnHelp registers a hook for when help is invoked.
func OnHelp(c *cobra.Command, h func(cmd *cobra.Command, args []string) error) {
	// Register the hook
	helpHooks = append(helpHooks, &commandHook{
		cmd:  c,
		hook: h,
	})
	// Search if a help hook is already registered for the command.
	// We don't need to integrate again.
	for _, ch := range helpHooks {
		if c == ch.cmd {
			return
		}
	}
	// Integrate the help hooks for this command
	helpFunc := c.HelpFunc()
	c.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		for _, ch := range helpHooks {
			if ch.cmd == cmd {
				if err := ch.hook(cmd, args); err != nil {
					return
				}
			}
		}
		helpFunc(cmd, args)
	})
}
