// Package clock retrieves time data from the server and forwards it to the client
package clock

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/pisign/pisign-backend/types"
)

// API for clock
type API struct {
	types.BaseAPI
	Config types.ClockConfig
	time   time.Time
}

// NewAPI creates a new clock api
func NewAPI(configChan chan types.ClientMessage, pool types.Pool) *API {
	a := new(API)
	a.BaseAPI.Init("clock", configChan, pool)
	if a.Pool != nil {
		a.Pool.Register(a)
	}
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
func (a *API) Configure(body types.ClientMessage) error {
	defer a.Pool.Save()
	a.ConfigurePosition(body.Position)
	log.Println("Configuring CLOCK:", body)
	oldConfig := a.Config

	if len(body.Config) == 0 {
		return nil
	}
	if err := json.Unmarshal(body.Config, &a.Config); err != nil {
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
func (a *API) Run(w types.Socket) {
	log.Println("Running CLOCK")
	ticker := time.NewTicker(1 * time.Second)
	defer func() {
		ticker.Stop()
		log.Printf("STOPPING CLOCK: %s\n", a.UUID)
	}()
	for {
		select {
		case body := <-a.ConfigChan:
			err := a.Configure(body)
			if err != nil {
				w.SendErrorMessage(err.Error())
			}
		case save := <-w.Close():
			a.Pool.Unregister(types.Unregister{
				API:  a,
				Save: save,
			})
			return
		case t := <-ticker.C:
			a.time = t
			w.Send(a.Data())
		}
	}
}
