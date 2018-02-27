package goutil

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Map map[string]interface{}

func (m Map) Set(key string, value interface{}) {
	m[key] = value
}

func (m Map) Get(key string) interface{} {
	return m[key]
}

func (m Map) GetInt64(key string) int64 {
	if v, ok := m[key]; ok {
		return Int64(v)
	}
	return 0
}

func (m Map) GetUint64(key string) uint64 {
	if v, ok := m[key]; ok {
		return Uint64(v)
	}
	return 0
}

func (m Map) GetBool(key string) bool {
	if v, ok := m[key]; ok {
		return Bool(v)
	}
	return false
}

func (m Map) GetString(key string) string {
	if v, ok := m[key]; ok {
		return String(v)
	}
	return ""
}

func (m Map) GetFloat64(key string) float64 {
	if v, ok := m[key]; ok {
		return Float64(v)
	}
	return 0
}

func (m Map) GetMap(key string) Map {
	if v, ok := m[key]; ok {
		return MapV(v)
	}
	return nil
}

func (m Map) GetArray(key string) []interface{} {
	if v, ok := m[key]; ok {
		return ArrayV(v)
	}
	return nil
}

func (m Map) GetStringArray(key string) []string {
	if v, ok := m[key]; ok {
		return StringArrayV(v)
	}
	return nil
}

func (m Map) GetInt64Array(key string) []int64 {
	if v, ok := m[key]; ok {
		return Int64ArrayV(v)
	}
	return nil
}

func (m Map) GetFloat64Array(key string) []float64 {
	if v, ok := m[key]; ok {
		return Float64ArrayV(v)
	}
	return nil
}

func (m Map) GetMapArray(key string) []Map {
	if v, ok := m[key]; ok {
		return MapArrayV(v)
	}
	return nil
}

func (m Map) Exist(key string) bool {
	_, found := m[key]
	return found
}

func (m Map) GetP(path string) (interface{}, error) {
	return getValueByPath(m, path)
}

func (m Map) GetBoolP(path string) bool {
	v, _ := m.GetP(path)
	return Bool(v)
}

func (m Map) GetInt64P(path string) int64 {
	v, _ := m.GetP(path)
	return Int64(v)
}

func (m Map) GetUint64P(path string) uint64 {
	v, _ := m.GetP(path)
	return Uint64(v)
}

func (m Map) GetFloat64P(path string) float64 {
	v, _ := m.GetP(path)
	return Float64(v)
}

func (m Map) GetStringP(path string) string {
	v, _ := m.GetP(path)
	return String(v)
}

func (m Map) GetMapP(path string) Map {
	v, _ := m.GetP(path)
	return MapV(v)
}

func (m Map) GetStringArrayP(path string) []string {
	v, _ := m.GetP(path)
	return StringArrayV(v)
}

func (m Map) GetInt64ArrayP(path string) []int64 {
	v, _ := m.GetP(path)
	return Int64ArrayV(v)
}

func (m Map) GetFloat64ArrayP(path string) []float64 {
	v, _ := m.GetP(path)
	return Float64ArrayV(v)
}

func (m Map) GetMapArrayP(path string) []Map {
	v, _ := m.GetP(path)
	return MapArrayV(v)
}

func getValueByPath(m Map, path string) (interface{}, error) {
	path = strings.TrimPrefix(path, "/")
	keys := strings.Split(path, "/")
	var v interface{} = m
	for i := 0; i < len(keys); i++ {
		if v == nil {
			break
		}
		switch reflect.TypeOf(v).Kind() {
		case reflect.Slice:
			ary := ArrayV(v)
			if keys[i] == "@len" { //check if have @len
				return len(ary), nil //return the array length
			}
			idx, err := strconv.Atoi(keys[i]) //get the target index.
			if err != nil {
				return nil, fmt.Errorf("invalid array index(/%v)", strings.Join(keys[:i+1], "/"))
			}
			if idx >= len(ary) || idx < 0 { //check index valid
				return nil, fmt.Errorf(
					"array out of index in path(/%v)", strings.Join(keys[:i+1], "/"),
				)
			}
			v = ary[idx]
		case reflect.Map:
			tm := MapV(v) //check map covert
			if tm == nil {
				return nil, fmt.Errorf(
					"invalid map in path(/%v)", strings.Join(keys[:i], "/"),
				)
			}
			v = tm.Get(keys[i])
		default:
			return nil, fmt.Errorf(
				"invalid type(%v) in path(/%v)",
				reflect.TypeOf(v).Kind(), strings.Join(keys[:i], "/"),
			)
		}
	}

	if v == nil {
		return nil, fmt.Errorf("value not found in path(/%v)", strings.Join(keys, "/"))
	}
	return v, nil
}
