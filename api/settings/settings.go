package settings

import "os"

type ISettings interface {
	Port() string
}

type settings struct {
	port string
}

func NewSettings() ISettings {
	return &settings{
		port: os.Getenv("PORT"),
	}
}

func (s *settings) Port() string {
	return s.port
}
