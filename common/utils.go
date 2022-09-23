package common

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

func RemoveStringFromSlice(in []string, remove string) []string {
	for i, v := range in {
		if v == remove {
			return append(in[:i], in[i+1:]...)
		}
	}
	return in
}
