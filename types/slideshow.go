package types

type SlideShowResponse struct {
	Speed      int
	FileImages []string
	UniqueTags []string
}

type SlideShowConfig struct {
	Speed        int
	IncludedTags []string
	UniqueTags   []string
}
