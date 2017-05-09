package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"./github"
	"./jenkins"
	"./proxy"
)

func main() {
	go func() {
		if err := startJenkinsProxy(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := startGithubProxy(); err != nil {
		log.Fatal(err)
	}
}

func startJenkinsProxy() error {
	server := os.Getenv("JENKINS_SERVER")
	user := os.Getenv("JENKINS_USER")
	token, err := ioutil.ReadFile("/run/secrets/jenkins_token")
	if err != nil {
		return err
	}

	h := http.NewServeMux()
	h.HandleFunc("/api/", func(rw http.ResponseWriter, r *http.Request) {
		res, err := jenkins.Get(user, string(token), server, r.RequestURI)
		if err != nil {
			http.Error(rw, err.Error(), 500)
			return
		}

		proxy.CopyResponse(rw, res)
	})

	return http.ListenAndServe(":8080", h)
}

func startGithubProxy() error {
	token, err := ioutil.ReadFile("/run/secrets/github_token")
	if err != nil {
		return err
	}

	h := http.NewServeMux()
	h.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		res, err := github.Get(string(token), r.RequestURI)
		if err != nil {
			http.Error(rw, err.Error(), 500)
			return
		}

		proxy.CopyResponse(rw, res)
	})

	return http.ListenAndServe(":8888", h)
}
