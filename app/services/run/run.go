package run

import (
	"fmt"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)

func Run(port string) {
	runService := gin.Default()
	//dockerPool := getPool(10)

	runService.POST("/run", func(c *gin.Context) {
		content := c.PostForm("content")
		extension := c.PostForm("extension")

		resCode := http.StatusOK
		res, err := runQuery([]byte(content), extension)
		if err != nil {
			fmt.Printf("Erreur lors du run: %v\n", err)
			resCode = http.StatusBadRequest
		}

		c.String(resCode, string(res))
	})

	runService.Run(port)
}
