package main

import (
	"coding-test/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/person", handler.SelectDataAll)
	r.GET("/person/:id", handler.SelectData)
	r.POST("/person", handler.InsertData)
	r.PUT("/person/:id", handler.UpdateData)
	r.DELETE("/person/:id", handler.DeleteData)

	r.Run("localhost:3000")
}

//post는 body로
//request type, response type은 보통 구조체에서 정의해서 씀
