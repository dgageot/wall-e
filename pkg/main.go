package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"./jenkins"
	"./proxy"
)

func main() {
	server := os.Getenv("JENKINS_SERVER")
	user := os.Getenv("JENKINS_USER")
	token, err := ioutil.ReadFile("/run/secrets/jenkins_token")
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("[" + string(token) + "]")

	h := http.NewServeMux()

	h.HandleFunc("/api/", func(rw http.ResponseWriter, r *http.Request) {
		res, err := jenkins.Get(user, string(token), server, r.RequestURI)
		if err != nil {
			http.Error(rw, err.Error(), 500)
			return
		}

		proxy.CopyResponse(rw, res)
	})

	err = http.ListenAndServe(":8080", h)
	log.Fatal(err)
}
