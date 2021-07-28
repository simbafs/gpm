// Package gpm provide API to control reverse proxy
// example: 
// func main() {
//     gpm.AddProxy("localhost:3000", "https://github.com")
//     gpm.AddProxy("yt.localhost:3000", "https://youtube.com")
//     gpm.AddProxy("gh.localhost:3000", "https://github.com")
//     gpm.RemoveProxy("localhost:3000")
//     gpm.Listen(":3000")
// }
package gpm

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

// proxyRoute represent a proxy route
// Example:
// {
//     From: "aurl.simba-fs.dev",
//     TOo:  "http://localhost:3000",
// }
type proxyRoute struct {
	From string
	To   string
}

//	proxyRoutes store all proxy routes
var proxyRoutes = map[string]proxyRoute{}

// AddProxy add a proxy route
func AddProxy(from, to string) {
	proxyRoutes[from] = proxyRoute{from, to}
}

// AddProxy add a proxy route
func RemoveProxy(from string) {
	delete(proxyRoutes, from)
}

func routeProxy(c *gin.Context) {
	// get proxy route
	host, ok := proxyRoutes[c.Request.Host]
	if !ok {
		c.HTML(http.StatusBadRequest, "400.html", nil)
		return
	}

	// parse url
	remote, err := url.Parse(host.To)
	if err != nil {
		c.HTML(http.StatusBadRequest, "400.html", nil)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = c.Param("proxyPath")
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}

// Listen start a server on addr
func Listen(addr string) {
	app := gin.Default()
	app.LoadHTMLGlob("view/*")

	app.Any("/*proxyPath", routeProxy)
	app.Run(addr)
}

