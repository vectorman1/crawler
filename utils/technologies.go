package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type ScannerData struct {
	Categories   map[string]ImportedCategory   `json:"categories,omitempty"`
	Technologies map[string]ImportedTechnology `json:"technologies,omitempty"`
}

type ImportedCategory struct {
	Name     string `json:"name,omitempty"`
	Priority int    `json:"priority,omitempty"`
}

type ImportedTechnology struct {
	Cats        []int             `json:"cats,omitempty"`
	Description string            `json:"description,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	Html        []string          `json:"html,omitempty"` //problematic
	Icon        string            `json:"icon,omitempty"`
	Implies     []string          `json:"implies,omitempty"` //problematic
	Scripts     []string          `json:"scripts,omitempty"` //problematic - replace for import "scripts":\s"(.+)" with "scripts": [ "\1" ]
	Cookies     map[string]string `json:"cookies,omitempty"`
	Website     string            `json:"website,omitempty"`
	CertIssuer  string            `json:"certIssuer,omitempty"`
	// TODO: Add remaining fields
}

var _technologies *ScannerData

func ReadUrls(loc string) ([]string, error) {
	file, err := os.Open(loc)

	if err != nil {
		log.Fatalf("Failed to read urls %v", err)
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}

func NewScannerData() *ScannerData {
	wd, _ := os.Getwd()
	file, err := ioutil.ReadFile(fmt.Sprintf("%s\\technologies.json", wd))

	if err != nil {
		log.Fatalln(err)
		return nil
	}

	var technologies ScannerData

	err = json.Unmarshal(file, &technologies)

	if err != nil {
		log.Fatalln(err)
		return nil
	}

	_technologies = &technologies

	return _technologies
}

func GetCategories(t ImportedTechnology) []ImportedCategory {
	var result []ImportedCategory
	for _, tc := range t.Cats {
		for k, _ := range _technologies.Categories {
			convertedKey := fmt.Sprintf("%d", tc)

			if convertedKey == k {
				result = append(result, _technologies.Categories[k])
			}
		}
	}
	return result
}
