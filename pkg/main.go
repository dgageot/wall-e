package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"./github"
	"./jenkins"
	"./proxy"
	"github.com/gorilla/mux"
)

func main() {
	if err := startProxy(); err != nil {
		log.Fatal(err)
	}
}

func startProxy() error {
	server := os.Getenv("JENKINS_SERVER")
	user := os.Getenv("JENKINS_USER")
	tokenJenkins, err := ioutil.ReadFile("/run/secrets/jenkins_token")
	if err != nil {
		return err
	}

	tokenGithub, err := ioutil.ReadFile("/run/secrets/github_token")
	if err != nil {
		return err
	}

	r := mux.NewRouter()
	r.HandleFunc("/jobs", func(rw http.ResponseWriter, r *http.Request) {
		res, err := jenkins.Get(user, string(tokenJenkins), server, "/api/json?tree=jobs[name,builds[building,number,result,runs[building,number,result]]{0,5}]")
		if err != nil {
			http.Error(rw, err.Error(), 500)
			log.Print(err)
			return
		}

		proxy.CopyResponse(rw, res)
	})

	r.HandleFunc("/pulls", func(rw http.ResponseWriter, r *http.Request) {
		res, err := github.Get(string(tokenGithub), "/repos/docker/pinata/pulls")
		if err != nil {
			http.Error(rw, err.Error(), 500)
			log.Print(err)
			return
		}

		proxy.CopyResponse(rw, res)
	})

	r.HandleFunc("/status/{sha1}", func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		res, err := github.Get(string(tokenGithub), "/repos/docker/pinata/commits/"+vars["sha1"]+"/status")
		if err != nil {
			http.Error(rw, err.Error(), 500)
			log.Print(err)
			return
		}

		proxy.CopyResponse(rw, res)
	})

	return http.ListenAndServe(":8080", r)
}
