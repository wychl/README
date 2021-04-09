## go数据库自定义类型

```go
package main

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Address struct {
	Name string
}

func (a Address) Value() (driver.Value, error) {
	data, err := json.Marshal(a)

	return string(data), err
}

func (a *Address) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), a)
}

type Customer struct {
	CustomerID int64
	Address    Address
}

func main() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db = db.Debug()

	db.AutoMigrate(&Customer{})

	c := &Customer{
		CustomerID: 10,
		Address: Address{
			Name: "shanghai",
		},
    }
    
    //插入数据
	err = db.Model(c).Create(c).Error
	if err != nil {
		fmt.Println(err.Error())
		return
	}

    //更新数据
	// err = db.Model(&Customer{}).Where(c).UpdateColumns(map[string]interface{}{"address": Address{
	// 	Name: "rrr",
	// }}).Error
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }

    //查询数据
	var result Customer
	err = db.Model(&Customer{}).Last(&result).Error
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("%#v\n", result)

}


```