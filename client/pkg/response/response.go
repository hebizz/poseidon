package e

import (
  "strconv"

  "github.com/gin-gonic/gin"
  log "k8s.io/klog"
)

type Gin struct {
  C *gin.Context
}

type Response struct {
  Code    string      `json:"c"`
  Message string      `json:"msg"`
  Data    interface{} `json:"data"`
}

func (g *Gin) Response(httpCode int, code int, err error, d interface{}) {
  log.Info(GetMsg(code))
  if err != nil {
    log.Errorf("error msg:%s", err)
  }
  g.C.JSON(httpCode, Response{
    Code:    strconv.Itoa(code),
    Message: GetMsg(code),
    Data:    d,
  })
}

