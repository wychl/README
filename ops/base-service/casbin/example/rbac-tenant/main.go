package tenants

import (
	"fmt"
	"log"

	"github.com/casbin/casbin"
)

// TestTenants test tenants
func TestTenants() {
	e := casbin.NewEnforcer("model.conf", "tenants/policy.csv")

	fmt.Printf("RBAC TENANTS test start\n") // output for debug

	// superAdmin
	if e.Enforce("superAdmin", "gy", "project", "read") {
		log.Println("superAdmin can read project in gy")
	} else {
		log.Fatal("ERROR: superAdmin can not read project in gy")
	}

	if e.Enforce("superAdmin", "gy", "project", "write") {
		log.Println("superAdmin can write project in gy")
	} else {
		log.Fatal("ERROR: superAdmin can not write project in gy")
	}

	if e.Enforce("superAdmin", "jn", "project", "read") {
		log.Println("superAdmin can read project in jn")
	} else {
		log.Fatal("ERROR: superAdmin can not read project in jn")
	}

	if e.Enforce("superAdmin", "jn", "project", "write") {
		log.Println("superAdmin can write project in jn")
	} else {
		log.Fatal("ERROR: superAdmin can not write project in jn")
	}

	// admin
	if e.Enforce("quyuan", "gy", "project", "read") {
		log.Println("quyuan can read project in gy")
	} else {
		log.Fatal("ERROR: quyuan can not read project in gy")
	}

	if e.Enforce("quyuan", "gy", "project", "write") {
		log.Println("quyuan can write project in gy")
	} else {
		log.Fatal("ERROR: quyuan can not write project in gy")
	}

	if e.Enforce("quyuan", "jn", "project", "read") {
		log.Fatal("ERROR: quyuan can read project in jn")
	} else {
		log.Println("quyuan can not read project in jn")
	}

	if e.Enforce("quyuan", "jn", "project", "write") {
		log.Fatal("ERROR: quyuan can write project in jn")
	} else {
		log.Println("quyuan can not write project in jn")
	}

	if e.Enforce("quyuan", "gy", "asse", "read") {
		log.Fatal("ERROR: quyuan can read asse in gy")
	} else {
		log.Println("quyuan can not read asse in gy")
	}

	if e.Enforce("quyuan", "gy", "asse", "write") {
		log.Fatal("ERROR: quyuan can write asse in gy")
	} else {
		log.Println("quyuan can not write asse in gy")
	}

	if e.Enforce("quyuan", "jn", "asse", "read") {
		log.Println("quyuan can read asse in jn")
	} else {
		log.Fatal("ERROR: quyuan can not read asse in jn")
	}

	if e.Enforce("quyuan", "jn", "asse", "write") {
		log.Println("quyuan can write asse in jn")
	} else {
		log.Fatal("ERROR: quyuan can not write asse in jn")
	}

	// wenyin
	if e.Enforce("wenyin", "gy", "asse", "write") {
		log.Println("wenyin can write asse in gy")
	} else {
		log.Fatal("ERROR: wenyin can not write asse in gy")
	}

	if e.Enforce("wenyin", "jn", "asse", "write") {
		log.Fatal("ERROR: wenyin can write asse in jn")
	} else {
		log.Println("wenyin can not write asse in jn")
	}

	// shangshang
	if e.Enforce("shangshang", "jn", "project", "write") {
		log.Println("shangshang can write project in jn")
	} else {
		log.Fatal("ERROR: shangshang can not write project in jn")
	}

	if e.Enforce("shangshang", "gy", "project", "write") {
		log.Fatal("ERROR: shangshang can write project in gy")
	} else {
		log.Println("shangshang can not write project in gy")
	}
}
