package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/edgarSucre/comic-parser/api/domain"
)

type Server struct {
	provider domain.ComicProvider
}

func NewServer(p domain.ComicProvider) *Server {
	return &Server{p}
}

func (s *Server) GetChapter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	tempChap := strings.TrimPrefix(r.URL.Path, "/")
	if tempChap == "" {
		tempChap = "0"
	}

	chapter, err := strconv.Atoi(tempChap)
	if err != nil || chapter < 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"err": "%s"}`, "invalid Chapter: must be a postive number")
		return
	}

	comic, err := s.provider.GetCommic(int(chapter))
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

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
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

func (s *Server) Start(port string) error {
	router := http.NewServeMux()
	router.HandleFunc("/", s.GetChapter)

	log.Printf("Listening on port: %s\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%s", port), router)
}
