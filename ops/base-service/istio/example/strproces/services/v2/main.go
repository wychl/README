package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/", func(c *gin.Context) {

		req := struct {
			Input string
		}{}

		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusOK, err.Error())
			return
		}
		defer c.Request.Body.Close()

		err = json.Unmarshal(data, &req)
		if err != nil {
			c.String(http.StatusOK, err.Error())
			return
		}

		c.String(http.StatusOK, strings.ToLower(req.Input))

	})

	http.ListenAndServe(":8080", router)
}
