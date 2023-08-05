package entities

type Avatar struct {
	fileName string
	url      string
}

func NewAvatar(fileName string, url string) Avatar {
	return Avatar{
		fileName,
		url,
	}
}

func (a Avatar) FileName() string {
	return a.fileName
}

func (a Avatar) URL() string {
	return a.url
}
