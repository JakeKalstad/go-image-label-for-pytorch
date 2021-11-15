package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type ImageSet struct {
	Name   string
	Images []string
}

func moveFile(source, dest string) {
	sourceFile, err := os.Open(source)
	if err != nil {
		panic(err)
	}
	defer sourceFile.Close()

	newFile, err := os.Create(dest)
	if err != nil {
		panic(err)
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, sourceFile)
	if err != nil {
		panic(err)
	}
}

type Config struct {
	DataFile           string
	ImagePrefix        string
	DefaultPrefix      string
	SecondaryPrefix    string
	SecondaryPredicate string
	IgnorePredicate    string
	Outfile            string
}

func main() {
	combineString := "%s/%s"
	configBytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	config := Config{}
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		panic(err)
	}
	if len(config.Outfile) == 0 {
		panic(errors.New("No output file name defined `classes`"))
	}
	if len(config.ImagePrefix) == 0 {
		panic(errors.New("No image directory defined `imagefolder`"))
	}
	if len(config.DataFile) == 0 {
		panic(errors.New("No data file defined `myinput.json`"))
	}
	if len(config.DefaultPrefix) == 0 {
		panic(errors.New("No default directory defined `my_data_set`"))
	}
	if len(config.SecondaryPredicate) > 0 && len(config.SecondaryPrefix) == 0 {
		panic(errors.New("No secondary directory defined with a predicate defined:" + config.SecondaryPredicate))
	}
	if _, err := os.Stat(config.DefaultPrefix); os.IsNotExist(err) {
		err = os.Mkdir(config.DefaultPrefix, os.FileMode(0777))
		if err != nil {
			panic(err)
		}
	}
	if len(config.SecondaryPredicate) > 0 {
		if _, err := os.Stat(config.SecondaryPrefix); os.IsNotExist(err) {
			err = os.Mkdir(config.SecondaryPrefix, os.FileMode(0777))
			if err != nil {
				panic(err)
			}
		}
	}

	data, err := ioutil.ReadFile(config.DataFile)
	if err != nil {
		panic(err)
	}
	imageSets := []ImageSet{}
	err = json.Unmarshal(data, &imageSets)
	if err != nil {
		panic(err)
	}
	classes := map[int]string{}

	for idx, b := range imageSets {
		classes[idx] = b.Name
		for _, img := range b.Images {
			if strings.Contains(img, config.IgnorePredicate) {
				continue
			}
			dir := fmt.Sprintf(combineString, config.DefaultPrefix, strconv.Itoa(idx))
			containsPredicate := strings.Contains(strings.ToLower(img), strings.ToLower(config.SecondaryPredicate))
			if len(config.SecondaryPredicate) > 0 && containsPredicate {
				dir = fmt.Sprintf(combineString, config.SecondaryPrefix, strconv.Itoa(idx))
			}
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				err = os.Mkdir(dir, os.FileMode(0777))
				if err != nil {
					panic(err)
				}
			}
			moveFile(fmt.Sprintf(combineString, config.ImagePrefix, img), fmt.Sprintf(combineString, dir, img))
		}
	}
	classBytes, err := json.Marshal(classes)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(config.Outfile+".json", classBytes, os.FileMode(0777))
	if err != nil {
		panic(err)
	}
}
