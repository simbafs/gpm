package host

import (
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	Config "github.com/simba-fs/gpm/config"
)

type Host struct {
	ErrPage string
	Config  *Config.Config
}

func (h *Host) Init(c *Config.Config) {
	h.ErrPage = `<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="UTF-8">
		<title>Page Not Found</title>
		<style>
* {
	font-family: Roboto, Arial, sans-serif
}
		</style>
	</head>
	<body>
		<h1>Page Not Found</h1>
		<p>You may enter a wrong URL. Check your spell first. If it's OK, contact to the maintainer of this website.</p>
	</body>
</html>`
	h.Config = c
}

// Set sets a proxy route
func (h *Host) Set(from, to string) {
	h.Config.Host[from] = Config.Host{
		From: from,
		To:   to,
	}
}

// Remove removes a proxy route
func (h *Host) Remove(from string) {
	delete(h.Config.Host, from)
}

func (h *Host) routeProxy(c *gin.Context) {
	// get proxy route
	host, ok := Config.Host{}, false
	for _, v := range h.Config.Host {
		if v.From == c.Request.Host {
			host, ok = v, true
			break
		}
	}
	if !ok {
		c.Data(http.StatusBadRequest, "text/html", []byte(h.ErrPage))
		return
	}

	// parse url
	remote, err := url.Parse(host.To)
	if err != nil {
		c.Data(http.StatusBadRequest, "text/html", []byte(h.ErrPage))
		return
	}

	switch remote.Scheme {
	case "http", "https":
		h.routeProxyHttp(c, host)
	case "static":
		h.routeProxyStatic(c, host)
	}
}

// Listen starts a server on addr
func (h *Host) Listen() {
	gin.SetMode(gin.ReleaseMode)
	app := gin.Default()

	app.Any("/*proxyPath", h.routeProxy)
	fmt.Printf("Server start at %s\n", h.Config.Address)
	app.Run(h.Config.Address)
}

func (h *Host) SetConfig() {
	err := os.MkdirAll(h.Config.Storage, fs.ModeDir|fs.ModePerm)
	if err != nil {
		panic(err)
	}

	// set hosts
	for _, v := range h.Config.Host {
		h.Set(v.From, v.To)
	}
}
