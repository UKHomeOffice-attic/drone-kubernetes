package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func zeroReplicas(artifact Artifact, token string) (bool, error) {
	if debug {
		log.Println("zeroReplicas")
	}
	json := `{"spec": {"replicas": 0}}`
	req := ReqEnvelope{
		Verb:  "PATCH",
		Token: token,
		Url:   fmt.Sprintf("%s/%s", artifact.Url, artifact.Metadata.Name),
		Json:  []byte(json),
	}
	res, err := doRequest(req)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second * 5)
	return res, err
}

func deleteArtifact(artifact Artifact, token string) (bool, error) {
	if debug {
		log.Println("deleteArtifact")
	}
	if strings.Contains(artifact.Kind, "ReplicationController") {
		res, err := zeroReplicas(artifact, token)
		if err != nil {
			log.Fatal(err)
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

func existsArtifact(artifact Artifact, token string) (bool, error) {
	aUrl := fmt.Sprintf("%s/%s", artifact.Url, artifact.Metadata.Name)
	if debug {
		log.Println("existsArtifact " + aUrl)
	}
	req := ReqEnvelope{
		Url:   aUrl,
		Token: token,
		Verb:  "GET",
	}
	return doRequest(req)

}

func createArtifact(artifact Artifact, token string) {
	if debug {
		log.Println("createArtifact ")
	}
	deployments = append(deployments, artifact.Metadata.Name)
	param := ReqEnvelope{
		Url:   artifact.Url,
		Token: token,
		Json:  artifact.Data,
		Verb:  "POST",
	}
	doRequest(param)
}

func updateArtifact(artifact Artifact, token string) {
	if debug {
		log.Println("updateArtifact ")
	}
	artifact.Url = fmt.Sprintf("%s/%s", artifact.Url, artifact.Metadata.Name)
	param := ReqEnvelope{
		Url:   artifact.Url,
		Token: token,
		Json:  artifact.Data,
		Verb:  "PUT",
	}
	doRequest(param)
}
