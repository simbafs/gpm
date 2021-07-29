package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

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

// Listen starts a server on addr
func Listen(addr string) {
	gin.SetMode(gin.ReleaseMode)
	app := gin.Default()
	app.LoadHTMLGlob("./view/*")

	app.Any("/*proxyPath", routeProxy)
	app.Run(addr)
}

