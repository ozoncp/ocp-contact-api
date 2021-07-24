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

	readConfig := func() {
		file, err := os.Open(configPath)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("File %v was opened\n", configPath)
		defer func() {
			if err := file.Close(); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("File %v was closed\n", configPath)
			}
		}()
		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf(" Read content: %v\n", string(bytes))
		}
	}

	for i := 0; i < 10; i++{
		readConfig()
	}

	defer func(){
		if err := os.Remove(configPath); err == nil {
			fmt.Printf("File %v was removed\n", configPath)
		} else {
			fmt.Printf("File %v wasn't removed, error %v\n", configPath, err)
		}
	}()
}
