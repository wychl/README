package main

import (
	"fmt"

	"github.com/casbin/casbin"
)

func main() {
	e := casbin.NewEnforcer("model.conf", "policy.csv")
	sub := "alice" // the user that wants to access a resource.
	obj := "data1" // the resource that is going to be accessed.
	act := "read"  // the operation that the user performs on the resource.

	if e.Enforce(sub, obj, act) == true {
		fmt.Println("allow")
		// permit alice to read data1
	} else {
		fmt.Println("deny")
		// deny the request, show an error
	}

	roles := e.GetAllSubjects()

	fmt.Println(roles)
}
