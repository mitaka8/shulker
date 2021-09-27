package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"gitlab.com/mitaka8/shulker"
	_ "gitlab.com/mitaka8/shulker/console"
	_ "gitlab.com/mitaka8/shulker/discord"
	_ "gitlab.com/mitaka8/shulker/irc"
	_ "gitlab.com/mitaka8/shulker/minecraft"
	"gitlab.com/mitaka8/shulker/registry"
)

func main() {
	configPath := flag.String("config", "config.json", "Path to the configuration file")
	flag.Parse()
	configBuffer, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("Cannot read config file: %v", err)
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
