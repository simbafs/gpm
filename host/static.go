package host

import (
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	Config "github.com/simba-fs/gpm/config"
	// Log "github.com/simba-fs/gpm/log"
)

var fileservers = map[string]http.Handler{}
// var log = Log.NewLog("host/static")

func (h *Host) routeProxyStatic(c *gin.Context, host Config.Host) {
	// block access to .git
	url, ok := strings.Split(c.Request.URL.Path, "/"), true
	for _, v := range url {
		if v == ".git" {
			ok = false
		}
	}
	if !ok {
		c.Redirect(http.StatusFound, "/")
	}

	hostName := strings.SplitN(host.To, "://", 2)[1]

	filePath := path.Join(h.Config.Storage, hostName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.Data(http.StatusBadRequest, "text/html", []byte(h.ErrPage))
	}

	fileserver, ok := fileservers[hostName]
	if ok {
		fileserver.ServeHTTP(c.Writer, c.Request)
	} else {
		filesystem := os.DirFS(filePath)
		fileserver = http.FileServer(http.FS(filesystem))
		fileservers[hostName] = fileserver
		fileserver.ServeHTTP(c.Writer, c.Request)
	}
	c.Abort()
}
