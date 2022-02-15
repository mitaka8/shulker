package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	_ "embed"

	"gitlab.com/mitaka8/shulker"
	_ "gitlab.com/mitaka8/shulker/console"
	_ "gitlab.com/mitaka8/shulker/discord"
	_ "gitlab.com/mitaka8/shulker/irc"
	_ "gitlab.com/mitaka8/shulker/matrix"
	_ "gitlab.com/mitaka8/shulker/minecraft"
	"gitlab.com/mitaka8/shulker/registry"
)

//go:embed config.sample.json
var defaultConfig []byte

func main() {
	configPath := flag.String("config", "config.json", "Path to the configuration file")
	flag.Parse()
	configBuffer, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Printf("Written default config file to: %v, %v\n", *configPath, err)
		fmt.Printf("%s\n", defaultConfig)
		ioutil.WriteFile(*configPath, defaultConfig, 0644)
		return
	}
	config := make([]registry.Config, 1)
	err = json.Unmarshal(configBuffer, &config)
	if err != nil {
		log.Fatalf("Cannot parse config file: %v", err)
		return
	}
	shulker.NewFromConfig(config)

}
