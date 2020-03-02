package types

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
)

func NewImageDB() *ImageDB {
	var imgDB ImageDB

	imgFile, err := os.Open("./assets/images/images.json")

	if err != nil {
		log.Println(err.Error())
		return &ImageDB{}
	}

	defer imgFile.Close()

	byteValue, _ := ioutil.ReadAll(imgFile)

	if err = json.Unmarshal(byteValue, &imgDB); err != nil {
		log.Println(err.Error())
		return &ImageDB{}
	}

	return &imgDB
}

// TaggedImage is an instance of an image with tags
type TaggedImage struct {
	Tags     []string
	FilePath string
}

// ImageDB holds the images that we have access to
type ImageDB struct {
	// A map of tags to file locations?
	Images     []TaggedImage
	NumImages  int
	UniqueTags []string // the tag "all" means to show all?
}

func (db *ImageDB) GetNumImages() int {
	return len(db.Images)
}

func (db *ImageDB) AddImage(filepath string, file multipart.File, tagsList []string) error {
	out, err := os.Create(filepath)

	defer out.Close()
	if err != nil {
		return errors.New("Server error saving file")
	}

	_, err = io.Copy(out, file) // file not files[i] !

	if err != nil {
		return err
	}

	// If everything was successful, append the new image to the json file
	db.Images = append(db.Images, TaggedImage{
		Tags:     tagsList,
		FilePath: filepath,
	})

	return nil
}
