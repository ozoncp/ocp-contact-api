package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"database"`
	SSL      string `yaml:"ssl"`
	Driver   string `yaml:"driver"`
}

type Grpc struct {
	Address string `yaml:"address"`
}

type Request struct {
	BatchSize int `yaml:"batchsize"`
}

type Project struct {
	Name    string `yaml:"name"`
	Author  string `yaml:"author"`
	Version string `yaml:"version"`
}

type Kafka struct {
	Topic   string   `yaml:"topic"`
	Brokers []string `yaml:"brokers"`
}

type Prometheus struct {
	Uri  string `yaml:"uri"`
	Port string `yaml:"port"`
}

type Config struct {
	Project  Project  `yaml:"project"`
	Grpc     Grpc     `yaml:"grpc"`
	Database Database `yaml:"database"`
	Request  Request  `yaml:"request"`
	Kafka    Kafka    `yaml:"kafka"`
	Prometheus Prometheus `yaml:"prometheus"`
}

func Read(path string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}
