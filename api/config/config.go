package config

import (
	"fmt"
	"os"
)

type ENV struct {
	ServerPort string
	ComicHost  string
}

var env *ENV

func GetEnvironment() (*ENV, error) {
	if env != nil {
		return env, nil
	}

	port := os.Getenv("API_PORT")
	if port == "" {
		return nil, fmt.Errorf("config error: %s", "Couldn't read 'API_PORT' environment variable")
	}

	comicHost := os.Getenv("COMIC_HOST")
	if comicHost == "" {
		return nil, fmt.Errorf("config error: %s", "Couldn't read 'COMIC_HOST' environment variable")

	}

	env = &ENV{port, comicHost}

	return env, nil

}
