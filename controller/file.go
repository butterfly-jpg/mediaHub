package controller

import (
	"mediahub/pkg/storage"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	os storage.Storage
}

func NewController(os storage.Storage) *Controller {
	return &Controller{os}
}

func (c *Controller) Upload(ctx *gin.Context) {

}
