package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/jinzhu/inflection"
	"github.com/spf13/cobra"
	"grest.dev/grest"
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
	err := addEndPoint(true)
	if err == nil {
		fmt.Println("Success!")
	} else {
		fmt.Println("Failed!", err.Error())
	}
}

func input(s survey.Prompt, res any, opts ...survey.AskOpt) {
	err := survey.AskOne(s, res, opts...)
	if err != nil && err == terminal.InterruptErr {
		os.Exit(0)
	}
}

func addEndPoint(isUpdateOpenAPI bool) error {
	templatePath := "src/codegentemplate"
	input(&survey.Input{
		Message: "Template path:",
		Default: "src/codegentemplate",
	}, &templatePath)

	endPointPath := ""
	input(&survey.Input{
		Message: "RESTful API path:",
		Help:    `The path that uniquely identifies a RESTful API, for example "/api/contacts"`,
	}, &endPointPath, survey.WithValidator(func(val any) error {
		if str, ok := val.(string); !ok || !regexp.MustCompile(`^/([a-zA-Z0-9_-]+/?)*$`).MatchString(str) {
			return fmt.Errorf(`"%v" not a valid RESTful API path.`+"\n", val)
		}
		return nil
	}))
	p := strings.Split(endPointPath, "/")
	endPoint := p[len(p)-1]
	singularName := inflection.Singular(strings.ReplaceAll(endPoint, "_", " "))

	packagePathPrefix := ""
	input(&survey.Input{
		Message: "Package path prefix:",
		Default: "src",
	}, &packagePathPrefix)

	packagePath := strings.ToLower(strings.ReplaceAll(singularName, " ", ""))
	input(&survey.Input{
		Message: "Package path:",
		Default: packagePath,
		Help:    `The path that uniquely identifies a package, for example "` + packagePath + `"`,
	}, &packagePath, survey.WithValidator(func(val any) error {
		if str, ok := val.(string); !ok || !regexp.MustCompile(`^(?i)[a-z0-9]+([a-z0-9._-]*[a-z0-9]+)?(/([a-z0-9._-]*[a-z0-9]+)?)*$`).MatchString(str) {
			return fmt.Errorf(`"%v" is not a valid package path.`+grest.Fmt("\n\nSee ", grest.FmtHiWhite)+grest.Fmt("https://go.dev/ref/mod#glos-package-path\n", grest.FmtHiBlue), val)
		}
		return nil
	}))

	modelStructName := strings.ReplaceAll(strings.Title(singularName), " ", "")
	input(&survey.Input{
		Message: "Model struct name:",
		Default: modelStructName,
	}, &modelStructName, survey.WithValidator(func(val any) error {
		if str, ok := val.(string); !ok || !regexp.MustCompile(`^[A-Z][A-Za-z0-9]*$`).MatchString(str) {
			return fmt.Errorf(`"%v" is not a valid struct name`+"\n", val)
		}
		return nil
	}))

	isAddField := true
	newFields := []map[string]string{}
	for isAddField {
		input(&survey.Confirm{
			Message: "Add field?",
		}, &isAddField)
		if !isAddField {
			continue
		}
		fmt.Println("----------Add Field----------")

		fieldName := ""
		input(&survey.Input{
			Message: "Field name:",
		}, &fieldName, survey.WithValidator(func(val any) error {
			if str, ok := val.(string); !ok || !regexp.MustCompile(`^([a-zA-Z0-9_-]+/?)*$`).MatchString(str) {
				return fmt.Errorf(`"%v" not a valid field name.`+"\n", val)
			}
			return nil
		}))

		fieldType := "NullString"
		input(&survey.Select{
			Message: "Field type:",
			Options: []string{
				"NullUUID",
				"NullString",
				"NullText",
				"NullJSON",
				"NullBool",
				"NullInt64",
				"NullFloat64",
				"NullDate",
				"NullTime",
				"NullDateTime",
			},
			Default: "NullString",
		}, &fieldType)

		newFields = append(newFields, map[string]string{"name": fieldName, "type": fieldType})
		fmt.Println()
	}
	newFieldStr := ""
	for _, nf := range newFields {
		temp := `StructFieldName app.field_type !json:"field_name" db:"m.field_name" gorm:"column:field_name"!` + "\n"
		temp = strings.ReplaceAll(temp, "StructFieldName", grest.String{}.PascalCase(nf["name"]))
		temp = strings.ReplaceAll(temp, "field_type", nf["type"])
		temp = strings.ReplaceAll(temp, "field_name", nf["name"])
		temp = strings.ReplaceAll(temp, "!", "`")
		newFieldStr += temp
	}

	packagePathWithPrefix := packagePathPrefix + "/" + packagePath
	err := os.MkdirAll(packagePathWithPrefix, 0755)
	if err != nil {
		return err
	}

	filepath.Walk(templatePath,
		func(fileName string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			newFileName := strings.Replace(fileName, filepath.FromSlash("src/codegentemplate/codegentemplate"), packagePathWithPrefix+"/"+packagePath, 1)
			if info.IsDir() {
				if newFileName == "codegentemplate" {
					return nil
				}
				return os.MkdirAll(newFileName, 0755)
			}
			file, err := os.Open(fileName)
			if err != nil {
				return err
			}
			content, err := io.ReadAll(file)
			if err != nil {
				return err
			}
			fmt.Println("writting file :", newFileName)
			newContent := strings.ReplaceAll(string(content), "codegentemplate", packagePath)
			newContent = strings.ReplaceAll(newContent, "CodeGenTemplate", modelStructName)
			newContent = strings.ReplaceAll(newContent, "end_point", endPoint)
			newContent = strings.ReplaceAll(newContent, "// AddField : DONT REMOVE THIS COMMENT\n", newFieldStr)
			return os.WriteFile(newFileName, []byte(newContent), 0755)
		})
	grest.FormatFile(packagePathWithPrefix)

	goModFile, err := os.Open("go.mod")
	if err != nil {
		return err
	}
	goModContent, err := io.ReadAll(goModFile)
	if err != nil {
		return err
	}
	baseModulePath := strings.Split(string(goModContent), "\n")[0]
	baseModulePath = strings.Replace(baseModulePath, "module ", "", 1)

	for _, fileName := range []string{"src/migrator.go", "src/router.go"} {
		fmt.Println("updating file :", fileName)
		file, err := os.Open(fileName)
		if err != nil {
			return err
		}
		content, err := io.ReadAll(file)
		if err != nil {
			return err
		}
		importSection := "// import : DONT REMOVE THIS COMMENT"
		newImportSection := `"` + baseModulePath + "/" + packagePathWithPrefix + `"` + "\n" + importSection
		newContent := strings.Replace(string(content), importSection, newImportSection, 1)

		registerTableSection := "// RegisterTable : DONT REMOVE THIS COMMENT"
		newRegisterTableSection := `app.DB().RegisterTable("main", ` + packagePath + "." + modelStructName + "{})\n" + registerTableSection
		newContent = strings.Replace(newContent, registerTableSection, newRegisterTableSection, 1)

		addRouteSection := "// AddRoute : DONT REMOVE THIS COMMENT"
		newAddRouteSection := `
			app.Server().AddRoute("/codegentemplate", "POST", codegentemplate.REST().Create, codegentemplate.OpenAPI().Create())
			app.Server().AddRoute("/codegentemplate", "GET", codegentemplate.REST().Get, codegentemplate.OpenAPI().Get())
			app.Server().AddRoute("/codegentemplate/{id}", "GET", codegentemplate.REST().GetByID, codegentemplate.OpenAPI().GetByID())
			app.Server().AddRoute("/codegentemplate/{id}", "PUT", codegentemplate.REST().UpdateByID, codegentemplate.OpenAPI().UpdateByID())
			app.Server().AddRoute("/codegentemplate/{id}", "PATCH", codegentemplate.REST().PartiallyUpdateByID, codegentemplate.OpenAPI().PartiallyUpdateByID())
			app.Server().AddRoute("/codegentemplate/{id}", "DELETE", codegentemplate.REST().DeleteByID, codegentemplate.OpenAPI().DeleteByID())

			// AddRoute : DONT REMOVE THIS COMMENT`
		newAddRouteSection = strings.ReplaceAll(newAddRouteSection, "/codegentemplate", endPointPath)
		newAddRouteSection = strings.ReplaceAll(newAddRouteSection, "codegentemplate", packagePath)
		newContent = strings.Replace(newContent, addRouteSection, newAddRouteSection, 1)
		err = os.WriteFile(fileName, []byte(newContent), 0755)
		if err != nil {
			return err
		}
		grest.FormatFile(fileName)
	}

	if isUpdateOpenAPI {
		return updateOpenAPI()
	}
	return nil
}

func updateOpenAPI() error {
	fmt.Println()
	fmt.Println("prepare open api doc...")
	fmt.Println()
	os.Setenv("IS_GENERATE_OPEN_API_DOC", "true")
	return exec.Command("go", "run", "main.go", "update").Run()
}
