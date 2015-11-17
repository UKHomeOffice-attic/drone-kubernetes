package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/drone/drone-plugin-go/plugin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var debug bool

type WebHook struct {
	Timestamp int64
	Images    []string
	Namespace string
	Source    string
	Target    string
	Url       string
	Token     string
}

type ReqEnvelope struct {
	Verb  string
	Token string
	Json  []byte
	Url   string
}

type Artifact struct {
	ApiVersion string
	Kind       string
	Data       []byte
	Metadata   struct {
		Name string
	}
	Url string
}

func zeroReplicas(artifact Artifact, token string) (bool, error) {

	json := `{"spec": {"replicas": 0}}`
	req := ReqEnvelope{
		Verb:  "PATCH",
		Token: token,
		Url:   fmt.Sprintf("%s/%s", artifact.Url, artifact.Metadata.Name),
		Json:  []byte(json),
	}
	res, err := doRequest(req)
	if err != nil {
		log.Panic("%s", err)
	}
	time.Sleep(time.Second * 5)
	return res, err
}

func deleteArtifact(artifact Artifact, token string) (bool, error) {
	if strings.Contains(artifact.Kind, "ReplicationController") {
		res, e := zeroReplicas(artifact, token)
		if e != nil {
			log.Panic("%s", e)
			os.Exit(1)
		}
		if res {
			if debug {
				log.Println("Replicas set to Zero")
			}
		}
	}

	url := fmt.Sprintf("%s/%s", artifact.Url, artifact.Metadata.Name)
	if debug {
		log.Println(url)
	}
	param := ReqEnvelope{
		Url:   url,
		Token: token,
		Verb:  "DELETE",
	}
	return doRequest(param)
}

func doRequest(param ReqEnvelope) (bool, error) {
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
			os.Exit(1)
		}
		log.Printf("%s\n", string(contents))
	}

	if err != nil {
		log.Panic("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()

		if response.StatusCode == 200 {
			return true, err
		}
	}
	return false, err
}

func existsArtifact(artifact Artifact, token string) (bool, error) {
	aUrl := fmt.Sprintf("%s/%s", artifact.Url, artifact.Metadata.Name)

	req := ReqEnvelope{
		Url:   aUrl,
		Token: token,
		Verb:  "GET",
	}
	return doRequest(req)

}

func createArtifact(artifact Artifact, token string) {
	deployments = append(deployments, artifact.Metadata.Name)
	param := ReqEnvelope{
		Url:   artifact.Url,
		Token: token,
		Json:  artifact.Data,
		Verb:  "POST",
	}
	doRequest(param)

}

func readArtifactFromFile(workspace string, artifactFile string, apiserver string, namespace string) (Artifact, error) {
	file, e := ioutil.ReadFile(workspace + "/" + artifactFile)
	// fmt.Println(string(file))
	if e != nil {
		log.Panic(e)
		os.Exit(1)
	}
	artifact := Artifact{}
	json.Unmarshal(file, &artifact)
	artifact.Data = file
	if artifact.Kind == "ReplicationController" {
		artifact.Url = fmt.Sprintf("%s/api/v1/namespaces/%s/replicationcontrollers", apiserver, namespace)
	}
	if artifact.Kind == "Service" {
		artifact.Url = fmt.Sprintf("%s/api/v1/namespaces/%s/services", apiserver, namespace)
	}

	return artifact, e
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func sendWebhook(wh *WebHook) {

	jwh, err := json.Marshal(wh)
	if err != nil {
		log.Panic(err)
		return
	}
	req := ReqEnvelope{
		Verb:  "POST",
		Token: wh.Token,
		Url:   wh.Url,
		Json:  []byte(jwh),
	}
	doRequest(req)
}

var deployments []string

func main() {
	var vargs = struct {
		ReplicationControllers []string `json:replicationcontrollers`
		Services               []string `json:services`
		ApiServer              string   `json:apiserver`
		Token                  string   `json:token`
		Namespace              string   `json:namespace`
		Debug                  string   `json:debug`
		Webhook                string   `json:webhook`
		Source                 string   `json:source`
		WebHookToken           string   `json:webhook_token`
	}{}

	workspace := plugin.Workspace{}
	plugin.Param("workspace", &workspace)
	plugin.Param("vargs", &vargs)
	plugin.Parse()

	// Iterate over rcs and svcs
	for _, rc := range vargs.ReplicationControllers {
		artifact, err := readArtifactFromFile(workspace.Path, rc, vargs.ApiServer, vargs.Namespace)
		if err != nil {
			log.Panic(err)
			return
		}
		if b, _ := existsArtifact(artifact, vargs.Token); b {
			deleteArtifact(artifact, vargs.Token)
			time.Sleep(time.Second * 5)
		}
		createArtifact(artifact, vargs.Token)
	}
	for _, rc := range vargs.Services {
		artifact, err := readArtifactFromFile(workspace.Path, rc, vargs.ApiServer, vargs.Namespace)
		if err != nil {
			log.Panic(err)
			return
		}
		createArtifact(artifact, vargs.Token)
	}
	wh := &WebHook{
		Timestamp: makeTimestamp(),
		Images:    deployments,
		Namespace: vargs.Namespace,
		Source:    vargs.Source,
		Target:    vargs.ApiServer,
		Url:       vargs.Webhook,
		Token:     vargs.WebHookToken,
	}
	sendWebhook(wh)
}
