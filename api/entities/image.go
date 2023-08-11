package entities

type Image struct {
	fileName             string
	originalURL          string
	originalContentType  string
	formattedURL         string
	formattedContentType string
}

func NewImage(fileName string, originalURL string, originalContentType string, formattedURL string, formattedContentType string) Image {
	return Image{
		fileName,
		originalURL,
		originalContentType,
		formattedURL,
		formattedContentType,
	}
}

func (i Image) FileName() string {
	return i.fileName
}

func (i Image) OriginalURL() string {
	return i.originalURL
}

func (i Image) FormattedURL() string {
	return i.formattedURL
}

func (i Image) OriginalContentType() string {
	return i.originalContentType
}

func (i Image) FormattedContentType() string {
	return i.formattedContentType
}
