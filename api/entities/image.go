package entities

type Image struct {
	fileName    string
	url         string
	contentType string
}

func NewImage(fileName string, url string, contentType string) Image {
	return Image{
		fileName,
		url,
		contentType,
	}
}

func (i Image) FileName() string {
	return i.fileName
}

func (i Image) Url() string {
	return i.url
}

func (i Image) ContentType() string {
	return i.contentType
}
