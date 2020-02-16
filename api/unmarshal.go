package api

import (
	"encoding/json"
	"log"

	"github.com/pisign/pisign-backend/types"
	"github.com/pisign/pisign-backend/utils"
)

// Unmarshal an API of unknown type
func Unmarshal(body *json.RawMessage) types.API {
	log.Printf("Unmarshalling!\n")
	var name struct {
		Name string
	}
	err := utils.ParseJSON(*body, &name)
	if err != nil {
		log.Println("Could not unmarshal API: ", err)
		return nil
	}
	log.Printf("Creating new API %s\n", name.Name)

	newAPI, err := NewAPI(name.Name, nil, nil)
	if err != nil {
		log.Printf("Unknown API type: %v\n", err)
		return nil
	}

	utils.ParseJSON(*body, newAPI)
	return newAPI
}
