package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"regexp"
)

func Unique(slice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func Contains(slice []string, t string) bool {
	for _, v := range slice {
		if v == t {
			return true
		}
	}

	return false
}

func GetExternalLinks(body string) ([]string, error) {
	re := regexp.MustCompile(URL_REGEX)
	matchedLinks := re.FindAllStringSubmatch(body, -1)

	var result []string
	for _, link := range matchedLinks {
		u, err := url.ParseRequestURI(link[1])
		if err != nil {
			continue
		}
		if u.Hostname() == "" {
			continue
		}
		l := fmt.Sprintf("%s://%s", "http", u.Hostname())
		result = append(result, l)
	}

	result = Unique(result)

	return result, nil
}

func SaveToDiskAsJson(val interface{}) {
	bytes, err := json.MarshalIndent(val, "", "\t")
	if err != nil {
		return
	}
	err = ioutil.WriteFile("fingerprints.json", bytes, 0644)
	if err != nil {
		return
	}

	fmt.Println("saved fingerprints results to json")
}
