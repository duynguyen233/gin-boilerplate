package controller

import (
	"net/http"
	"wine-be/internal/model"
	"wine-be/pkg/logger"

	"github.com/gin-gonic/gin"
)

type DefaultRoutes struct {
	l logger.Interface
}

func NewDefaultRoutes(handler *gin.RouterGroup, l logger.Interface) *DefaultRoutes {
	r := &DefaultRoutes{l}
	handler.GET("/ping", r.ping)
	return r
}

// @Summary     Ping default
// @Description Ping default
// @ID          ping
// @Tags  	    Default
// @Accept      json
// @Produce     json
// @Success     200 {object} model.Response[string]
// @Failure     500 {object} model.ErrorResponse
// @Router      /ping [get]
func (r *DefaultRoutes) ping(c *gin.Context) {
	c.JSON(http.StatusOK, model.Response[string]{
		Message: "Ping successfully",
		Data:    "pong",
	})
}
