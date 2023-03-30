package cmd

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"grest.dev/cmd/codegentemplate/app"
	"grest.dev/grest"
)

//go:embed all:codegentemplate
var f embed.FS

type cmdInit struct{}

func CmdInit() *cobra.Command {
	return &cobra.Command{
		Use:     "init",
		Example: "grest init",
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
	err := runInit()
	if err == nil {
		fmt.Println("Success!")
	} else {
		fmt.Println("Failed!", err.Error())
	}
}

func runInit() error {
	var qs = []*survey.Question{
		{
			Name: "module-path",
			Prompt: &survey.Input{
				Message: "Module path:",
				Help: "A module path is the canonical name for a module, declared with the module directive in the module’s go.mod file. " +
					"A module’s path is the prefix for package paths within the module.\n\n" +
					grest.Fmt("See ", grest.FmtHiWhite) + grest.Fmt("https://go.dev/ref/mod#module-path\n", grest.FmtHiBlue),
			},
			Validate: func(val any) error {
				if str, ok := val.(string); !ok || !regexp.MustCompile(`^(?i)[a-z0-9]+([a-z0-9._-]*[a-z0-9]+)?(/([a-z0-9._-]*[a-z0-9]+)?)*$`).MatchString(str) {
					return errors.New("\"" + str + "\" is not a valid module path.\n\n" +
						grest.Fmt("See ", grest.FmtHiWhite) + grest.Fmt("https://go.dev/ref/mod#go-mod-file-ident\n", grest.FmtHiBlue))
				}
				return nil
			},
		},
		{
			Name:   "project-name",
			Prompt: &survey.Input{Message: "Project name:"},
			Validate: func(val any) error {
				if val == nil {
					return errors.New("Project name is required")
				}
				return nil
			},
		},
		{
			Name:   "project-description",
			Prompt: &survey.Input{Message: "Project description:"},
		},
		{
			Name: "database",
			Prompt: &survey.Select{
				Message: "Choose a database:",
				Options: []string{"postgres", "mysql", "sqlserver", "clickhouse", "sqlite", "other"},
				Default: "sqlite",
			},
		},
		{
			Name:   "is-add-end-point",
			Prompt: &survey.Confirm{Message: "Add your first end point?"},
		},
	}
	answer := struct {
		ModulePath         string `survey:"module-path"`
		ProjectName        string `survey:"project-name"`
		ProjectDescription string `survey:"project-description"`
		Database           string `survey:"database"`
		IsAddEndPoint      bool   `survey:"is-add-end-point"`
	}{}

	err := survey.Ask(qs, &answer)
	if err != nil {
		return err
	}
	if answer.ProjectDescription == "" {
		answer.ProjectDescription = "The My App API allows you to perform all the operations that you do with our applications." +
			"My App API is built using REST principles which ensures predictable URLs, uses standard HTTP response codes, authentication, " +
			"and verbs that makes writing applications easy."
	}

	err = writeGoModFile(answer.ModulePath)
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
			newContent := strings.ReplaceAll(string(content), "grest.dev/cmd/codegentemplate", answer.ModulePath)
			newContent = strings.ReplaceAll(newContent, `o.Info.Description = ""`, `o.Info.Description = "`+answer.ProjectDescription+`"`)
			newContent = strings.ReplaceAll(newContent, "My App API", answer.ProjectName)
			if answer.Database != "other" {
				newContent = strings.ReplaceAll(newContent, "sqlite", answer.Database)
			}
			newContent = strings.ReplaceAll(newContent, "23.03.161330", time.Now().Format("2006.01.021504"))
			newContent = strings.ReplaceAll(newContent, "f4cac8b77a8d4cb5881fac72388bb226", app.Crypto().NewToken())
			newContent = strings.ReplaceAll(newContent, "wAGyTpFQX5uKV3JInABXXEdpgFkQLPTf", app.Crypto().NewToken())
			newContent = strings.ReplaceAll(newContent, "0de0cda7d2dd4937a1c4f7ddc43c580f", app.Crypto().NewToken())
			return os.WriteFile(newFileName, []byte(newContent), 0755)
		})
	if answer.Database == "other" {
		fmt.Println()
		fmt.Println("You choose", grest.Fmt("other", grest.FmtHiCyan, grest.FmtBold), "database, modify",
			grest.Fmt("app/db.go", grest.FmtHiGreen, grest.FmtBold), "with your own database driver setup.")
		fmt.Println("See", grest.Fmt("https://gorm.io/docs/connecting_to_the_database.html", grest.FmtHiBlue))
		fmt.Println()
	}
	if answer.IsAddEndPoint {
		fmt.Println("----------Add End Point----------")
		err = addEndPoint(false)
		if err != nil {
			return err
		}
	}
	err = exec.Command("go", "mod", "tidy").Run()
	if err != nil {
		return err
	}
	return updateOpenAPI()
}

func writeGoModFile(name string) error {
	filename := "go.mod"
	fmt.Println("writting file :", filename)
	return os.WriteFile(filename, []byte("module "+name+"\n\ngo "+runtime.Version()[2:6]), 0755)
}
