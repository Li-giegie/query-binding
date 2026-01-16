package query_binding

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

type Base struct {
	Int      int
	Int8     int8
	Int16    int16
	Int32    int32
	Int64    int64
	Uint     uint
	Uint8    uint8
	Uint16   uint16
	Uint32   uint32
	Uint64   uint64
	Float32  float32
	Float64  float64
	String   string
	Bytes    []byte
	Time     time.Time
	PInt     *int
	PInt8    *int8
	PInt16   *int16
	PInt32   *int32
	PInt64   *int64
	PUint    *uint
	PUint8   *uint8
	PUint16  *uint16
	PUint32  *uint32
	PUint64  *uint64
	PFloat32 *float32
	PFloat64 *float64
	PString  *string
	PBytes   *[]byte
	PTime    *time.Time
}

func TestQueryBase(t *testing.T) {
	var query = map[string][]string{
		"Int":      {"-1"},
		"Int8":     {"-2"},
		"Int16":    {"-3"},
		"Int32":    {"-4"},
		"Int64":    {"-5"},
		"Uint":     {"6"},
		"Uint8":    {"7"},
		"Uint16":   {"8"},
		"Uint32":   {"9"},
		"Uint64":   {"10"},
		"Float32":  {"11.1"},
		"Float64":  {"12.2"},
		"String":   {"hello"},
		"Bytes":    {"1", "2", "65"},
		"Time":     {time.Now().Format(time.RFC3339)},
		"PInt":     {"-13"},
		"PInt8":    {"-14"},
		"PInt16":   {"-15"},
		"PInt32":   {"-16"},
		"PInt64":   {"-17"},
		"PUint":    {"18"},
		"PUint8":   {"19"},
		"PUint16":  {"20"},
		"PUint32":  {"21"},
		"PUint64":  {"22"},
		"PFloat32": {"23.123"},
		"PFloat64": {"24.456"},
		"PString":  {"pointer string"},
		"PBytes":   {"1", "2", "3"},
		"PTime":    {time.Now().Format(time.RFC3339)},
	}
	var param Base
	err := Mapping(query, reflect.ValueOf(&param), "")
	if err != nil {
		t.Error(err)
	}
	Println("param", param)
}

type Slice struct {
	Int      []int
	Int8     []int8
	Int16    []int16
	Int32    []int32
	Int64    []int64
	Uint     []uint
	Uint8    []uint8
	Uint16   []uint16
	Uint32   []uint32
	Uint64   []uint64
	Float32  []float32
	Float64  []float64
	String   []string
	Bytes    []byte
	Time     []time.Time
	PInt     []*int
	PInt8    []*int8
	PInt16   []*int16
	PInt32   []*int32
	PInt64   []*int64
	PUint    []*uint
	PUint8   []*uint8
	PUint16  []*uint16
	PUint32  []*uint32
	PUint64  []*uint64
	PFloat32 []*float32
	PFloat64 []*float64
	PString  []*string
	PBytes   *[]byte
	PTime    []*time.Time
}

