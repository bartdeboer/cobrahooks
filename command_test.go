package cobrahooks

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = executeCommandC(root, args...)
	return output, err
}

func executeCommandWithContext(ctx context.Context, root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	err = root.ExecuteContext(ctx)

	return buf.String(), err
}

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

func checkStringContains(t *testing.T, got, expected string) {
	if !strings.Contains(got, expected) {
		t.Errorf("Expected to contain: \n %v\nGot:\n %v\n", expected, got)
	}
}

func checkStringOmits(t *testing.T, got, expected string) {
	if strings.Contains(got, expected) {
		t.Errorf("Expected to not contain: \n %v\nGot: %v", expected, got)
	}
}

func emptyRun(_ *cobra.Command, _ []string) {}

func TestHooks(t *testing.T) {
	var (
		persPreArgs  string
		preArgs      string
		runArgs      string
		postArgs     string
		persPostArgs string
	)

	c := &Command{
		&cobra.Command{
			Use: "c",
			PersistentPreRun: func(_ *cobra.Command, args []string) {
				persPreArgs = strings.Join(args, " ")
			},
			PreRun: func(_ *cobra.Command, args []string) {
				preArgs = strings.Join(args, " ")
			},
			Run: func(_ *cobra.Command, args []string) {
				runArgs = strings.Join(args, " ")
			},
			PostRun: func(_ *cobra.Command, args []string) {
				postArgs = strings.Join(args, " ")
			},
			PersistentPostRun: func(_ *cobra.Command, args []string) {
				persPostArgs = strings.Join(args, " ")
			},
		},
	}

	output, err := executeCommand(c.Command, "one", "two")
	if output != "" {
		t.Errorf("Unexpected output: %v", output)
	}
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if persPreArgs != "one two" {
		t.Errorf("Expected persPreArgs %q, got %q", "one two", persPreArgs)
	}
	if preArgs != "one two" {
		t.Errorf("Expected preArgs %q, got %q", "one two", preArgs)
	}
	if runArgs != "one two" {
		t.Errorf("Expected runArgs %q, got %q", "one two", runArgs)
	}
	if postArgs != "one two" {
		t.Errorf("Expected postArgs %q, got %q", "one two", postArgs)
	}
	if persPostArgs != "one two" {
		t.Errorf("Expected persPostArgs %q, got %q", "one two", persPostArgs)
	}
}

