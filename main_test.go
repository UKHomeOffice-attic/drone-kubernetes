package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestIsJson(t *testing.T) {
	jsonFile, err := ioutil.ReadFile("example/test.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	assert.True(t, isJSON(string(jsonFile)), "Read and validate a json file")

}

func TestYaml2Json(t *testing.T) {
	yamlFile, err := ioutil.ReadFile("example/simple-test.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	file := yaml2Json(yamlFile, "")

	yamlFile, err = ioutil.ReadFile("example/variable-test.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	file = yaml2Json(yamlFile, "")
	assert.Nil(t, file, "Variables not resolved should return null")
	assert.False(t, isJSON(string(file)), "Is not a valid json")

	yamlFile, err = ioutil.ReadFile("example/variable-test.yaml")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	file = yaml2Json(yamlFile, "23")
	assert.True(t, isJSON(string(file)), "Variables are resolved, json is well formed")
}
