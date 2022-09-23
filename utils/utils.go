package utils

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

func LoadFromFile(file string, target interface{}) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, target)
}

func StringsContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
