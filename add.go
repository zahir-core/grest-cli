package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

type cmdAdd struct{}

func CmdAdd() *cobra.Command {
	return &cobra.Command{
		Use:   "add",
		Short: cmdAdd{}.Summary(),
		Long:  cmdAdd{}.Description(),
		Run:   cmdAdd{}.Run,
	}
}

func (cmdAdd) Summary() string {
	return "Add a new end point for the current app"
}

func (cmdAdd) Description() string {
	return `
Create a new end point for the current grest app by automatically generating the basic code.
It will guess which kind of file to create based on the path provided.

Ensure you run this within the root directory of your app.
`
}

func (cmdAdd) Run(c *cobra.Command, args []string) {
	err := addEndPoint()
	if err == nil {
		fmt.Println("Success!")
	} else {
		fmt.Println("Failed!", err.Error())
	}
}

func addEndPoint() error {
	fmt.Println("add end point: TODO")
	os.Setenv("IS_GENERATE_OPEN_API_DOC", "true")
	fmt.Println("prepare open api doc...")
	return exec.Command("go", "run", "main.go").Run()
}