func TestPersistentHooks(t *testing.T) {
	// var (
	// 	parentPersPreArgs  string
	// 	parentPreArgs      string
	// 	parentRunArgs      string
	// 	parentPostArgs     string
	// 	parentPersPostArgs string
	// )

	// var (
	// 	childPersPreArgs  string
	// 	childPreArgs      string
	// 	childRunArgs      string
	// 	childPostArgs     string
	// 	childPersPostArgs string
	// )

	var (
		persParentPersPreArgs  string
		persParentPreArgs      string
		persParentRunArgs      string
		persParentPostArgs     string
		persParentPersPostArgs string
	)

	var (
		persChildPersPreArgs  string
		persChildPreArgs      string
		persChildPreArgs2     string
		persChildRunArgs      string
		persChildPostArgs     string
		persChildPersPostArgs string
	)

	parentCmd := &Command{
		&cobra.Command{
			Use: "parent",
			// PersistentPreRun: func(_ *cobra.Command, args []string) {
			// 	parentPersPreArgs = strings.Join(args, " ")
			// },
			// PreRun: func(_ *cobra.Command, args []string) {
			// 	parentPreArgs = strings.Join(args, " ")
			// },
			// Run: func(_ *cobra.Command, args []string) {
			// 	parentRunArgs = strings.Join(args, " ")
			// },
			// PostRun: func(_ *cobra.Command, args []string) {
			// 	parentPostArgs = strings.Join(args, " ")
			// },
			// PersistentPostRun: func(_ *cobra.Command, args []string) {
			// 	parentPersPostArgs = strings.Join(args, " ")
			// },
		},
	}

	childCmd := &Command{
		&cobra.Command{
			Use: "child",
			// PersistentPreRun: func(_ *cobra.Command, args []string) {
			// 	childPersPreArgs = strings.Join(args, " ")
			// },
			// PreRun: func(_ *cobra.Command, args []string) {
			// 	childPreArgs = strings.Join(args, " ")
			// },
			// Run: func(_ *cobra.Command, args []string) {
			// 	childRunArgs = strings.Join(args, " ")
			// },
			// PostRun: func(_ *cobra.Command, args []string) {
			// 	childPostArgs = strings.Join(args, " ")
			// },
			// PersistentPostRun: func(_ *cobra.Command, args []string) {
			// 	childPersPostArgs = strings.Join(args, " ")
			// },
		},
	}
	parentCmd.AddCommand(childCmd.Command)

	parentCmd.OnPersistentPreRun(func(_ *cobra.Command, args []string) error {
		persParentPersPreArgs = strings.Join(args, " ")
		return nil
	})
	parentCmd.OnPreRun(func(_ *cobra.Command, args []string) error {
		persParentPreArgs = strings.Join(args, " ")
		return nil
	})
	parentCmd.OnRun(func(_ *cobra.Command, args []string) error {
		persParentRunArgs = strings.Join(args, " ")
		return nil
	})
	parentCmd.OnPostRun(func(_ *cobra.Command, args []string) error {
		persParentPostArgs = strings.Join(args, " ")
		return nil
	})
	parentCmd.OnPersistentPostRun(func(_ *cobra.Command, args []string) error {
		persParentPersPostArgs = strings.Join(args, " ")
		return nil
	})

	childCmd.OnPersistentPreRun(func(_ *cobra.Command, args []string) error {
		persChildPersPreArgs = strings.Join(args, " ")
		return nil
	})
	childCmd.OnPreRun(func(_ *cobra.Command, args []string) error {
		persChildPreArgs = strings.Join(args, " ")
		return nil
	})
	childCmd.OnPreRun(func(_ *cobra.Command, args []string) error {
		persChildPreArgs2 = strings.Join(args, " ") + " three"
		return nil
	})
	childCmd.OnRun(func(_ *cobra.Command, args []string) error {
		persChildRunArgs = strings.Join(args, " ")
		return nil
	})
	childCmd.OnPostRun(func(_ *cobra.Command, args []string) error {
		persChildPostArgs = strings.Join(args, " ")
		return nil
	})
	childCmd.OnPersistentPostRun(func(_ *cobra.Command, args []string) error {
		persChildPersPostArgs = strings.Join(args, " ")
		return nil
	})

	output, err := executeCommand(parentCmd.Command, "child", "one", "two")
	if output != "" {
		t.Errorf("Unexpected output: %v", output)
	}
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	// if EnablePersistentRunOverride && parentPersPreArgs != "" {
	// 	t.Errorf("Expected blank parentPersPreArgs, got %q", parentPersPreArgs)
	// }
	// if parentPersPreArgs != "one two" {
	// 	t.Errorf("Expected parentPersPreArgs %q, got %q", "one two", parentPersPreArgs)
	// }
	// if parentPreArgs != "" {
	// 	t.Errorf("Expected blank parentPreArgs, got %q", parentPreArgs)
	// }
	// if parentRunArgs != "" {
	// 	t.Errorf("Expected blank parentRunArgs, got %q", parentRunArgs)
	// }
	// if parentPostArgs != "" {
	// 	t.Errorf("Expected blank parentPostArgs, got %q", parentPostArgs)
	// }
	// if EnablePersistentRunOverride && parentPersPostArgs != "" {
	// 	t.Errorf("Expected blank parentPersPostArgs, got %q", parentPersPostArgs)
	// }
	// if parentPersPostArgs != "one two" {
	// 	t.Errorf("Expected parentPersPostArgs %q, got %q", "one two", parentPersPostArgs)
	// }
	// if childPersPreArgs != "one two" {
	// 	t.Errorf("Expected childPersPreArgs %q, got %q", "one two", childPersPreArgs)
	// }
	// if childPreArgs != "one two" {
	// 	t.Errorf("Expected childPreArgs %q, got %q", "one two", childPreArgs)
	// }
	// if childRunArgs != "one two" {
	// 	t.Errorf("Expected childRunArgs %q, got %q", "one two", childRunArgs)
	// }
	// if childPostArgs != "one two" {
	// 	t.Errorf("Expected childPostArgs %q, got %q", "one two", childPostArgs)
	// }
	// if childPersPostArgs != "one two" {
	// 	t.Errorf("Expected childPersPostArgs %q, got %q", "one two", childPersPostArgs)
	// }

	// Test On*Run hooks

	if persParentPersPreArgs != "one two" {
		t.Errorf("Expected persParentPersPreArgs %q, got %q", "one two", persParentPersPreArgs)
	}
	if persParentPreArgs != "" {
		t.Errorf("Expected blank persParentPreArgs, got %q", persParentPreArgs)
	}
	if persParentRunArgs != "" {
		t.Errorf("Expected blank persParentRunArgs, got %q", persParentRunArgs)
	}
	if persParentPostArgs != "" {
		t.Errorf("Expected blank persParentPostArg, got %q", persParentPostArgs)
	}
	if persParentPersPostArgs != "one two" {
		t.Errorf("Expected persParentPersPostArgs %q, got %q", "one two", persParentPersPostArgs)
	}

	if persChildPersPreArgs != "one two" {
		t.Errorf("Expected persChildPersPreArgs %q, got %q", "one two", persChildPersPreArgs)
	}
	if persChildPreArgs != "one two" {
		t.Errorf("Expected persChildPreArgs %q, got %q", "one two", persChildPreArgs)
	}
	if persChildPreArgs2 != "one two three" {
		t.Errorf("Expected persChildPreArgs %q, got %q", "one two three", persChildPreArgs2)
	}
	if persChildRunArgs != "one two" {
		t.Errorf("Expected persChildRunArgs %q, got %q", "one two", persChildRunArgs)
	}
	if persChildPostArgs != "one two" {
		t.Errorf("Expected persChildPostArgs %q, got %q", "one two", persChildPostArgs)
	}
	if persChildPersPostArgs != "one two" {
		t.Errorf("Expected persChildPersPostArgs %q, got %q", "one two", persChildPersPostArgs)
	}
}

