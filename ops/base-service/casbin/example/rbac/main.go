package rbac

import (
	"fmt"
	"log"

	"github.com/casbin/casbin"
)

func TestRBAC() {
	e := casbin.NewEnforcer("model.conf", "tenants/policy.csv")

	fmt.Printf("RBAC test start\n") // output for debug

	// superAdmin
	if e.Enforce("superAdmin", "project", "read") {
		log.Println("superAdmin can read project")
	} else {
		log.Fatal("ERROR: superAdmin can not read project")
	}

	if e.Enforce("superAdmin", "project", "write") {
		log.Println("superAdmin can write project")
	} else {
		log.Fatal("ERROR: superAdmin can not write project")
	}

	// admin
	if e.Enforce("quyuan", "project", "read") {
		log.Println("quyuan can read project")
	} else {
		log.Fatal("ERROR: quyuan can not read project")
	}

	if e.Enforce("quyuan", "project", "write") {
		log.Println("quyuan can write project")
	} else {
		log.Fatal("ERROR: quyuan can not write project")
	}

	if e.Enforce("quyuan", "asse", "read") {
		log.Println("quyuan can read asse")
	} else {
		log.Fatal("ERROR: quyuan can not read asse")
	}

	if e.Enforce("quyuan", "asse", "write") {
		log.Println("quyuan can write asse")
	} else {
		log.Fatal("ERROR: quyuan can not write asse")
	}

	// zhuangjia
	if e.Enforce("wenyin", "project", "read") {
		log.Fatal("ERROR: wenyin can read project")
	} else {
		log.Println("wenyin can not read project")
	}

	if e.Enforce("wenyin", "project", "write") {
		log.Println("wenyin can write project")
	} else {
		log.Fatal("ERROR: wenyin can not write project")
	}

	if e.Enforce("wenyin", "asse", "read") {
		log.Fatal("ERROR: wenyin can read asse")
	} else {
		log.Println("wenyin can not read asse")
	}

	if e.Enforce("wenyin", "asse", "write") {
		log.Println("wenyin can write asse")
	} else {
		log.Fatal("ERROR: wenyin can not write asse")
	}

	// shangshang
	if e.Enforce("shangshang", "project", "read") {
		log.Println("shangshang can read project")
	} else {
		log.Fatal("ERROR: shangshang can not read project")
	}

	if e.Enforce("shangshang", "project", "write") {
		log.Fatal("ERROR: shangshang can write project")
	} else {
		log.Println("shangshang can not write project")
	}

	if e.Enforce("shangshang", "asse", "read") {
		log.Println("shangshang can read asse")
	} else {
		log.Fatal("ERROR: shangshang can not read asse")
	}

	if e.Enforce("shangshang", "asse", "write") {
		log.Fatal("ERROR: shangshang can write asse")
	} else {
		log.Println("shangshang can not write asse")
	}
}
