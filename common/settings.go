package common

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// These are intended to be generic settings that every app is expected to implement.
// packages in common can safely expect to have access to these settings in their constructors.
type CommonSettings interface {
	Appname() string
	Hostname() string
	Env() string
	IsDev() bool
}

type commonSettings struct {
	env      string
	appname  string
	hostname string
}

func NewCommonSettings() CommonSettings {
	env := os.Getenv("ENV")
	appname := os.Getenv("APP_NAME")

	if env == "dev" {
		err := godotenv.Load(fmt.Sprintf(".env/.env.%v.dev", appname))
		if err != nil {
			panic(err)
		}
	}

	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}

	return &commonSettings{
		env,
		appname,
		hostname,
	}
}

func (s *commonSettings) Appname() string {
	return s.appname
}

func (s *commonSettings) Hostname() string {
	return s.hostname
}

func (s *commonSettings) Env() string {
	return s.env
}

func (s *commonSettings) IsDev() bool {
	return s.env == "dev"
}
