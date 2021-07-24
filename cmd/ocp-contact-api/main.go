package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	fmt.Println("This is an Ozon Contact API")
	const configPath string = "config.cfg"

	fmt.Printf("Creating a file %v\n", configPath)
	file, err := os.Create(configPath)
	if err != nil {
		fmt.Printf("Cannot create file by path %v, error %v\n", configPath, err)
		return
	}

	file.WriteString("config file content")
	file.Close()

	readConfig := func() (string, error){
		file, err := os.Open(configPath)
		if err != nil {
			return "", err
		}
		fmt.Printf("File %v was opened\n", configPath)

		defer file.Close()

		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			return "", err
		}

		return string(bytes), nil
	}

	for i := 0; i < 10; i++{
		config, err := readConfig()
		if err != nil {
			fmt.Printf("Cannot read config by path %v, error %v\n", configPath, err)
			return
		}
		fmt.Printf(" Read content: %v\n", config)
	}

	defer os.Remove(configPath)
}
