package api

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"

	"github.com/pisign/pisign-backend/types"
	"github.com/pisign/pisign-backend/utils"
)

// Unmarshal an API of unknown type
func Unmarshal(body *json.RawMessage) types.API {
	var name struct {
		Name string
	}
	err := utils.ParseJSON(*body, &name)
	if err != nil {
		log.Println("Could not unmarshal API: ", err)
		return nil
	}

	newAPI, err := NewAPI(name.Name, nil, nil, uuid.New())
	if err != nil {
		log.Printf("Unknown API type: %v\n", err)
		return nil
	}

	utils.ParseJSON(*body, newAPI)
	return newAPI
}
