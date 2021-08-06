package main

import (
	"flag"
	"os"
	"path"

	"github.com/pelletier/go-toml/v2"
	Config "github.com/simba-fs/gpm/config"
	Git "github.com/simba-fs/gpm/git"
	Host "github.com/simba-fs/gpm/host"
	Log "github.com/simba-fs/gpm/log"
)

var config = Config.Config{}

func choice(choice ...string) string {
	for _, v := range choice {
		if v != "" {
			return v
		}
	}
	return ""
}

func init() {
	log := Log.NewLog("main")
	// log.Debugf("debug %s", "test")
	// log.Info("info")
	// log.Notice("notice")
	// log.Warning("warning")
	// log.Error("err")
	// log.Critical("crit")

	// parse cmd flags
	cmdHostConfig := Config.HostSlice{}
	cmdStaticConfig := Config.StaticSlice{}

	storagePath := flag.String("storage", "", "directory to store files such as static files (default \"./storage\")")
	logPath := flag.String("log", "", "directory to store log files (default \"./log\")")
	configPath := flag.String("file", "gpm.toml", "path to config file")
	flag.Var(&cmdHostConfig, "host", "from->to, ex: gh.localhost:3000--https://github.com")
	flag.Var(&cmdStaticConfig, "static", "repo^branch^name, ex: github.com/simba-fs/gpm^main^blog")
	address := flag.String("address", "", "listening address (default \"0.0.0.0:3000\")")
	logLevel := flag.String("logLevel", "info", "set log level.\nAvailable value: critical, error, warning, notice, info, debug, 0, 1, 2, 3, 4, 5")
	flag.Parse()

	// read config file and parse
	configFile, err := os.ReadFile(*configPath)
	if err == nil {
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
	config.Storage = choice(*storagePath, config.Storage, "./storage")
	if !path.IsAbs(config.Storage) {
		config.Storage = path.Join(cwd, config.Storage)
	}
	config.LogLevel = choice(*logLevel, config.LogLevel)
	// log file
	config.Log = choice(*logPath, config.Log)
	if config.Log != "" && !path.IsAbs(config.Log) {
		config.Log = path.Join(cwd, config.Log)
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
			Name:   k,
			Repo:   v.Repo,
			Branch: v.Branch,
			Path:   path.Join(config.Storage, k),
		}
	}

	Log.Init(&config)

	log.Debugf("config: %v\n", config)
}

func main() {
	git := Git.Git{}
	go git.Init(&config)

	host := Host.Host{}
	host.Init(&config)
	host.Listen()
}
