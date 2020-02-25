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
func NewAPI(sockets map[types.Socket]bool, pool types.Pool, id uuid.UUID) *API {
	a := new(API)
	a.BaseAPI.Init("clock", sockets, pool, id)
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
		if a.Pool != nil && a.Sockets != nil {
			a.Pool.Save()
			a.Send(a.Data())
		}
	}()
	if err := a.BaseAPI.Configure(message); err != nil {
		return err
	}

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
	}
	return nil
}

// Data gets the current time!
func (a *API) Data() (interface{}, error) {
	return types.ClockResponse{
		Time: a.time.In(a.loc()).String(),
	}, nil
}

// Run main entry point to clock API
func (a *API) Run() {
	go a.BaseAPI.Run()
	log.Println("Running CLOCK")
	a.Send(a.Data())
	ticker := time.NewTicker(1 * time.Second)
	defer func() {
		ticker.Stop()
		log.Printf("STOPPING CLOCK: %s\n", a.UUID)
	}()
	for a.Running() {
		select {
		case body := <-a.ConfigChan:
			if err := a.Configure(body); err != nil {
				a.Send(nil, err)
			}
		case msg := <-a.CloseChan:
			for socket := range msg.Sockets {
				delete(a.Sockets, socket)
				socket.Close()
			}
			if len(a.Sockets) == 0 {
				a.Stop()
				a.Pool.Unregister(types.Unregister{
					API:  a,
					Save: msg.Save,
				})
			}
		case t := <-ticker.C:
			a.time = t
			a.Send(a.Data())
		default:
			continue
		}
	}
}
