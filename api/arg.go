package api

import (
	"github.com/hzxiao/goutil"
	"github.com/juju/errors"
	"strconv"
	"strings"
)

type Arg struct {
	Key          string
	Value        interface{}
	DefaultValue string
	Type         string
	Require      bool
}

func CheckURLArg(formValue map[string][]string, args []*Arg) (goutil.Map, error) {
	argMap := goutil.Map{}
	if len(args) == 0 {
		return argMap, nil
	}

	for _, arg := range args {
		vs := formValue[arg.Key]
		var v string
		if len(vs) == 0 {
			if arg.Require {
				return nil, errors.Errorf("require %v field", arg.Key)
			}
			if arg.DefaultValue != "" {
				v = arg.DefaultValue
			} else {
				continue
			}
		} else {
			v = vs[0]
		}
		switch arg.Type {
		case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint32", "uint64":
			i, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, errors.Errorf("convert to int64 err(%v) by key(%v), value(%v)", err, arg.Key, v)
			}
			argMap.Set(arg.Key, i)
		case "float32", "float64":
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, errors.Errorf("convert to float64 err(%v) by key(%v), value(%v)", err, arg.Key, v)
			}
			argMap.Set(arg.Key, f)
		case "bool":
			v = strings.ToLower(v)
			if v == "1" || v == "true" {
				argMap.Set(arg.Key, true)
			} else if v == "0" || v == "false" {
				argMap.Set(arg.Key, false)
			} else {
				return nil, errors.Errorf("invalid bool type of key(%v), value(%v)", arg.Key, v)
			}
		case "string":
			if v != "" {
				argMap.Set(arg.Key, v)
			}
		default:
			return nil, errors.Errorf("unknown type(%v) of key(%v)", arg.Type, arg.Key)
		}
	}

	return argMap, nil
}

func TakeByKeys(inMap goutil.Map, keys ...string) goutil.Map {
	outMap := goutil.Map{}
	for i := range keys {
		v, ok := inMap[keys[i]]
		if ok {
			outMap.Set(keys[i], v)
		}
	}
	return outMap
}
