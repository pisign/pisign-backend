package slideshow

import (
	"errors"
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
	a.BaseAPI.Init("slideshow", sockets, pool, id)

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
		log.Printf("Configuring Slideshow:\n %s\n", string(message.Config))
		if err := utils.ParseJSON(message.Config, &a.Config); err != nil {
			log.Println("Could not properly configure Slideshow")
			return errors.New("could not properly configure Slideshow")
		}

		// Add custom checks for config fields here (see the `time` api as an example)

		log.Println("Slideshow configuration successful:", a)
	}
	log.Println("Slideshow config set to ", a.Config)

	return nil
}

// Data performs a bulk of the computational logic
func (a *API) Data() (interface{}, error) {
	// Get the filepaths to include for a set of tags
	// Also send all the possible tags so the client can choose later which tags to use
	// when configuring the api
	imgDB := a.Pool.GetImageDB().(*types.ImageDB)
	images := make([]string, imgDB.NumImages-1)
	uniqueTagsSet := make(map[string]bool)

	for _, item := range imgDB.Images {
		// Add to the tag to the tags set
		added := false
		for _, tag := range item.Tags {
			uniqueTagsSet[tag] = true
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

	uniqueTags := make([]string, len(uniqueTagsSet)-1)

	// Get the unique tags from the set
	for tag := range uniqueTagsSet {
		uniqueTags = append(uniqueTags, tag)
	}

	// TODO there must be some issue with making slices that causes
	// us to have empty lists... fix that so we don't have to rely on
	// these functions to clean that up!
	uniqueTags = utils.DeleteEmpty(uniqueTags)
	imagesNotEmpty := utils.DeleteEmpty(images)

	imgDB.UniqueTags = uniqueTags
	a.FilePathsForTags = imagesNotEmpty

	log.Println("Sending with tags", uniqueTags)
	log.Println("Sending files", imagesNotEmpty)

	return types.SlideShowResponse{
		FileImages: imagesNotEmpty,
		Speed:      a.Config.Speed,
		UniqueTags: uniqueTags,
	}, nil
}

// Run main entry point to the API
func (a *API) Run() {
	// Start up the BaseAPI to handle core API stuff
	go a.BaseAPI.Run()

	log.Println("Running Slideshow")

	// Send data to client (using default config values)
	a.Send(a.Data())
	ticker := time.NewTicker(1 * time.Second)

	defer func() {
		ticker.Stop()
		log.Printf("STOPPING Slideshow: %s\n", a.UUID)
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
