package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const INSERT_TEXT string = "INSERT NEW LINES HERE"

// CreateDirectory checks if a directory exists and if not creates it
// taken from https://www.socketloop.com/tutorials/golang-check-if-directory-exist-and-create-if-does-not-exist
func CreateDirectory(dirName string) error {
	src, err := os.Stat(dirName)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(dirName, 0755)
		if errDir != nil {
			return err
		}
		log.Printf("Created directory: '%s'\n", dirName)
		return nil
	}

	if err != nil {
		return err
	}

	if src.Mode().IsRegular() {
		return fmt.Errorf("CreateDirectory: '%s' already exists as a file, not a directory", dirName)
	}
	return nil
}

func CreateFile(path string) (*os.File, error) {
	emptyFile, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return emptyFile, nil
}

func AddExtension(fname string, ext string) string {
	return fmt.Sprintf("%s.%s", fname, ext)
}

// Text insertion inspired by https://siongui.github.io/2017/01/30/go-insert-line-or-string-to-file/
func InsertText(path string, text string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { WrapError(f.Close()) }()

	var lines []string
	scanner := bufio.NewScanner(f)
	targetFound := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, INSERT_TEXT) {
			targetFound = true
			lines = append(lines, strings.Split(text, "\n")...)
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if !targetFound {
		return errors.New(fmt.Sprintf("No insertion point found in %s\n", path))
	}

	fileContent := strings.Join(lines, "\n")
	return ioutil.WriteFile(path, []byte(fileContent), 0644)
}

// FileSearch searches for the given string anywhere in the file
func FileSearch(path string, str string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer func() { WrapError(f.Close()) }()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		log.Printf("Line: %s, search string: %s\n", line, str)
		if strings.Contains(line, str) {
			return true
		}
	}
	return false
}

// AppendText appends text
func AppendText(path string, text string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() { WrapError(f.Close()) }()

	_, err = f.WriteString(text)
	if err != nil {
		return err
	}
	return nil
}
