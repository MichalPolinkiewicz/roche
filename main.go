package main

import (
	"github.com/MichalPolinkiewicz/roche/cmd/server/rest"
	"github.com/MichalPolinkiewicz/roche/service/postman"
	"log"
	"net/http"
	"time"
)

func main() {
	// TODO - read cfg from file/envs, replace default logger

	postmanServiceConfig, err := postman.NewPostmanServiceConfig("https://postman-echo.com/get?", time.Hour, []string{"message"})
	if err != nil {
		log.Fatal(err)
	}
	postmanService, err := postman.NewPostmanService(postmanServiceConfig, http.Client{})
	if err != nil {
		log.Fatal(err)
	}

	messageHandler, err := rest.NewMessageHandler(postmanService)
	messageEndpoint, err := rest.NewEndpoint("/message", http.MethodPost, messageHandler.Handle)
	if err != nil {
		log.Fatal(err)
	}

	restServerConfig, err := rest.NewRestServerConfig("8082")
	if err != nil {
		log.Fatal(err)
	}
	restServer := rest.NewRestServer(restServerConfig, []*rest.Endpoint{messageEndpoint})
	if err := restServer.Run(); err != nil {
		log.Fatal(err)
	}
}
