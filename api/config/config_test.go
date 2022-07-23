package config_test

import (
	"os"
	"strings"
	"testing"

	"github.com/edgarSucre/comic-parser/api/config"
)

func TestGetEnvironment(t *testing.T) {
	_, err := config.GetEnvironment()
	if !strings.Contains(err.Error(), "API_PORT") {
		t.Error("Failure, Expected error to contain API_PORT")
	}

	os.Setenv("API_PORT", "8080")

	_, err = config.GetEnvironment()
	if !strings.Contains(err.Error(), "COMIC_HOST") {
		t.Error("Failure, Expected error to contain COMIC_HOST")
	}

	os.Setenv("COMIC_HOST", "https://xkcd.com/")

	env, err := config.GetEnvironment()

	if err != nil {
		t.Error("Failed to retrieve environment")
	}

	if env.ServerPort != "8080" {
		t.Errorf("Failure, Expected: 8080, Got: %s", env.ServerPort)
	}

	if env.ComicHost != "https://xkcd.com/" {
		t.Errorf("Failure, Expected: https://xkcd.com/, Got: %s", env.ComicHost)
	}

	os.Setenv("COMIC_HOST", "this should not affect the environment")
	env, _ = config.GetEnvironment()

	if env.ComicHost != "https://xkcd.com/" {
		t.Errorf("Failure, Expected: https://xkcd.com/, Got: %s", env.ComicHost)
	}
}
