package query_binding

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

var (
	Validator StructValidator
	Default   QueryBinding
)

type StructValidator interface {
	ValidateStruct(any) error
	Engine() any
}

type UnmarshalParam interface {
	UnmarshalParam(vals []string) error
}

type QueryBinding struct{}

func (*QueryBinding) Name() string {
	return "query"
}

func (q *QueryBinding) Bind(req *http.Request, obj any) error {
	rv := reflect.ValueOf(obj)
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("obj is not a struct")
	}
	values := req.URL.Query()
	err := Mapping(values, rv, "form")
	if err != nil {
		return err
	}
	if Validator == nil {
		return nil
	}
	return Validator.ValidateStruct(obj)
}

func Mapping(m map[string][]string, value reflect.Value, tagName string) (err error) {
	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() == reflect.Struct {
		for i := 0; i < value.NumField(); i++ {
			fv := value.Field(i)
			if !fv.CanSet() {
				continue
			}
			ft := value.Type().Field(i)
			fn := ft.Name
			if tagName != "" {
				fn = ft.Tag.Get(tagName)
			}
			qv := m[fn]
			if fn == "-" || len(qv) == 0 && !ft.Anonymous || (fn == "" && !ft.Anonymous) {
				continue
			}
			if ft.Anonymous {
				if up, ok := fv.Addr().Interface().(UnmarshalParam); ok {
					if err = up.UnmarshalParam(qv); err != nil {
						return err
					}
					continue
				}
				if err = Mapping(m, fv, tagName); err != nil {
					return err
				}
				continue
			}
			if err = mapping(fv, qv); err != nil {
				return err
			}
		}
	}
	return nil
}

func mapping(v reflect.Value, vals []string) error {
	if fn, ok := v.Addr().Interface().(UnmarshalParam); ok {
		return fn.UnmarshalParam(vals)
	}
	switch v.Kind() {
	case reflect.String:
		v.SetString(vals[0])
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(vals[0], 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(n)
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(vals[0], 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(n)
	case reflect.Uint8:
		n, err := strconv.ParseUint(vals[0], 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(n)
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(vals[0], 64)
		if err != nil {
			return err
		}
		v.SetFloat(n)
	case reflect.Bool:
		n, err := strconv.ParseBool(vals[0])
		if err != nil {
			return err
		}
		v.SetBool(n)
	case reflect.Struct:
		switch v.Interface().(type) {
		case time.Time:
			t, err := time.Parse(time.RFC3339, vals[0])
			if err != nil {
				return err
			}
			v.Set(reflect.ValueOf(t))
		default:
			err := json.Unmarshal([]byte(vals[0]), v.Addr().Interface())
			if err != nil {
				return err
			}
		}
	case reflect.Ptr:
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		err := mapping(v.Elem(), vals)
		if err != nil {
			return err
		}
	case reflect.Slice:
		val := reflect.MakeSlice(v.Type(), len(vals), len(vals))
		for i := 0; i < len(vals); i++ {
			err := mapping(val.Index(i), vals[i:])
			if err != nil {
				return err
			}
		}
		v.Set(val)
	case reflect.Array:
		if len(vals) > v.Len() {
			return errors.New("array length mismatch")
		}
		for i := 0; i < len(vals); i++ {
			err := mapping(v.Index(i), vals[i:])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Bind(req *http.Request, obj any) error {
	return Default.Bind(req, obj)
}
