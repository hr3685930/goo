package gin

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

type GinHTTP struct {

}

func (*GinHTTP) HTTP(s *http.Server, route func(e *gin.Engine)) error {
    gin.SetMode(gin.ReleaseMode)
    g := gin.New()
    route(g)
    s.Handler = g
    return s.ListenAndServe()
}