func TestHelpCommandExecutedOnChild(t *testing.T) {
	rootCmd := &Command{&cobra.Command{Use: "root", Run: emptyRun}}
	childCmd := &Command{&cobra.Command{Use: "child", Long: "Long description", Run: emptyRun}}
	grandChildCmd := &Command{&cobra.Command{Use: "grandchild", Run: emptyRun}}
	rootCmd.AddCommand(childCmd.Command)
	childCmd.AddCommand(grandChildCmd.Command)

	var (
		rootFlagValue       string
		childFlagValue      string
		grandChildFlagValue string
		grandChildHelpTest  string
	)

	childCmd.PersistentFlags().StringVarP(&rootFlagValue, "root-flag", "", "", "")
	childCmd.PersistentFlags().StringVarP(&childFlagValue, "child-flag", "", "", "")
	grandChildCmd.PersistentFlags().StringVarP(&grandChildFlagValue, "grand-child-flag", "", "", "")
	grandChildCmd.PersistentFlags().StringVarP(&grandChildHelpTest, "grand-child-help-test", "", "", "")

	rootCmd.OnPersistentPreRun(func(cmd *cobra.Command, args []string) error {
		cmd.OutOrStdout().Write([]byte("Hello from hook 1 "))
		flag := cmd.Flags().Lookup("root-flag")
		flag.DefValue = "root new default value"
		return nil
	}, RunOnHelp)

	childCmd.OnPersistentPreRun(func(cmd *cobra.Command, args []string) error {
		cmd.OutOrStdout().Write([]byte("Hello from hook 2 "))
		flag := cmd.Flags().Lookup("child-flag")
		flag.DefValue = "child new default value"
		return nil
	}, RunOnHelp)

	childCmd.OnPreRun(func(cmd *cobra.Command, args []string) error {
		cmd.OutOrStdout().Write([]byte("Hello from hook 3 "))
		flag := cmd.Flags().Lookup("child-flag")
		flag.DefValue = "child new default value overwrite"
		return nil
	}, RunOnHelp, Persistent) // <- Persistent option

	childCmd.OnPreRun(func(cmd *cobra.Command, args []string) error {
		cmd.OutOrStdout().Write([]byte("Hello from hook 4 "))
		flag := cmd.Flags().Lookup("child-flag")
		flag.DefValue = "child new default value should not overwrite"
		return nil
	}, RunOnHelp)

	grandChildCmd.OnPersistentPreRun(func(cmd *cobra.Command, args []string) error {
		cmd.OutOrStdout().Write([]byte("Hello from hook 5 "))
		flag := cmd.Flags().Lookup("grand-child-flag")
		flag.DefValue = "grandchild new default value"
		return nil
	}, RunOnHelp)

	grandChildCmd.OnHelp(func(cmd *cobra.Command, args []string) error {
		cmd.OutOrStdout().Write([]byte("Hello from hook 6 "))
		flag := cmd.Flags().Lookup("grand-child-help-test")
		flag.DefValue = "grandchild help test value"
		return nil
	}, RunOnHelp)

	output1, err1 := executeCommand(rootCmd.Command, "child", "grandchild")
	if err1 != nil {
		t.Errorf("Unexpected error: %v", err1)
	}

	checkStringContains(t, output1,
		"Hello from hook 1 "+
			"Hello from hook 2 "+
			"Hello from hook 3 "+
			"Hello from hook 5 ")
	checkStringOmits(t, output1, "Hello from hook 4")
	checkStringOmits(t, output1, "Hello from hook 6")

	output2, err2 := executeCommand(rootCmd.Command, "child", "grandchild", "--help")
	if err2 != nil {
		t.Errorf("Unexpected error: %v", err2)
	}

	checkStringContains(t, output2,
		"Hello from hook 1 "+
			"Hello from hook 2 "+
			"Hello from hook 3 "+
			"Hello from hook 5 "+
			"Hello from hook 6 ")
	checkStringOmits(t, output2, "Hello from hook 4")
	checkStringContains(t, output2, "(default \"root new default value\")")
	checkStringContains(t, output2, "(default \"child new default value overwrite\")")
	checkStringContains(t, output2, "(default \"grandchild new default value\")")
	checkStringContains(t, output2, "(default \"grandchild help test value\")")
}
