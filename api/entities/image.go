package entities

type Image struct {
	fileName     string
	originalURL  string
	thumbnailURL string
}

func NewImage(fileName string, originalURL string, thumbnailURL string) Image {
	return Image{
		fileName,
		originalURL,
		thumbnailURL,
	}
}

func (i Image) FileName() string {
	return i.fileName
}

func (i Image) OriginalURL() string {
	return i.originalURL
}

func (i Image) ThumbnailURL() string {
	return i.thumbnailURL
}
