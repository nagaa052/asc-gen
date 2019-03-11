package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/ghodss/yaml"
	"github.com/nagaa052/asc-gen/pkg/asc"
	"github.com/pkg/errors"
)

var (
	templatePath      string = "./templates/output.md.tpl"
	defaultConfigPath string = "./config/config.yaml"
)

func main() {

	config, err := getConfig()
	if err != nil {
		log.Fatalf("%+v", err)
		os.Exit(1)
	}

	app, err := asc.New(config)
	if err != nil {
		log.Fatalf("%+v", err)
		os.Exit(1)
	}

	err = app.Run()
	if err != nil {
		log.Fatalf("%+v", err)
		os.Exit(1)
	}
}

func getConfig() (*asc.AssConfig, error) {
	cyml := os.Getenv("ASC_CONFIG")
	if cyml == "" {
		buf, err := ioutil.ReadFile(defaultConfigPath)
		if err != nil {
			return nil, errors.Wrapf(err, "failed read config file")
		}
		cyml = string(buf)
	}
	config := new(asc.AssConfig)
	err := yaml.Unmarshal([]byte(cyml), &config)
	if err != nil {
		return nil, errors.Wrapf(err, "failed Unmarshal config")
	}
	config.TemplatePath = templatePath
	return config, nil
}
