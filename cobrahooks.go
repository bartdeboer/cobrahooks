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
	cmd        *cobra.Command
	hook       func(cmd *cobra.Command, args []string) error
	runOnHelp  bool
	persistent bool
}

var (
	preRunHooks            []*commandHook
	persistentPreRunHooks  []*commandHook
	runHooks               []*commandHook
	persistentPostRunHooks []*commandHook
	postRunHooks           []*commandHook
	helpHooks              []*commandHook
)

type HookOptions struct {
	runOnHelp  bool
	persistent bool
}

func RunOnHelp(o *HookOptions) { o.runOnHelp = true }

func Persistent(o *HookOptions) { o.persistent = true }

// OnRun registers a Run hook onto the command.
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

func runPreRunHooks(cmd *cobra.Command, args []string, isHelpRun bool) error {
	for _, ch := range preRunHooks {
		if ch.cmd == cmd && (!isHelpRun || ch.runOnHelp) {
			if err := ch.hook(cmd, args); err != nil {
				return err
			}
		}
	}
	return nil
}

// OnPreRun registers a PreRun hook on the command.
func (c *Command) OnPreRun(h func(cmd *cobra.Command, args []string) error, options ...func(*HookOptions)) {
	OnPreRun(c.Command, h, options...)
}

// OnPreRun registers a PreRun hook on the command.
func OnPreRun(c *cobra.Command, h func(cmd *cobra.Command, args []string) error, options ...func(*HookOptions)) {
	var opts HookOptions
	for _, option := range options {
		option(&opts)
	}
	if opts.persistent {
		OnPersistentPreRun(c, h, options...)
		return
	}
	// Register the hook
	preRunHooks = append(preRunHooks, &commandHook{
		cmd:  c,
		hook: h,
	})
	if opts.runOnHelp {
		initHelpHooks(c)
	}
	if c.PreRunE != nil {
		return
	}
	c.PreRunE = func(cmd *cobra.Command, args []string) error {
		return runPreRunHooks(cmd, args, false)
	}
}

// OnPostRun registers a PostRun hook on the command.
func (c *Command) OnPostRun(h func(cmd *cobra.Command, args []string) error, options ...func(*HookOptions)) {
	OnPostRun(c.Command, h, options...)
}

// OnPostRun registers a PostRun hook on the command.
func OnPostRun(c *cobra.Command, h func(cmd *cobra.Command, args []string) error, options ...func(*HookOptions)) {
	var opts HookOptions
	for _, option := range options {
		option(&opts)
	}
	if opts.persistent {
		OnPersistentPostRun(c, h)
		return
	}
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

func runPersistentPreRunHooks(cmd *cobra.Command, args []string, isHelpRun bool) error {
	var runChain []*commandHook
	// Walk up the command chain
	for p := cmd; p != nil; p = p.Parent() {
		// find any registered PersistentPreRun hooks and build the run chain
		for _, ch := range persistentPreRunHooks {
			if ch.cmd == p && (!isHelpRun || ch.runOnHelp) {
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

// OnPersistentPostRun registers a PreRun hook on the command and all of its childs
func (c *Command) OnPersistentPreRun(h func(cmd *cobra.Command, args []string) error, options ...func(*HookOptions)) {
	OnPersistentPreRun(c.Command, h, options...)
}

// OnPersistentPostRun registers a PreRun hook on the command and all of its childs
func OnPersistentPreRun(c *cobra.Command, h func(cmd *cobra.Command, args []string) error, options ...func(*HookOptions)) {
	var opts HookOptions
	for _, option := range options {
		option(&opts)
	}
	// Register the hook
	// Prepend to the array to ensure hooks are executed
	// in the same order as how they are registered for the command
	// (since the runChain will be executed in reversed order from parent to child)
	persistentPreRunHooks = append([]*commandHook{
		&commandHook{
			cmd:       c,
			hook:      h,
			runOnHelp: opts.runOnHelp,
		},
	}, persistentPreRunHooks...)

	if opts.runOnHelp {
		initHelpHooks(c)
	}

	if c.PersistentPreRunE != nil {
		return
	}
	c.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		return runPersistentPreRunHooks(cmd, args, false)
	}
}

// OnPersistentPostRun registers a PostRun hook on the command and all of its childs
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

// OnHelp registers a hook when help is invoked for the command
func (c *Command) OnHelp(h func(cmd *cobra.Command, args []string) error, options ...func(*HookOptions)) {
	OnHelp(c.Command, h, options...)
}

var isHelpHooksInitialized bool

func initHelpHooks(c *cobra.Command) {
	if isHelpHooksInitialized {
		return
	}
	r := c.Root()
	if r == nil {
		return
	}
	// Integrate with the root command
	helpFunc := r.HelpFunc()
	r.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		if err := runPersistentPreRunHooks(cmd, args, true); err != nil {
			return
		}
		if err := runPreRunHooks(cmd, args, true); err != nil {
			return
		}
		for p, isParent := cmd, false; p != nil; p, isParent = p.Parent(), true {
			for _, ch := range helpHooks {
				if ch.cmd == p && (!isParent || ch.persistent == true) {
					if err := ch.hook(cmd, args); err != nil {
						return
					}
				}
			}
		}
		helpFunc(cmd, args)
	})
	isHelpHooksInitialized = true
}

// OnHelp registers a hook for when help is invoked.
func OnHelp(c *cobra.Command, h func(cmd *cobra.Command, args []string) error, options ...func(*HookOptions)) {
	initHelpHooks(c)
	var opts HookOptions
	for _, option := range options {
		option(&opts)
	}
	// Register the hook
	helpHooks = append(helpHooks, &commandHook{
		cmd:        c,
		hook:       h,
		persistent: opts.persistent,
	})
}