func TestQuerySlice(t *testing.T) {
	querys := map[string][]string{
		// 整数类型（有符号） - 最多 3 个值
		"Int":   []string{"-1", "0", "1"},
		"Int8":  []string{"-1", "0", "1"},
		"Int16": []string{"-1", "0", "1"},
		"Int32": []string{"-1", "0", "1"},
		"Int64": []string{"-1", "0", "1"},

		// 整数类型（无符号） - 最多 3 个值
		"Uint":   []string{"0", "1", "100"},
		"Uint8":  []string{"0", "1", "255"},
		"Uint16": []string{"0", "1", "65535"},
		"Uint32": []string{"0", "1", "4294967295"},
		"Uint64": []string{"0", "1", "18446744073709551615"},

		// 浮点数 - 最多 3 个值
		"Float32": []string{"-1.0", "0.0", "1.0"},
		"Float64": []string{"-1.0", "0.0", "1.0"},

		// 字符串 - 最多 3 个值
		"String": []string{"hello", "world", "123"},

		// ✅ Bytes 字段：必须是 ASCII 数字字符，最多 3 个
		"Bytes": []string{"0", "1", "2"},

		// 时间 - 最多 3 个值
		"Time": []string{
			"2023-01-01T00:00:00Z",
			"2023-12-06T12:00:00Z",
			"1970-01-01T00:00:00Z",
		},

		// 指针或切片类型字段 - 最多 3 个值（底层类型合法值）
		"PInt":     []string{"-1", "0", "1"},
		"PInt8":    []string{"-1", "0", "1"},
		"PInt16":   []string{"-1", "0", "1"},
		"PInt32":   []string{"-1", "0", "1"},
		"PInt64":   []string{"-1", "0", "1"},
		"PUint":    []string{"0", "1", "100"},
		"PUint8":   []string{"0", "1", "255"},
		"PUint16":  []string{"0", "1", "65535"},
		"PUint32":  []string{"0", "1", "4294967295"},
		"PUint64":  []string{"0", "1", "18446744073709551615"},
		"PFloat32": []string{"-1.0", "0.0", "1.0"},
		"PFloat64": []string{"-1.0", "0.0", "1.0"},
		"PString":  []string{"hello", "world", "123"},
		// ✅ PBytes 字段：必须是 ASCII 数字字符，最多 3 个
		"PBytes": []string{"0", "1", "2"},
		"PTime": []string{
			"2023-01-01T00:00:00Z",
			"2023-12-06T12:00:00Z",
			"1970-01-01T00:00:00Z",
		},
	}
	var param Slice
	err := Mapping(querys, reflect.ValueOf(&param), "")
	if err != nil {
		t.Error(err)
	}
	Println("param", param)
}

type C1 time.Time

func (c *C1) UnmarshalParam(vals []string) error {
	t, err := time.Parse(time.RFC3339, vals[0])
	if err != nil {
		return err
	}
	*c = C1(t)
	return nil
}

func (c C1) MarshalJSON() ([]byte, error) {
	return []byte("\"" + time.Time(c).Format(time.RFC3339) + "\""), nil
}

type C2 []string

func (c *C2) UnmarshalParam(vals []string) error {
	*c = vals
	return nil
}

type Custom struct {
	ABC string `form:"abc"`
	C1  C1     `form:"c1"`
	C2  C2     `form:"c2"`
	C3  []C1   `form:"c3"`
}

func TestQueryCustom(t *testing.T) {
	query := map[string][]string{
		"abc": {"1", "2"},
		"c1":  {time.Now().Format(time.RFC3339)},
		"c2":  {"a", "b", "c"},
		"c3":  {time.Now().Format(time.RFC3339)},
	}
	var param Custom
	err := Mapping(query, reflect.ValueOf(&param), "form")
	if err != nil {
		t.Fatal(err)
	}
	Println("param", param)
}

type Common struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}

type Anonymous struct {
	Name string `form:"name"`
	C    int    `form:"c"`
	Common
	CC Common
}

func TestQueryAnonymous(t *testing.T) {
	query := map[string][]string{
		"name":     {"1", "2"},
		"page":     {"10"},
		"pageSize": {"100"},
		"c":        {"1"},
		"CC":       {"1", "2"},
	}
	var param Anonymous
	err := Mapping(query, reflect.ValueOf(&param), "form")
	if err != nil {
		t.Fatal(err)
	}

	Println("param", param)

}

func Println(a ...any) {
	for _, a2 := range a {
		switch v := a2.(type) {
		case string:
			println(v)
		case *string:
			if v != nil {
				println(*v)
			}
		case []byte:
			println(string(v))
		case *[]byte:
			if v != nil {
				println(string(*v))
			}
		case error:
			println(v.Error())
		default:
			data, err := json.MarshalIndent(v, "", "\t")
			if err != nil {
				panic(err)
			}
			println(string(data))
		}
	}
}
