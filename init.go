package cmd

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"grest.dev/cmd/codegentemplate/app"
)

//go:embed all:codegentemplate
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
		err := runInit(args[0])
		if err == nil {
			fmt.Println("Success!")
		} else {
			fmt.Println("Failed!", err.Error())
		}
	}
}

func writeGoModFile(name string) error {
	filename := "go.mod"
	fmt.Println("writting file :", filename)
	return os.WriteFile(filename, []byte("module "+name+"\n\ngo "+runtime.Version()[2:6]), 0755)
}

func runInit(name string) error {
	err := writeGoModFile(name)
	if err != nil {
		return err
	}
	fs.WalkDir(f, "codegentemplate",
		func(fileName string, info fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			newFileName := strings.Replace(fileName, "codegentemplate/", "", 1)
			if info.IsDir() {
				if newFileName == "codegentemplate" {
					return nil
				}
				return os.MkdirAll(newFileName, 0755)
			}
			content, err := f.ReadFile(fileName)
			if err != nil {
				return err
			}
			fmt.Println("writting file :", newFileName)
			newContent := strings.ReplaceAll(string(content), "grest.dev/cmd/codegentemplate", name)
			newContent = strings.ReplaceAll(newContent, "23.03.161330", time.Now().Format("2006.01.021504"))
			newContent = strings.ReplaceAll(newContent, "f4cac8b77a8d4cb5881fac72388bb226", app.NewToken())
			newContent = strings.ReplaceAll(newContent, "wAGyTpFQX5uKV3JInABXXEdpgFkQLPTf", app.NewToken())
			newContent = strings.ReplaceAll(newContent, "0de0cda7d2dd4937a1c4f7ddc43c580f", app.NewToken())
			return os.WriteFile(newFileName, []byte(newContent), 0755)
		})
	fmt.Println("go mod tidy...")
	err = exec.Command("go", "mod", "tidy").Run()
	if err != nil {
		return err
	}
	os.Setenv("IS_GENERATE_OPEN_API_DOC", "true")
	fmt.Println("prepare open api doc...")
	return exec.Command("go", "run", "main.go").Run()
}
