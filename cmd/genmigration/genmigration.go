package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"bytes"
	"github.com/gedex/inflector"
	"github.com/serenize/snaker"
	"go/format"
	"gopkg.in/godo.v2/util"
	"io/ioutil"
	"path/filepath"

	"text/template"

	"github.com/abiosoft/ishell"
)

// FirstCharLower lowers first char of string
func FirstCharLower(s string) string {
	ss := strings.Split(s, "")
	ss[0] = strings.ToLower(ss[0])

	return strings.Join(ss, "")
}

// FirstCharUpper uppers first char of string
func FirstCharUpper(s string) string {
	ss := strings.Split(s, "")
	ss[0] = strings.ToUpper(ss[0])

	return strings.Join(ss, "")
}

// GetFirstChar returns the first char of string
func GetFirstChar(s string) string {
	return s[0:1]
}

// FuncMap is a set of functions to use in templates
var FuncMap = template.FuncMap{
	"pluralize":   inflector.Pluralize,
	"singularize": inflector.Singularize,
	"title":       strings.Title,
	"firstLower":  FirstCharLower,
	"toLower":     strings.ToLower,
	"toSnakeCase": snaker.CamelToSnake,
	"firstChar":   GetFirstChar,
}

func GenerateFile(outputPath string, data interface{}) {
	path := filepath.Join("migrations/template.tmpl")
	body, _ := ioutil.ReadFile(path)
	tmpl := template.Must(template.New("model").Option("missingkey=error").Funcs(FuncMap).Parse(string(body)))

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	utils.CheckErr(err)

	src, _ := format.Source(buf.Bytes())
	dstPath := filepath.Join(outputPath)

	if !util.FileExists(filepath.Dir(dstPath)) {
		if err := os.Mkdir(filepath.Dir(dstPath), 0644); err != nil {
			fmt.Println(err)
		}
	}
	if err := ioutil.WriteFile(dstPath, src, 0644); err != nil {
		fmt.Println(err)
	}
}

// SelectedModel is a structure that holds model name and selected methods
type SelectedModel struct {
	Namespace   string
	ModelName   string
	Methods     []string
	MigrationId string
}

// SelectModels asks user which models and methods to generate
func SelectModels() (selectedModels []SelectedModel) {
	shell := ishell.New()

	shell.Println("Generating migrations")

	files, err := ioutil.ReadDir("models")
	if err != nil {
		fmt.Println(err)
	}

	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, strings.Title(inflector.Singularize(strings.TrimRight(file.Name(), ".go"))))
	}

	// Step 1: Ask user which models to use
	choices := shell.Checklist(fileNames,
		"Please select models you want to generate the matching store:",
		nil)
	if len(choices) == 0 {
		shell.Println("Please choose at least one model (by pressing spacebar on each one you want to select)")
		return nil
	}

	for _, file := range choices {
		var selectedModel SelectedModel
		selectedModel.ModelName = FirstCharUpper(inflector.Pluralize(fileNames[file]))
		selectedModels = append(selectedModels, selectedModel)
	}

	return selectedModels
}

func main() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	nowTimeString := time.Now().Format("200601021504") // current date

	selectedModels := SelectModels()

	for _, selectedModel := range selectedModels {
		fmt.Println(selectedModel)
		fileName := fmt.Sprintf("migrations/%s_%s.go", nowTimeString, strings.ToLower(selectedModel.ModelName))
		longFileName := path.Join(workingDirectory, fileName)
		selectedModel.MigrationId = nowTimeString
		GenerateFile(longFileName, selectedModel)

		fmt.Printf("Successfully created migration %s\n", fileName)
	}
}
