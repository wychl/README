package main

import (
	"fmt"
	"net/http"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
)

var e *casbin.Enforcer

func init() {
	e = casbin.NewEnforcer("model.conf", "policy.csv")
}

func main() {
	router := gin.Default()
	router.Use(authorization)
	router.GET("/users/:id", func(c *gin.Context) {
		fmt.Println(c.Param("id"))
		c.String(http.StatusOK, "ok")
	})

	http.ListenAndServe(":8080", router)
}

func authorization(c *gin.Context) {
	path := c.Request.URL.Path
	user := c.Request.Header.Get("user")
	method := c.Request.Method

	sub := user
	obj := path
	act := method

	if e.Enforce(sub, obj, act) == true {
		c.Next()
		fmt.Println("allow")
	} else {
		c.Abort()
		fmt.Println("deny")
		c.String(http.StatusUnauthorized, "deny")
	}
}
