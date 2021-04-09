package main

import (
	"fmt"
)

var appConfig *Config

func main() {
	fmt.Println(
		appConfig.ToApplication().GenerateYaml(),
	)

}
