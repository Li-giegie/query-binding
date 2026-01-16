# query-binding
是一个增强Gin框架中的QueryBinding自定义切片类型的实现
不分类型，只要实现了UnmarshalParam接口，一定会被回调，这一点与官方的实现差别很大

```go
// Gin 内置接口
type UnmarshalParam interface {
    UnmarshalParam(val string) error
}
// 改为 支持自定义切片类型
type UnmarshalParam interface {
	UnmarshalParam(vals []string) error
}
```

## 获取
```go
go get github.com/Li-giegie/query-binding
```

## 使用
```go
package main

import (
	"errors"
	"fmt"
	"github.com/Li-giegie/query-binding"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"time"
)

type TimeX struct {
	Time time.Time
}

func (t *TimeX) UnmarshalParam(vals []string) error {
	fmt.Println("UnmarshalParam", vals)
	for _, s := range []string{time.RFC3339, time.DateTime, time.DateOnly} {
		d, err := time.Parse(s, vals[0])
		if err == nil {
			t.Time = d
			return nil
		}
	}
	return errors.New("time format error")
}

type Param struct {
	Time  TimeX   `form:"time"`
	Times []TimeX `form:"times"`
}

func main() {
	// 把Gin默认的验证器 binding.Validator 注册到 query_binding.Validator
	query_binding.Validator = binding.Validator
	// 把扩展版本 query_binding.Default 注册到 binding.Query
	binding.Query = &query_binding.Default

	eng := gin.Default()
	eng.GET("/", func(c *gin.Context) {
		var param Param
		err := c.ShouldBindQuery(&param)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, param)
	})
	fmt.Println(eng.Run(":8080"))
}

// curl --location --request GET 'http://localhost:8080/?time=2025-01-01&times=2025-02-01&times=2025-03-01%2015:14:01' \
//--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
//--header 'Accept: */*' \
//--header 'Host: localhost:8080' \
//--header 'Connection: keep-alive'

// out
// {
//    "Time": {
//        "Time": "2025-01-01T00:00:00Z"
//    },
//    "Times": [
//        {
//            "Time": "2025-02-01T00:00:00Z"
//        },
//        {
//            "Time": "2025-03-01T15:14:01Z"
//        }
//    ]
//}
```
