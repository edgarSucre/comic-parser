package main

import (
	"log"

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
	log.Fatal(server.Start(env.ServerPort))
}
