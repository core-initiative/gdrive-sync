package main

import (
    "os"
    "io"
    "gopkg.in/yaml.v2"
)

type Config struct {
    Folder        string `yaml:"folder"`
    Schedule      string `yaml:"schedule"`
    TargetFolderID string `yaml:"targetFolderID"`
}

func LoadConfig(filename string) (*Config, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    data, err := io.ReadAll(file)
    if err != nil {
        return nil, err
    }

    var config Config
    err = yaml.Unmarshal(data, &config)
    if err != nil {
        return nil, err
    }

    return &config, nil
}
