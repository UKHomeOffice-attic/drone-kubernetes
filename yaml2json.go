package main

import (
	"bytes"
	"encoding/json"
	"gopkg.in/yaml.v2"
	"log"
	"strconv"
	"text/template"
)

func transformData(in interface{}) (out interface{}, err error) {
	switch in.(type) {
	case map[interface{}]interface{}:
		o := make(map[string]interface{})
		for k, v := range in.(map[interface{}]interface{}) {
			sk := ""
			switch k.(type) {
			case string:
				sk = k.(string)
			case int:
				sk = strconv.Itoa(k.(int))
			default:
				log.Fatal(err)
			}
			v, err = transformData(v)
			if err != nil {
				return nil, err
			}
			o[sk] = v
		}
		return o, nil
	case []interface{}:
		in1 := in.([]interface{})
		len1 := len(in1)
		o := make([]interface{}, len1)
		for i := 0; i < len1; i++ {
			o[i], err = transformData(in1[i])
			if err != nil {
				return nil, err
			}
		}
		return o, nil
	default:
		return in, nil
	}
	return in, nil
}

func yaml2Json(file []byte, tag string) []byte {

	var data interface{}

	t := template.New("artifact")
	t, _ = t.Parse(string(file))
	yamdata := new(bytes.Buffer)
	ddata := DroneData{TAG: tag}
	t.Execute(yamdata, ddata)
	err := yaml.Unmarshal(yamdata.Bytes(), &data)

	data, err = transformData(data)
	if err != nil {
		log.Fatal(err)
	}
	file, err = json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	if string(file) != "null" {
		return file
	} else {
		return nil
	}
}

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil

}
