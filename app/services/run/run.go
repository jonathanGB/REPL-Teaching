package run

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)

func Run(port string) {
	runService := gin.Default()
	dockerPool := getPool(10)

	runService.POST("/run", func(c *gin.Context) {
		content := c.PostForm("content")
		extension := c.PostForm("extension")

		resCode := http.StatusOK
		res, err := dockerPool.Query([]byte(content), extension)
		if err != nil {
			resCode = http.StatusBadRequest
		}

		c.JSON(resCode, gin.H{
			"response": string(res),
		})
	})

	runService.Run(port)
}
