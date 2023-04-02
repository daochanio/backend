package settings

import "os"

type ISettings interface {
	Port() string
	DbConnectionString() string
}

type settings struct {
	port               string
	dbConnectionString string
}

func NewSettings() ISettings {
	return &settings{
		port:               os.Getenv("PORT"),
		dbConnectionString: os.Getenv("DB_CONNECTION_STRING"),
	}
}

func (s *settings) Port() string {
	return s.port
}

func (s *settings) DbConnectionString() string {
	return s.dbConnectionString
}
