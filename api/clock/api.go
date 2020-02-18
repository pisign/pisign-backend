// Package clock retrieves time data from the server and forwards it to the client
package clock

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/pisign/pisign-backend/types"
)

// API for clock
type API struct {
	types.BaseAPI
	Config types.ClockConfig
	time   time.Time
}

// NewAPI creates a new clock api
func NewAPI(socket types.Socket, pool types.Pool, id uuid.UUID) *API {
	a := new(API)
	a.BaseAPI.Init("clock", socket, pool, id)
	a.Config.Location = "Local"
	return a
}

func (a *API) loc() *time.Location {
	t, err := time.LoadLocation(a.Config.Location)
	if err != nil {
		t, _ = time.LoadLocation("Local")
	}
	return t
}

// Configure for clock
func (a *API) Configure(message types.ClientMessage) error {
	defer func() {
		if a.Pool != nil && a.Socket != nil {
			a.Pool.Save()
			a.Socket.Send(a.Data())
		}
	}()
	a.BaseAPI.Configure(message)

	switch message.Action {
	case types.ConfigureAPI, types.Initialize:
		log.Println("Configuring CLOCK:", message)
		oldConfig := a.Config
		if err := json.Unmarshal(message.Config, &a.Config); err != nil {
			log.Println("Could not properly configure clock")
			a.Config = oldConfig
			return errors.New("could not properly configure clock")
		}

		if _, err := time.LoadLocation(a.Config.Location); err != nil {
			log.Printf("Could not load timezone %s: %s", a.Config.Location, err)
			a.Config.Location = oldConfig.Location // Revert to old location
			return errors.New("could not load timezone " + a.Config.Location)
		}

		log.Println("Clock configuration successful:", a)
	case types.ChangeAPI:
		a.Pool.Switch(a, message.Name)
	default:
		return errors.New("Invalid ClientMessage.Action")
	}
	return nil
}

// Data gets the current time!
func (a *API) Data() interface{} {
	return types.ClockResponse{
		Time: a.time.In(a.loc()).String(),
		BaseMessage: types.BaseMessage{
			Status:       types.StatusSuccess,
			ErrorMessage: "",
		},
	}
}

// Run main entry point to clock API
func (a *API) Run() {
	log.Println("Running CLOCK")
	ticker := time.NewTicker(1 * time.Second)
	defer func() {
		ticker.Stop()
		log.Printf("STOPPING CLOCK: %s\n", a.UUID)
	}()
	for {
		select {
		case body := <-a.Socket.Config():
			err := a.Configure(body)
			if err != nil {
				a.Socket.SendErrorMessage(err.Error())
			}
		case <-a.StopChan:
			return
		case save := <-a.Socket.Close():
			a.Pool.Unregister(types.Unregister{
				API:  a,
				Save: save,
			})
			return
		case t := <-ticker.C:
			a.time = t
			a.Socket.Send(a.Data())
		}
	}
}
