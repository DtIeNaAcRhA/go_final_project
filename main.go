package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/DtIeNaAcRhA/go_final_project/pkg/api"
	"github.com/DtIeNaAcRhA/go_final_project/pkg/db"
)

func main() {
	envDBFile := os.Getenv("TODO_DBFILE")
	if envDBFile == "" {
		envDBFile = "scheduler.db"
		log.Print("The default value of the db")
	} else {
		log.Print("DB from a variable environment")
	}

	err := db.Init(envDBFile)
	if err != nil {
		log.Fatalf("database opening error: %v", err)
	}
	log.Println("Successful connection to the database")
	defer db.Close()

	api.SetEnvPassword()
	api.SetSecret()

	webDir := "web"

	envPort := os.Getenv("TODO_PORT")
	if envPort == "" {
		envPort = "7540"
		log.Print("The default value of the port")
	} else {
		log.Print("Port from a variable environment")
	}
	log.Printf("The server starts on the port %s", envPort)

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	api.Init()

	if err := http.ListenAndServe(fmt.Sprintf(":%s", envPort), nil); err != nil {
		log.Printf("Start server error: %s", err.Error())
		return
	}
}
