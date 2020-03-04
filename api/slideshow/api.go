package slideshow

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/pisign/pisign-backend/types"
	"github.com/pisign/pisign-backend/utils"
)

type API struct {
	types.BaseAPI
	FilePathsForTags []string

	Config types.SlideShowConfig
	// Add other fields as necessary, (keep lowercase to avoid being stored in json)
}

// NewAPI creates a new API
func NewAPI(sockets map[types.Socket]bool, pool types.Pool, id uuid.UUID) *API {
	a := new(API)
	a.BaseAPI.Init(types.APISlideshow, sockets, pool, id)

	return a
}

// Configure based on client message
func (a *API) Configure(message types.ClientMessage) error {
	// Make sure the client widget is immediately sent new data with new config options
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
		log.Printf("Configuring %s:\n %+v\n", a, string(message.Config))
		if err := utils.ParseJSON(message.Config, &a.Config); err != nil {
			log.Printf("Could not properly configure %s\n", a)
			return errors.New(fmt.Sprintf("could not properly configure %s\n", a))
		}

		// Add custom checks for config fields here (see the `time` api as an example)

		log.Printf("%s configuration successful: %+v\n", a, a)
	}
	log.Printf("%s config set to %+v\n", a, a.Config)

	return nil
}

// Data performs a bulk of the computational logic
func (a *API) Data() (interface{}, error) {
	// Get the filepaths to include for a set of tags
	// Also send all the possible tags so the client can choose later which tags to use
	// when configuring the api
	imgDB := a.Pool.GetImageDB()
	// Make sure the unique tags are up to date
	a.Config.UniqueTags = imgDB.UniqueTags
	images := make([]string, imgDB.GetNumImages())

	for _, item := range imgDB.Images {
		// Add to the tag to the tags set
		added := false
		for _, tag := range item.Tags {
			if utils.StringInSlice(tag, a.Config.IncludedTags) {
				if !added {
					if len(item.FilePath[1:]) > 0 {
						images = append(images, item.FilePath[1:])
						added = true
					}
				}
			}

			// Include all? TODO add the custom ALL tag or something to also include it
			if len(a.Config.IncludedTags) == 0 {
				if len(item.FilePath[1:]) != 0 {
					images = append(images, item.FilePath[1:])
				}
			}
		}
	}

	// TODO there must be some issue with making slices that causes
	// us to have empty lists... fix that so we don't have to rely on
	// these functions to clean that up!
	imagesNotEmpty := utils.DeleteEmpty(images)

	a.FilePathsForTags = imagesNotEmpty

	log.Println("Sending files", imagesNotEmpty)

	return types.SlideShowResponse{
		FileImages:      imagesNotEmpty,
		SlideShowConfig: a.Config,
	}, nil
}

// Run main entry point to the API
func (a *API) Run() {
	// Start up the BaseAPI to handle core API stuff
	go a.BaseAPI.Run()

	log.Printf("Running %s\n", a)

	// Send data to client (using default config values)
	a.Send(a.Data())
	ticker := time.NewTicker(10 * time.Second)

	defer func() {
		ticker.Stop()
		log.Printf("STOPPING %s\n", a)
	}()

	// Create a new channel to recieve termination messages on
	stop := a.AddStopChan()
	for {
		select {
		case body := <-a.ConfigChan: // Configuration messages
			if err := a.Configure(body); err != nil {
				a.Send(nil, err)
			}
		case <-ticker.C:
			a.Send(a.Data())
		case <-stop: // Terminate
			return
		}
	}
}
