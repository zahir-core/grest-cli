package cmd

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Tag struct {
	Key   string
	Value string
}

func FormatApp() {
	filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(path, ".go") {
				FormatFile(path)
			}
			return nil
		})
}

func FormatFile(fileName string) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	for _, node := range f.Decls {
		genDecl, isGenDecl := node.(*ast.GenDecl)
		if isGenDecl {
			for _, spec := range genDecl.Specs {
				typeSpec, isTypeSpec := spec.(*ast.TypeSpec)
				if isTypeSpec {
					fmt.Println("Formatting", fileName, typeSpec.Name.Name)
					structType, isStructType := typeSpec.Type.(*ast.StructType)
					if isStructType {
						mapTag, maxTagLen := ParseTag(structType.Fields.List)
						RewriteTag(structType.Fields.List, mapTag, maxTagLen)
					}
				}
			}
		}
	}

	var buf bytes.Buffer
	err = format.Node(&buf, fset, f)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(fileName, buf.Bytes(), 0)
	if err != nil {
		log.Fatal(err)
	}
}

func ParseTag(fields []*ast.Field) (mapTag map[string][]Tag, maxTagLen map[string]int) {
	mapTag = map[string][]Tag{}
	maxTagLen = map[string]int{}
	for _, field := range fields {
		if len(field.Names) > 0 {
			if field.Tag == nil {
				field.Tag = &ast.BasicLit{}
			}

			var tags []Tag
			ftv, _ := strconv.Unquote(field.Tag.Value)
			tgs := strings.Split(strings.ReplaceAll(ftv, `:"`, "==="), `"`)
			for _, tg := range tgs {
				t := strings.Split(strings.Trim(tg, " "), "===")
				if len(t) > 1 {
					key := t[0]
					value := t[1]
					lenVal := len(value)
					ml, isMaxLenExist := maxTagLen[key]
					if !isMaxLenExist || lenVal > ml {
						maxTagLen[key] = lenVal
					}
					tags = append(tags, Tag{Key: key, Value: value})
				}
			}
			mapTag[field.Names[0].Name] = tags
		}
	}

	return mapTag, maxTagLen
}

func RewriteTag(fields []*ast.Field, mapTag map[string][]Tag, maxTagLen map[string]int) {
	for _, field := range fields {
		if len(field.Names) > 0 {
			tags, isExist := mapTag[field.Names[0].Name]
			if isExist {
				if field.Tag == nil {
					field.Tag = &ast.BasicLit{}
				}
				field.Tag.Value = FormattedTagString(tags, maxTagLen)
			}
		}
	}
}

func FormattedTagString(tags []Tag, maxTagLen map[string]int) string {
	if len(tags) == 0 {
		return ""
	}
	sortedTags := []Tag{}
	sort.Slice(tags, func(i, j int) bool { return tags[i].Key < tags[j].Key })
	for _, tag := range tags {
		if tag.Key == "json" {
			sortedTags = append(sortedTags, tag)
		}
	}
	for _, tag := range tags {
		if tag.Key == "db" {
			sortedTags = append(sortedTags, tag)
		}
	}
	for _, tag := range tags {
		if tag.Key == "gorm" {
			sortedTags = append(sortedTags, tag)
		}
	}
	for _, tag := range tags {
		if tag.Key == "validate" {
			sortedTags = append(sortedTags, tag)
		}
	}
	for _, tag := range tags {
		if tag.Key != "json" && tag.Key != "db" && tag.Key != "gorm" && tag.Key != "validate" {
			sortedTags = append(sortedTags, tag)
		}
	}
	newTag := ""
	for _, t := range sortedTags {
		newTag += t.Key + ":" + TagValueWithDelimiter(t.Value, maxTagLen[t.Key])
	}
	return "`" + strings.Trim(newTag, " ") + "`"
}

func TagValueWithDelimiter(str string, maxLen int) string {
	tag := `"` + str + `"`
	n := maxLen - len(str) + 1
	for i := 0; i < n; i++ {
		tag += " "
	}
	return tag
}
