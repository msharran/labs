package main

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"path/filepath"
)

type Option struct {
	Environment  string
	DesiredCount int
}

var opts = []Option{
	{
		Environment:  "stage",
		DesiredCount: 1,
	},
	{
		Environment:  "prod",
		DesiredCount: 2,
	},
}

//go:embed templates/*
var templateFiles embed.FS

func main() {

	fmt.Println(filepath.Glob("templates*"))

	// files, err := templateFiles.ReadDir("templates")
	// check(err)

	// var files []string
	// err := fs.WalkDir(templateFiles, "templates", func(path string, d fs.DirEntry, err error) error {
	// 	if !d.IsDir() && strings.HasSuffix(path, ".tpl") {
	// 		files = append(files, path)
	// 	}
	// 	return nil
	// })
	// check(err)

	// files, err := fs.Glob(templateFiles, "templates/snippets/*")
	// check(err)
	// files = append(files, "templates/base/foo.tpl")
	// fmt.Println("files::::::::::::", files)

	// content, err := mergeFiles(files...)
	// check(err)

	// // b, err := templateFiles.ReadFile("templates/foo.yaml.tpl")
	// // check(err)

	// fmt.Println("content::::::::::::")
	// fmt.Println(content)

	// t := template.Must(template.New("foo").Parse(content))

	// for _, opt := range opts {
	// 	fmt.Println("out::::::::::::")
	// 	t.Execute(os.Stdout, opt)
	// }
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func mergeFiles(fileNames ...string) (content string, err error) {
	var readers []io.Reader
	for _, file := range fileNames {
		fmt.Println(file)
		content, err := templateFiles.ReadFile(file)
		if err != nil {
			return "", err
		}
		readers = append(readers, bytes.NewBuffer(content))
	}

	reader := io.MultiReader(readers...)
	b, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
