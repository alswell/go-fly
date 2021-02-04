package fmtio

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

func LoadJSON(filename string, v interface{}) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

func DumpJSON(filename string, v interface{}) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, 0644)
}

func PrintJSON(v interface{}) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func LoadYAML(filename string, v interface{}) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, v)
}

func DumpYAML(filename string, v interface{}) error {
	b, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, 0644)
}

func PrintYAML(v interface{}) error {
	b, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}
