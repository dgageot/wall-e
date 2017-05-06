package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dgageot/wall-e/jenkins"
	"github.com/dgageot/wall-e/proxy"
)

func main() {
	server := os.Getenv("JENKINS_SERVER")
	user := os.Getenv("JENKINS_USER")
	token := os.Getenv("JENKINS_TOKEN")

	h := http.NewServeMux()

	h.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		res, err := jenkins.Get(user, token, server, r.RequestURI)
		if err != nil {
			http.Error(rw, err.Error(), 500)
			return
		}

		proxy.CopyResponse(rw, res)
	})

	err := http.ListenAndServe(":8080", h)
	log.Fatal(err)
}
