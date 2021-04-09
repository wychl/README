package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/", func(c *gin.Context) {
		var resp Response
		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			resp.Code = 1
			resp.ErrorMessage = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		defer c.Request.Body.Close()

		fmt.Println(string(data))
	})
	r.Run(":5001")
}

// Response ...
type Response struct {
	Code         int    `json:"code"`
	ErrorMessage string `json:"error_message"`
}
