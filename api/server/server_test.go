package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/edgarSucre/comic-parser/api/domain"
	"github.com/edgarSucre/comic-parser/api/server"
)

type MockProvier struct{}

func (mp *MockProvier) GetCommic(id int) (domain.Comic, error) {
	//hardcode condition to force error
	if id == 99 {
		return domain.Comic{}, fmt.Errorf("internal error")
	}

	if id != 12 {
		return domain.Comic{}, fmt.Errorf("couldn't find the comic with id: %d", id)
	}

	return domain.Comic{
		Num:    12,
		Tittle: "Mocked Commic",
	}, nil
}

func TestGetChapter(t *testing.T) {
	srv := server.NewServer(&MockProvier{})

	t.Run("Invalid Chapter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/fail", nil)
		res := httptest.NewRecorder()

		srv.GetChapter(res, req)

		if res.Code != http.StatusBadRequest {
			t.Errorf("Failed, Expected: %d, Got: %d", http.StatusBadRequest, res.Code)
		}

		err := getErrorResponse(res.Body)
		if strings.Contains(err.Content, "Invalid Chapter") {
			t.Errorf("Failed, Expected: %s, to contain 'InvalidChaper'", err.Content)
		}
	})

	t.Run("Internal Error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/99", nil)
		res := httptest.NewRecorder()

		srv.GetChapter(res, req)

		code := http.StatusInternalServerError
		if res.Code != code {
			t.Errorf("Failed, Expected: %d, Got: %d", code, res.Code)
		}

		err := getErrorResponse(res.Body)
		contains := "Internal error"
		if strings.Contains(err.Content, contains) {
			t.Errorf("Failed, Expected: %s, to contain '%s'", err.Content, contains)
		}
	})

	t.Run("Chapter not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/1", nil)
		res := httptest.NewRecorder()

		srv.GetChapter(res, req)

		if res.Code != http.StatusNotFound {
			t.Errorf("Failed, Expected: %d, Got: %d", http.StatusNotFound, res.Code)
		}

		err := getErrorResponse(res.Body)
		contains := "Couldn't find the comic"
		if strings.Contains(err.Content, contains) {
			t.Errorf("Failed, Expected: %s, to contain '%s'", err.Content, contains)
		}
	})

	t.Run("Return Chapter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/12", nil)
		res := httptest.NewRecorder()

		srv.GetChapter(res, req)

		if res.Code != http.StatusOK {
			t.Errorf("Failed, Expected: %d, Got: %d", http.StatusOK, res.Code)
		}

		if ct := res.Header().Get("Content-Type"); ct != "application/json" {
			t.Errorf("Failed, Expected: %s, Got: %s", "application/json", ct)
		}

		response, _ := ioutil.ReadAll(res.Body)
		var comic domain.Comic
		json.Unmarshal(response, &comic)

		if comic.Tittle != "Mocked Commic" {
			t.Errorf("Failed, Expected: %s, Got: %s", "Mocked Commic", comic.Tittle)
		}

		if comic.Num != 12 {
			t.Errorf("Failed, Expected: %d, Got: %d", 12, comic.Num)
		}
	})
}

type errorResponse struct {
	Content string
}

func getErrorResponse(body *bytes.Buffer) errorResponse {
	response, _ := ioutil.ReadAll(body)

	var err errorResponse

	json.Unmarshal(response, &err)

	return err
}
