package api

import (
	"errors"
	"log"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pisign/pisign-backend/utils"
)

// Template code inspired by https://levelup.gitconnected.com/learn-and-use-templates-in-go-aa6146b01a38
type templateData struct {
	NameLower string
	NameTitle string
	NameUpper string
}

func AutoCreate(name string) error {
	name = strings.ToLower(name)
	data := templateData{
		NameLower: name,
		NameTitle: strings.Title(name),
		NameUpper: strings.ToUpper(name),
	}
	log.Printf("Creating new api %s!\n", name)
	allFiles := make([]string, 5)
	files, err := createMainFiles(data)
	if err != nil {
		return err
	}
	allFiles = append(allFiles, files...)

	files, err = createGlobalTypesFile(data)
	if err != nil {
		return err
	}
	allFiles = append(allFiles, files...)

	files, err = updateFactory(data)
	if err != nil {
		return err
	}
	allFiles = append(allFiles, files...)

	// TODO: Automatically fix file imports
	return nil
}

func createMainFiles(data templateData) ([]string, error) {
	// Create directory
	basePath := filepath.Join("api", data.NameLower)
	if err := utils.CreateDirectory(basePath); err != nil {
		return nil, err
	}

	// Generate Main file
	path := filepath.Join(basePath, utils.AddExtension(data.NameLower, "go"))
	outFile, err := utils.CreateFile(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		utils.WrapError(outFile.Close())
	}()

	tpl, err := template.ParseFiles("api/example/example.template")
	if err != nil {
		return nil, err
	}
	if err := tpl.Execute(outFile, data); err != nil {
		return nil, err
	}

	return []string{path}, nil
}

func updateFactory(data templateData) ([]string, error) {
	return nil, errors.New("Unimplemented")
}

func createGlobalTypesFile(data templateData) ([]string, error) {
	basePath := filepath.Join("types")

	// Generate Types file
	path := filepath.Join(basePath, utils.AddExtension(data.NameLower, "go"))
	log.Printf("Path: %s\n", path)
	outFile, err := utils.CreateFile(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		utils.WrapError(outFile.Close())
	}()

	tpl, err := template.ParseFiles("types/example.template")
	if err != nil {
		return nil, err
	}
	if err := tpl.Execute(outFile, data); err != nil {
		return nil, err
	}

	return []string{path}, nil
}
