package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

const errPage = `<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="UTF-8">
		<title>No Proxy Route Found</title>
		<style>
* {
	font-family: Roboto, Arial, sans-serif
}
		</style>
	</head>
	<body>
		<h1>No Proxy Route Found</h1>
		<p>You may enter a wrong URL. Check your spell first. If it's OK, contact to the maintainer of this website.</p>
	</body>
</html>`

type proxyRoute struct {
	From string `toml:"from"`
	To   string `toml:"to"`
}

//	proxyRoutes store all proxy routes
var proxyRoutes = map[string]proxyRoute{}

// Set sets a proxy route
func Set(from, to string) {
	proxyRoutes[from] = proxyRoute{from, to}
}

// Remove removes a proxy route
func Remove(from string) {
	delete(proxyRoutes, from)
}

func routeProxy(c *gin.Context) {
	// get proxy route
	host, ok := proxyRoutes[c.Request.Host]
	if !ok {
		c.Data(http.StatusBadRequest, "text/html", []byte(errPage))
		return
	}

	// parse url
	remote, err := url.Parse(host.To)
	if err != nil {
		c.Data(http.StatusBadRequest, "text/html", []byte(errPage))
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

// Listen starts a server on addr
func Listen(addr string) {
	gin.SetMode(gin.ReleaseMode)
	app := gin.Default()

	app.Any("/*proxyPath", routeProxy)
	app.Run(addr)
}

