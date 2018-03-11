package tools

import (
	"fmt"
	"github.com/hzxiao/goutil"
	"github.com/juju/errors"
	"gopkg.in/mgo.v2/bson"
	"regexp"
	"strconv"
	"strings"
)

func GenerateUniqueId() string {
	return bson.NewObjectId().Hex()
}

func ToBsonMap(m goutil.Map) bson.M {
	if m == nil {
		return nil
	}
	bm := bson.M{}
	for k := range m {
		bm[k] = m[k]
	}

	return bm
}

func Struct2BsonMap(v interface{}) (bson.M, error) {
	buf, err := bson.Marshal(v)
	if err != nil {
		return nil, err
	}

	var bm bson.M
	err = bson.Unmarshal(buf, &bm)
	if err != nil {
		return nil, err
	}
	return bm, nil
}

func ParseRegex(str string) string {
	var data = regexp.QuoteMeta(str)
	var array = strings.Split(data, " ")
	var result = ".*"
	for _, item := range array {
		result += item + ".*"
	}
	return result
}

func ReserveDecimalFractionOf(value float64, n int) (float64, error) {
	if n < 0 {
		return value, errors.New("n can not less than zero")
	}
	f := "%0." + strconv.Itoa(n) + "f"
	s := fmt.Sprintf(f, value)
	return strconv.ParseFloat(s, 64)
}
