package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"user-service/remote"
)

func OrderList(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}
	orders := remote.Call("order")
	c.JSON(http.StatusOK, gin.H{"userId": userId, "orderList": orders})
}
