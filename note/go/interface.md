# interface

## 操作 slice

```go
package main

import (
	"fmt"
	"reflect"
)

func main(){
    data:=[]int{1,2,3}
    err:=handler(data)
    if err!=nil{
        fmt.Println(err)
    }

}

func handler(in interface{})error{
    v := reflect.ValueOf(in)
    if v.Kind() != reflect.Slice {
       return fmt.Errorf("data type is %v",reflect.TypeOf(in))
    }

    num := v.Len()
    for i := 0; i < num; i++ {
    	fmt.Println(v.Index(i).Interface())
    }
    return nil
}
```

## reference

- [深入理解 Go Interface](http://legendtkl.com/2017/06/12/understanding-golang-interface/)
