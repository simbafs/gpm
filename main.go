package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/simba-fs/gpm/proxy"
)

// proxyRoute represent a proxy route
// Example:
// {
//     From: "aurl.simba-fs.dev",
//     To:  "http://localhost:3000",
// }
type proxyRoute struct {
	From string `toml:"from"`
	To   string `toml:"to"`
}

type proxyRoutes []proxyRoute

func (p *proxyRoutes) String() string {
	return ""
}

func (p *proxyRoutes) Set(value string) error {
	fromTo := strings.SplitN(value, "--", 2)
	pr := proxyRoute{
		From: fromTo[0],
		To:   fromTo[1],
	}
	*p = append(*p, pr)
	return nil
}

func main() {
	cmdProxyConfig := proxyRoutes{}

	configPath := flag.String("file", "gpm.toml", "path to config file")
	flag.Var(&cmdProxyConfig, "proxy", "from->to, ex: gh.localhost:3000--https://github.com")
	listen := flag.String("listen", "0.0.0.0:3000", "listening address")
	flag.Parse()

	config := map[string]proxyRoute{}
	configFile, _ := os.ReadFile(*configPath)

	toml.Unmarshal(configFile, &config)

	fmt.Printf("Load proxies from %s:\n", *configPath)
	for _, value := range config {
		fmt.Printf("\t%s -> %s\n", value.From, value.To)
		proxy.Set(value.From, value.To)
	}
	fmt.Printf("Load proxies from cmd flags:\n")
	for _, value := range cmdProxyConfig {
		fmt.Printf("\t%s -> %s\n", value.From, value.To)
		proxy.Set(value.From, value.To)
	}

	fmt.Printf("Server start at %s\n", *listen)
	proxy.Listen(*listen)
}
