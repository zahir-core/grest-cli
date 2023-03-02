package cmd

import (
	"embed"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

//go:embed all:openapi-ui
var f embed.FS

type cmdInit struct{}

func CmdInit() *cobra.Command {
	return &cobra.Command{
		Use:     "init",
		Example: "grest init example.com/org/backend",
		Short:   cmdInit{}.Summary(),
		Long:    cmdInit{}.Description(),
		Run:     cmdInit{}.Run,
	}
}

func (cmdInit) Summary() string {
	return "Initialize new app in the current directory"
}

func (cmdInit) Description() string {
	return `
Create a new grest app by automatically generating the basic code.
It will guess which kind of file to create based on the path provided.

Ensure you run this within the root directory of your app.
`
}

func (cmdInit) Run(c *cobra.Command, args []string) {
	if len(args) == 0 {
		c.Help()
	} else {
		writeGoModFile(args[0])
		writeOpenAPIUIFiles(args[0])
	}
}

func writeGoModFile(name string) error {
	if err := os.MkdirAll(name, 0755); err != nil {
		return err
	}
	gomod := "module " + name + "\n\ngo 1.20"
	filename := name + "/go.mod"
	fmt.Println("writting file :", filename)
	err := os.WriteFile(filename, []byte(gomod), 0755)
	return err
}

func writeOpenAPIUIFiles(name string) error {
	fmt.Println("writting @stoplight/elements files...")
	filePath := name + "/docs"

	os.RemoveAll(filePath)
	if err := os.MkdirAll(filePath, 0755); err != nil {
		return err
	}

	files := []string{
		"stoplight-elements-web-components.min.js",
		"stoplight-elements-styles.min.css",
		"index.html",
	}
	for _, fileName := range files {
		file, err := f.ReadFile("openapi-ui/" + fileName)
		if err != nil {
			return err
		}
		fmt.Println("writting file :", filePath+"/"+fileName)
		err = os.WriteFile(filePath+"/"+fileName, file, 0755)
		if err != nil {
			return err
		}
	}
	fmt.Println("@stoplight/elements files has been written")
	return nil
}
