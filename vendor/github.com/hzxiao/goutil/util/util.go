package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func Bool(v interface{}) bool {
	b, err := BoolE(v)
	if err != nil {
		return false
	}

	return b
}

func BoolE(v interface{}) (bool, error) {
	if v == nil {
		return false, fmt.Errorf("arg value is null")
	}
	b, ok := v.(bool)
	if ok {
		return b, nil
	}
	return false, fmt.Errorf("invalid bool value")
}

func Int64(v interface{}) int64 {
	val, err := Int64E(v)
	if err == nil {
		return val
	}
	return 0

}

func Int64E(v interface{}) (int64, error) {
	if v == nil {
		return 0, fmt.Errorf("arg value is null")
	}
	k := reflect.TypeOf(v)
	switch k.Kind() {
	case reflect.Int:
		return int64(v.(int)), nil
	case reflect.Int8:
		return int64(v.(int8)), nil
	case reflect.Int16:
		return int64(v.(int16)), nil
	case reflect.Int32:
		return int64(v.(int32)), nil
	case reflect.Int64:
		return v.(int64), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(Uint64(v)), nil
	case reflect.Float32, reflect.Float64:
		return int64(Float64(v)), nil
	case reflect.String:
		fv, err := strconv.ParseInt(v.(string), 10, 64)
		if err != nil {
			return 0, err
		}
		return fv, nil
	case reflect.Struct:
		if k.Name() == "Time" {
			return Timestamp(v.(time.Time)), nil
		}
		return 0, fmt.Errorf("incompactable kind(%v)", k.Kind())

	default:
		return 0, fmt.Errorf("incompactable kind(%v)", k.Kind())
	}
}

func Timestamp(t time.Time) int64 {
	return t.Local().UnixNano() / 1e6
}

func Uint64(v interface{}) uint64 {
	val, err := Uint64E(v)
	if err == nil {
		return val
	}
	return 0

}

func Uint64E(v interface{}) (uint64, error) {
	if v == nil {
		return 0, fmt.Errorf("arg value is null")
	}
	k := reflect.TypeOf(v)
	switch k.Kind() {
	case reflect.Uint:
		return uint64(v.(uint)), nil
	case reflect.Uint8:
		return uint64(v.(uint8)), nil
	case reflect.Uint16:
		return uint64(v.(uint16)), nil
	case reflect.Uint32:
		return uint64(v.(uint32)), nil
	case reflect.Uint64:
		return v.(uint64), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return uint64(Int64(v)), nil
	case reflect.Float32, reflect.Float64:
		return uint64(Float64(v)), nil
	case reflect.String:
		fv, err := strconv.ParseUint(v.(string), 10, 64)
		if err != nil {
			return 0, err
		}
		return fv, nil
	default:
		return 0, fmt.Errorf("incompactable kind(%v)", k.Kind().String())
	}
}

func Float64(v interface{}) float64 {
	val, err := Float64E(v)
	if err == nil {
		return val
	}
	return 0
}

func Float64E(v interface{}) (float64, error) {
	if v == nil {
		return 0, fmt.Errorf("arg value is null")
	}
	k := reflect.TypeOf(v)
	switch k.Kind() {
	case reflect.Float32:
		return float64(v.(float32)), nil
	case reflect.Float64:
		return float64(v.(float64)), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(Uint64(v)), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(Int64(v)), nil
	case reflect.String:
		if fv, err := strconv.ParseFloat(v.(string), 64); err == nil {
			return fv, nil
		} else {
			return 0, err
		}
	default:
		return 0, fmt.Errorf("incompactable kind(%v)", k.Kind().String())
	}
}

func String(v interface{}) string {
	if v == nil {
		return ""
	}
	switch reflect.TypeOf(v).Kind() {
	case reflect.String:
		return v.(string)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func MapV(v interface{}) Map {
	if mv, ok := v.(Map); ok {
		return mv
	} else if mv, ok := v.(map[string]interface{}); ok {
		return Map(mv)
	} else {
		return nil
	}
}

func ArrayV(v interface{}) []interface{} {
	if vals, ok := v.([]interface{}); ok {
		return vals
	}
	vals := reflect.ValueOf(v)
	if vals.Kind() != reflect.Slice {
		return nil
	}
	var vs = []interface{}{}
	for i := 0; i < vals.Len(); i++ {
		vs = append(vs, vals.Index(i).Interface())
	}
	return vs
}

func MapArrayV(v interface{}) []Map {
	var vals = ArrayV(v)
	if vals == nil {
		return nil
	}
	var ms = []Map{}
	for _, val := range vals {
		var mv = MapV(val)
		if mv == nil {
			return nil
		}
		ms = append(ms, mv)
	}
	return ms
}

func StringArrayV(v interface{}) []string {
	var vals = ArrayV(v)
	if vals == nil {
		return nil
	}
	var ms = []string{}
	for _, val := range vals {
		ms = append(ms, String(val))
	}
	return ms
}

func Int64ArrayV(v interface{}) []int64 {
	as := ArrayV(v)
	if as == nil {
		return nil
	}
	is := []int64{}
	for _, v := range as {
		iv, err := Int64E(v)
		if err != nil {
			return nil
		}
		is = append(is, iv)
	}
	return is
}

func Float64ArrayV(v interface{}) []float64 {
	as := ArrayV(v)
	if as == nil {
		return nil
	}
	is := []float64{}
	for _, v := range as {
		iv, err := Float64E(v)
		if err != nil {
			return nil
		}
		is = append(is, iv)
	}
	return is
}

func Struct2Map(s interface{}) Map {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil
	}
	var data = Map{}
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}

	return data
}

func Struct2Json(s interface{}) string {
	bys, err := json.Marshal(s)
	if err != nil {
		return "null"
	}

	return string(bys)
}
