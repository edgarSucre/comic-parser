package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/edgarSucre/comic-parser/api/config"
	"github.com/edgarSucre/comic-parser/api/server"
	"github.com/edgarSucre/comic-parser/api/xkcd"
)

func main() {
	env, err := config.GetEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	comicClient := xkcd.NewClient(env)
	server := server.NewServer(comicClient)

	router := http.NewServeMux()
	router.HandleFunc("/api/", server.GetChapter)
	router.Handle("/", http.FileServer(http.Dir("./public")))

	log.Printf("Listening on port: %s\n", env.ServerPort)
	http.ListenAndServe(fmt.Sprintf(":%s", env.ServerPort), router)
}
