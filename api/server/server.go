package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/edgarSucre/comic-parser/api/domain"
)

type Server struct {
	provider domain.ComicProvider
	cache    map[any][]byte
}

func NewServer(p domain.ComicProvider) *Server {
	return &Server{p, make(map[any][]byte)}
}

func (s *Server) ValidateIdMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chapter, err := getIdParam(r)
		if err != nil || chapter < 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "%s", err.Error())
			return
		}

		next(w, r)
	}
}

func (s *Server) DataCacheMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cControl := r.Header.Get("Cache-Control")

		if cControl != "no-cache" {
			chapter, _ := getIdParam(r)
			if content, ok := s.cache[chapter]; ok {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				w.Write(content)
				return
			}
		}

		next(w, r)
	}
}

func (s *Server) ImgCacheMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cControl := r.Header.Get("Cache-Control")

		if cControl != "no-cache" {
			url := r.URL.Query().Get("src")
			if content, ok := s.cache[url]; ok {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/octet-stream")
				w.Write(content)
				return
			}
		}

		next(w, r)
	}
}

func (s *Server) GetChapter(w http.ResponseWriter, r *http.Request) {
	// valid chapter alredy verified by middleware
	chapter, _ := getIdParam(r)

	comic, err := s.provider.GetCommic(chapter)
	if err != nil {
		setErrorResponse(w, err)
		return
	}

	response, err := json.Marshal(comic)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"err": "%s"}`, err.Error())
		return
	}

	s.cache[chapter] = response

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (s *Server) GetImage(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("src")
	w.Header().Set("Content-Type", "application/octet-stream")
	response, err := http.Get(url)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"err": "%s"}`, "couldn't find the image")
	}
	defer response.Body.Close()

	img, err := ioutil.ReadAll(response.Body)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"err": "%s"}`, "couldn't find the image")
	}

	s.cache[url] = img

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(img)
}

func setErrorResponse(w http.ResponseWriter, err error) {
	if strings.Contains(err.Error(), "couldn't find the comic") {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"err": "%s"}`, err.Error())
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, `{"err": "%s"}`, err.Error())
}

func getIdParam(r *http.Request) (int, error) {
	tempChap := strings.TrimPrefix(r.URL.Path, "/api/")
	if tempChap == "" {
		tempChap = "0"
	}

	chapter, err := strconv.Atoi(tempChap)
	if err != nil || chapter < 0 {
		return 0, fmt.Errorf(`{"err": "%s"}`, "invalid Chapter: must be a postive number")
	}
	return chapter, nil
}
