package pool

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/pisign/pisign-backend/api"

	"github.com/pisign/pisign-backend/types"
	"github.com/pisign/pisign-backend/utils"
)

const storageFolder string = "layouts"

// Layout of multiple Sockets
// Each layout is stored serverside to be retrieved later by the client
type Layout struct {
	Name string
	List []types.API
}

func getFilename(name string) string {
	err := utils.CreateDirectory(storageFolder)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s/%s.json", storageFolder, name)
}

// LoadLayout returns the stored layout of the given name
func LoadLayout(name string) Layout {
	log.Printf("Loading layout %s\n", name)

	dataFile, err := os.Open(getFilename(name))

	if err != nil {
		log.Printf("Error Opening file for layout %s: %v\n", name, err)
		return Layout{}
	}
	defer dataFile.Close()

	var layout struct {
		Name string
		List []*json.RawMessage
	}
	dataDecoder := json.NewDecoder(dataFile)
	err = dataDecoder.Decode(&layout)
	if err != nil {
		log.Printf("Error Decoding layout %s: %v\n", name, err)
		return Layout{}
	}
	if layout.Name != name {
		log.Printf("Layout name requested (%s) and retrieved from file (%s) do not match!\n", name, layout.Name)
		return Layout{}
	}
	var list []types.API
	for _, body := range layout.List {
		list = append(list, api.Unmarshal(body))
	}
	return Layout{Name: layout.Name, List: list}
}

// SaveLayout stores layout to a local json file
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
