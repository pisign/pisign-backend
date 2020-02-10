package widget

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Layout of multiple widgets
type Layout struct {
	Name    string
	Widgets []*Widget
}

func getFilename(name string) string {
	return fmt.Sprintf("layouts/%s.json", name)
}

// LoadLayout returns the current stored layout for the given name
func LoadLayout(name string) Layout {
	log.Printf("Loading layout %s\n", name)
	dataFile, err := os.Open(getFilename(name))

	if err != nil {
		log.Printf("Error Opening file for layout %s: %v\n", name, err)
		return Layout{}
	}
	defer dataFile.Close()

	var layout Layout
	dataDecoder := json.NewDecoder(dataFile)
	err = dataDecoder.Decode(&layout)
	if err != nil {
		log.Printf("Error Decoding layout %s: %v\n", name, err)
		return Layout{}
	}
	fmt.Printf("Layout: %+v\n", layout)
	return layout
}

// SaveLayout stores the widgets to a local json file
func SaveLayout(layout Layout) error {
	name := layout.Name
	log.Printf("Saving layout %s\n", name)
	dataFile, err := os.Create(getFilename(name))
	if err != nil {
		log.Printf("Error Creating file for layout %s: %v\n", name, err)
		return err
	}
	defer dataFile.Close()

	dataEncoder := json.NewEncoder(dataFile)
	dataEncoder.Encode(layout)
	return nil
}
