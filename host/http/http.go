package http

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	Config "github.com/simba-fs/gpm/config"
)

func Route(c *gin.Context) {
	h, ok := c.Get("host")
	if !ok {
		return
	}
	host := h.(*Config.Host)

	ErrPage, _ := c.Get("ErrPage")

	// parse url
	remote, err := url.Parse(host.To)
	if err != nil {
		c.Data(http.StatusBadRequest, "text/html", ErrPage.([]byte))
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
