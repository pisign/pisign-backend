package types

type SlideShowResponse struct {
	SlideShowConfig
	FileImages []string
}

type SlideShowConfig struct {
	Speed        int
	IncludedTags []string
	UniqueTags   []string
}
