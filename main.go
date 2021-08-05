package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/op/go-logging"
	"github.com/pelletier/go-toml/v2"
	Config "github.com/simba-fs/gpm/config"
	Git "github.com/simba-fs/gpm/git"
	Host "github.com/simba-fs/gpm/host"
)

var formatter = logging.MustStringFormatter(
	`%{color}%{module} %{longfunc} ▶ %{level:.4s} %{color:reset} %{message}`,
)

var level = map[string]logging.Level{
	"critical": logging.CRITICAL,
	"error": logging.ERROR,
	"warning": logging.WARNING,
	"notice": logging.NOTICE,
	"info": logging.INFO,
	"debug": logging.DEBUG,
	"0": logging.CRITICAL,
	"1": logging.ERROR,
	"2": logging.WARNING,
	"3": logging.NOTICE,
	"4": logging.INFO,
	"5": logging.DEBUG,

}

func choice(choice ...string) string {
	for _, v := range choice {
		if v != "" {
			return v
		}
	}
	return ""
}

func main() {
	fmt.Println("test")
	logging.SetFormatter(formatter)
	log := logging.MustGetLogger("main")
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
	configPath := flag.String("file", "gpm.toml", "path to config file")
	flag.Var(&cmdHostConfig, "host", "from->to, ex: gh.localhost:3000--https://github.com")
	flag.Var(&cmdStaticConfig, "static", "repo^branch^name, ex: github.com/simba-fs/gpm^main^blog")
	address := flag.String("address", "", "listening address (default \"0.0.0.0:3000\")")
	logLevel := flag.String("logLevel", "info", "set log level.\nAvailable value: critical, error, warning, notice, info, debug, 0, 1, 2, 3, 4, 5")
	flag.Parse()

	logging.SetLevel(level[*logLevel], "")

	// read config file and parse
	config := Config.Config{}
	configFile, err := os.ReadFile(*configPath)
	if err == nil {
		log.Noticef("Read config file %s\n", *configPath)
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
	config.LogLevel = choice(*logLevel, config.LogLevel)
	logging.SetLevel(level[config.LogLevel], "")
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

	log.Debugf("config: %v\n", config)

	git := Git.Git{}
	go git.Init(&config)

	host := Host.Host{}
	host.Init(&config)
	host.Listen()
}
