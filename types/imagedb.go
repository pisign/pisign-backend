package types

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
)

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
	db.NumImages++

	return nil
}
