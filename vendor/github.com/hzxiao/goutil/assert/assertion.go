package assert

import (
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func Equal(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		_, file, line, _ := runtime.Caller(1)
		t.Logf("\033[31m%s:%d:\n\n\t   %#v (expected)\n\n\t!= %#v (actual)\033[39m\n\n",
			filepath.Base(file), line, expected, actual)
		t.FailNow()
	}
}

func NotEqual(t *testing.T, expected, actual interface{}) {
	if reflect.DeepEqual(expected, actual) {
		_, file, line, _ := runtime.Caller(1)
		t.Logf("\033[31m%s:%d:\n\n\tnexp: %#v\n\n\tgot:  %#v\033[39m\n\n",
			filepath.Base(file), line, expected, actual)
		t.FailNow()
	}
}

func NotNil(t *testing.T, obj interface{}) {
	if isNil(obj) {
		_, file, line, _ := runtime.Caller(1)
		t.Logf("\033[31m%s:%d:\n\n\tExpected value not to be <nil>\033[39m\n\n",
			filepath.Base(file), line, obj)
		t.FailNow()
	}
}

func Nil(t *testing.T, obj interface{}) {
	if !isNil(obj) {
		_, file, line, _ := runtime.Caller(1)
		t.Logf("\033[31m%s:%d:\n\n\t   <nil> (expected)\n\n\t!= %#v (actual)\033[39m\n\n",
			filepath.Base(file), line, obj)
		t.FailNow()
	}
}
func isNil(obj interface{}) bool {
	if obj == nil {
		return true
	}

	value := reflect.ValueOf(obj)
	kind := value.Kind()
	if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
		return true
	}
	return false
}

func Len(t *testing.T, obj interface{}, length int) {
	ok, l := getLen(obj)
	if !ok {
		_, file, line, _ := runtime.Caller(1)
		t.Logf("\033[31m%s:%d:\n\n\t   can not get length of %#v\n\n\t \033[39m\n\n",
			filepath.Base(file), line, obj)
		t.FailNow()
	}
	if l != length {
		_, file, line, _ := runtime.Caller(1)
		t.Logf("\033[31m%s:%d:\n\n\t   %#v (expected)\n\n\t!= %#v (actual)\033[39m\n\n",
			filepath.Base(file), line, length, l)
		t.FailNow()
	}
}

//getLen try to get length of obj
//return (false, 0) if impossible
func getLen(x interface{}) (ok bool, length int) {
	v := reflect.ValueOf(x)
	defer func() {
		if err := recover(); err != nil {
			ok = false
		}
	}()
	return true, v.Len()
}

func True(t *testing.T, value bool) {
	if !value {
		_, file, line, _ := runtime.Caller(1)
		t.Logf("\033[31m%s:%d:\n\n\t   true (expected)\n\n\t!= false (actual)\033[39m\n\n",
			filepath.Base(file), line)
		t.FailNow()
	}
}

func False(t *testing.T, value bool) {
	if value {
		_, file, line, _ := runtime.Caller(1)
		t.Logf("\033[31m%s:%d:\n\n\t   false (expected)\n\n\t!= true (actual)\033[39m\n\n",
			filepath.Base(file), line)
		t.FailNow()
	}
}

func Error(t *testing.T, err error) bool {
	if err == nil {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("\033[31m%s:%d:\n\n\t   error (expected)\n\n\t!= <nil> (actual)\033[39m\n\n",
			filepath.Base(file), line)
		return false
	}
	return true
}

func NoError(t *testing.T, err error) bool {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("\033[31m%s:%d:\n\n\t   <nil> (expected)\n\n\t!= error (actual)\033[39m\n\n",
			filepath.Base(file), line)
		return false
	}
	return true
}
