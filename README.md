# Cobrahooks

A more versatile hooks layer for [Cobra](https://github.com/spf13/cobra). It allows to register multiple hooks for cobra commands. Also ensures persistent hooks are persistently executed.

Please note that this package makes use of Cobra's (Persistent)(Pre/Post)RunE Command fields. Defining your own hooks using these fields will interfere with any hooks registered through this package.

```go
parentCmd := &Command{
    &cobra.Command{
        Use: "parent",
    },
}

childCmd := &Command{
    &cobra.Command{
        Use: "child",
    },
}

parentCmd.AddCommand(childCmd.Command)

parentCmd.OnPersistentPreRun(func(_ *cobra.Command, args []string) error {
    fmt.Println("Your external behavior")
    return nil
})
```
