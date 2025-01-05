package internal

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"order-service/global"
	"order-service/model"
	"time"
)

func SaveOrUpdate(c *gin.Context) {
	var order model.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, "Invalid Data")
		return
	}
	order.CreatedAt = time.Now().Unix()
	global.GLOBAL_DB.Create(&order)
	c.JSON(http.StatusOK, order)
}

func GetById(c *gin.Context) {
	var order model.Order
	id := c.Param("id")
	if err := global.GLOBAL_DB.Where("id=?", id).First(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, "Order not found")
		return
	}
	c.JSON(http.StatusOK, order)
}
