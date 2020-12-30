package config

import (
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Run struct {
		EnableCIMode bool `yaml:"enable_ci_mode"`
	} `yaml:"run"`
	StackWithDriftDetection    []string `yaml:"stack_with_drift_detection"`
	StackWithoutDriftDetection []string `yaml:"stack_without_drift_detection"`
	StackTemplateMaps          []struct {
		StackName    string `yaml:"stack_name,omitempty"`
		TemplateName string `yaml:"template_name,omitempty"`
	} `yaml:"stack_template_maps,omitempty"`
}

func load(file string) (Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return Config{}, nil
	}
	defer f.Close()
	return loadReader(f)
}

func loadReader(fd io.Reader) (Config, error) {
	data, err := ioutil.ReadAll(fd)
	if err != nil {
		return Config{}, err
	}
	cfg := Config{}
	err = yaml.UnmarshalStrict(data, &cfg)
	return cfg, err
}

func Load(path string) (Config, error) {
	if path != "" {
		return load(path)
	}
	for _, f := range [4]string{
		".cfndiff.yml",
		".cfndiff.yaml",
		"cfndiff.yml",
		"cfndiff.yaml",
	} {
		cfg, err := load(f)
		if err != nil && os.IsNotExist(err) {
			continue
		}
		return cfg, err
	}
	return Config{}, nil
}
