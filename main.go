package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/pelletier/go-toml/v2"
	Config "github.com/simba-fs/gpm/config"
	Host "github.com/simba-fs/gpm/host"
)

func choice(choice ...string) string {
	for _, v := range choice {
		if v != "" {
			return v
		}
	}
	return ""
}

func main() {
	// parse cmd flags
	cmdHostConfig := Config.HostSlice{}
	cmdStaticConfig := Config.StaticSlice{}

	storagePath := flag.String("storage", "", "directory to store files such as static files (default \"./storage\")")
	configPath := flag.String("file", "gpm.toml", "path to config file")
	flag.Var(&cmdHostConfig, "host", "from->to, ex: gh.localhost:3000--https://github.com")
	flag.Var(&cmdStaticConfig, "static", "repo^branch^name, ex: github.com/simba-fs/gpm^main^blog")
	address := flag.String("address", "", "listening address (default \"0.0.0.0:3000\")")
	flag.Parse()

	// read config file and parse
	config := Config.Config{}
	configFile, err := os.ReadFile(*configPath)
	if err == nil {
		fmt.Printf("Read config file %s\n", *configPath)
		toml.Unmarshal(configFile, &config)
	}

	if err != nil {
		panic(err)
	}

	// merge config file and cmd flags
	// address
	config.Address = choice(*address, config.Address, "0.0.0.0:3000")
	// storage
	cwd, _ := os.Getwd()
	config.Storage = choice(*storagePath, "./storage")
	if !path.IsAbs(config.Storage) {
		config.Storage = path.Join(cwd, config.Storage)
	}
	// cmdHostConfig
	for _, v := range cmdHostConfig {
		config.Host[v.From] = v
	}
	// cmdStaticConfig
	for _, v := range cmdStaticConfig {
		config.Static[v.Name] = v
	}
	for k, v := range config.Static {
		config.Static[k] = Config.Static{
			Name: k,
			Repo: v.Repo,
			Branch: v.Branch,
			Path: path.Join(config.Storage, k),
		}
	}

	fmt.Printf("config: %v\n", config)

	host := Host.Host{}
	host.Init(&config)
	host.Listen()
}
