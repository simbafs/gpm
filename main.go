package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/simba-fs/gpm/proxy"
	"github.com/pelletier/go-toml/v2"
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

func main() {

	// config := map[string]proxyRoute{
	//     "yt": {"yt.localhost:3000", "https://youtube.com"},
	//     "gh": {"gh.localhost:3000", "https://github.com"},
	// }
	//
	// fmt.Printf("%#v\n", config)
	// configFile, err := toml.Marshal(config)
	// if err != nil {
	//     fmt.Printf("err: %v\n", err)
	//     return
	// }
	//
	// fmt.Printf("%#v\n", string(configFile))

	configPath := flag.String("file", "gpm.toml", "path to config file")
	listen := flag.String("listen", "0.0.0.0:3000", "listening address")
	flag.Parse()

	config := map[string]proxyRoute{}
	configFile, err := os.ReadFile(*configPath)
	if err != nil {
		fmt.Printf("Can't read config file %s\n", *configPath)
		return
	}

	toml.Unmarshal(configFile, &config)

	fmt.Printf("Load proxies:\n")
	for _, value := range config {
		fmt.Printf("\t%s -> %s\n", value.From, value.To)
		proxy.AddProxy(value.From, value.To)
	}

	proxy.Listen(*listen)
	fmt.Printf("Server start at %s\n", *listen)
}
