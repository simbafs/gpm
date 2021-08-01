package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/simba-fs/gpm/proxy"
)

type host struct {
	From string `toml:"from"`
	To   string `toml:"to"`
}

type static struct {
	Name   string `toml:"name"`
	Repo   string `toml:"repo"`
	Branch string `toml:"branch"`
}

type config struct {
	Address string            `toml:"address"`
	Host    map[string]host   `toml:"host"`
	Static  map[string]static `toml:"static"`
}

type staticSlice []static

func (s *staticSlice) String() string {
	return ""
}

func (s *staticSlice) Set(value string) error {
	repoBranch := strings.SplitN(value, "^", 3)
	*s = append(*s, static{repoBranch[0], repoBranch[1], repoBranch[2]})

	return nil
}

type hostSlice []host

func (p *hostSlice) String() string {
	return ""
}

func (p *hostSlice) Set(value string) error {
	fromTo := strings.SplitN(value, "--", 2)
	*p = append(*p, host{fromTo[0], fromTo[1]})
	return nil
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
	// parse cmd flags
	cmdHostConfig := hostSlice{}
	cmdStaticConfig := staticSlice{}

	configPath := flag.String("file", "gpm.toml", "path to config file")
	flag.Var(&cmdHostConfig, "host", "from->to, ex: gh.localhost:3000--https://github.com")
	flag.Var(&cmdStaticConfig, "static", "repo^branch^name, ex: github.com/simba-fs/gpm^main^blog")
	address := flag.String("address", "", "listening address (default \"0.0.0.0:3000\")")
	flag.Parse()

	// read config file and parse
	config := config{}
	configFile, err := os.ReadFile(*configPath)
	if err == nil {
		fmt.Printf("Read config file %s\n", *configPath)
		toml.Unmarshal(configFile, &config)
	}

	if err != nil {
		panic(err)
	}

	// merge config file and cmd flags
	// cmdHostConfig
	for _, v := range cmdHostConfig {
		config.Host[v.From] = v
	}
	// cmdStaticConfig
	for _, v := range cmdStaticConfig {
		config.Static[v.Name] = v
	}
	// address
	config.Address = choice(*address, config.Address, "0.0.0.0:3000")

	fmt.Printf("config: %v\n", config)
	fmt.Printf("Server start at %s\n", config.Address)
	proxy.Listen(config.Address)
}
