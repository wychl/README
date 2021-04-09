## unsafe.Pointer

## 使用unsafe.Pointer将float64转化为uint64数据类型

```go
package main

import (
	"fmt"
	"unsafe"
)

func Float64ToUint64(f float64) *uint64 {
	return (*uint64)(unsafe.Pointer(&f))
}

func main() {
	f := 12.56
	fmt.Println(f)
}
```

**解释**

- `&f`取`f`的地址。
- `unsafe.Pointer(&f)`将`*float64`类型转换为`unsafe.Pointer`类型。
- `(*uint64)(unsafe.Pointer(&f))`将`unsafe.Pointer`类型转换为`*uint64`类型。


