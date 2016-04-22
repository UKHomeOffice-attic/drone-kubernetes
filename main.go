package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var debug bool
var deployments []string

func main() {

	system := drone.System{}
	repo := drone.Repo{}
	build := drone.Build{}
	vargs := Params{}
	workspace := drone.Workspace{}

	plugin.Param("system", &system)
	plugin.Param("workspace", &workspace)
	plugin.Param("repo", &repo)
	plugin.Param("build", &build)
	plugin.Param("vargs", &vargs)
	plugin.MustParse()

	if vargs.Debug {
		debug = true
	}

	if debug {
		log.Println("Workspace Root: " + workspace.Root)
		log.Println("Workspace Path: " + workspace.Path)

		log.Println("Tag: " + vargs.Tag)
	}

	// Iterate over rcs and svcs
	for _, rc := range vargs.ReplicationControllers {
		artifact, err := readArtifactFromFile(workspace.Path, rc, vargs.ApiServer, vargs.Namespace, vargs.Tag)
		if err != nil {
			log.Fatal(err)
		}
		if debug {
			log.Println("Artifact loaded: " + artifact.Url)
		}
		if b, _ := existsArtifact(artifact, vargs.Token); b {
			deleteArtifact(artifact, vargs.Token)
			time.Sleep(time.Second * 5)
		}
		createArtifact(artifact, vargs.Token)
	}
	for _, rc := range vargs.Services {
		artifact, err := readArtifactFromFile(workspace.Path, rc, vargs.ApiServer, vargs.Namespace, vargs.Tag)
		if err != nil {
			log.Fatal(err)
		}
		createArtifact(artifact, vargs.Token)
	}
	for _, d := range vargs.Deployments {
		artifact, err := readArtifactFromFile(workspace.Path, d, vargs.ApiServer, vargs.Namespace, vargs.Tag)
		if err != nil {
			log.Fatal(err)
		}
		if b, _ := existsArtifact(artifact, vargs.Token); b {
			updateArtifact(artifact, vargs.Token)
		} else {
			createArtifact(artifact, vargs.Token)
		}
	}
}

func doRequest(param ReqEnvelope) (bool, error) {
	if debug {
		log.Println("doRequest ")
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	var req *http.Request
	var err error
	// post payload to each artifact
	if param.Json == nil {
		req, err = http.NewRequest(param.Verb, param.Url, nil)
	} else {
		req, err = http.NewRequest(param.Verb, param.Url, bytes.NewBuffer(param.Json))
	}

	if param.Verb == "PATCH" {
		req.Header.Set("Content-Type", "application/strategic-merge-patch+json ")
	} else {
		req.Header.Set("Content-Type", "application/json ")
	}

	if debug {
		log.Println("HTTP Request %s", param.Verb)
		log.Println("HTTP Request %s", param.Url)
		log.Println("HTTP Request %s", string(param.Json))
	}

	req.Header.Set("Authorization", "Bearer "+param.Token)
	response, err := client.Do(req)
	if debug {
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s\n", string(contents))
	}

	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()

		if response.StatusCode == 200 {
			return true, err
		}
	}
	return false, err
}

func readArtifactFromFile(workspace string, artifactFile string, apiserver string, namespace string, tag string) (Artifact, error) {
	artifactFilename := workspace + "/" + artifactFile
	if debug {
		log.Println("readArtifactFromFile " + artifactFilename)
	}
	file, err := ioutil.ReadFile(artifactFilename)
	if err != nil {
		log.Fatal(err)
	}
	artifact := Artifact{}
	if strings.HasSuffix(artifactFilename, ".yaml") {
		file = yaml2Json(file, tag)
	}

	json.Unmarshal(file, &artifact)
	artifact.Data = file

	if artifact.Kind == "ReplicationController" {
		artifact.Url = fmt.Sprintf("%s/api/v1/namespaces/%s/replicationcontrollers", apiserver, namespace)
	}
	if artifact.Kind == "Service" {
		artifact.Url = fmt.Sprintf("%s/api/v1/namespaces/%s/services", apiserver, namespace)
	}
	if artifact.Kind == "Secret" {
		artifact.Url = fmt.Sprintf("%s/api/v1/namespaces/%s/secrets", apiserver, namespace)
	}
	if artifact.Kind == "Deployment" {
		artifact.Url = fmt.Sprintf("%s/apis/extensions/v1beta1/namespaces/%s/deployments", apiserver, namespace)
	}
	return artifact, err
}
