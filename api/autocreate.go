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
	if err := createMainFiles(data); err != nil {
		return err
	}
	if err := updateFactory(data); err != nil {
		return err
	}
	if err := createGlobalTypesFile(data); err != nil {
		return err
	}
	return nil
}

func createMainFiles(data templateData) error {
	// Create directory
	basePath := filepath.Join("api", data.NameLower)
	if err := utils.CreateDirectory(basePath); err != nil {
		return err
	}

	// Generate Main file
	path := filepath.Join(basePath, utils.AddExtension(data.NameLower, "go"))
	outFile, err := utils.CreateFile(path)
	if err != nil {
		return err
	}
	defer func() {
		utils.WrapError(outFile.Close())
	}()

	tpl, err := template.ParseFiles("api/example/example.template")
	if err != nil {
		return err
	}
	if err := tpl.Execute(outFile, data); err != nil {
		return err
	}
	// TODO: Automatically fix file imports right when file is created

	return nil
}

func updateFactory(data templateData) error {
	return errors.New("Unimplemented")
}

func createGlobalTypesFile(data templateData) error {
	return errors.New("Unimplemented")
}
