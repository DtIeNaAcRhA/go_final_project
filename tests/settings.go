package tests

import (
	"log"
	"os"
	"strconv"
)

var Port = getPort()
var DBFile = "../scheduler.db"
var FullNextDate = true
var Search = true
var Token = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXNzd29yZF9oYXNoIjoiZmJmYjM4NmVmZWE2N2U4MTZmMmRkYTBhOGM5NGE5OGViMjAzNzU3YWViYjNmNTVmMTgzNzU1YTE5MmQ0NDQ2NyJ9.NlPttiCSXWmiwk40AlhYlLxVLdDk3XReo54moSQXRn0`

func getPort() int {
	envPort := os.Getenv("TODO_PORT")
	if envPort == "" {
		return 7540 //устанавливаем значение по умолчанию
	}
	port, err := strconv.Atoi(envPort)
	if err != nil {
		log.Fatalf("port initialization error: %v", err)
	}
	return port
}
