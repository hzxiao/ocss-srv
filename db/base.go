package db

import "gopkg.in/mgo.v2/bson"

func count(collectionName string, cond bson.M) (int, error) {
	return C(collectionName).Find(cond).Count()
}

func one(collectionName string, cond, selector bson.M, v interface{}) error {
	return C(collectionName).Find(cond).Select(selector).One(v)
}

func list(collectionName string, cond, selector bson.M, sort []string, skip, limit int, v interface{}) (int, error) {
	query := C(collectionName).Find(cond).Sort(sort...).Select(selector)
	count, err := query.Count()
	if err != nil {
		return 0, err
	}
	if skip > 0 {
		query = query.Skip(skip)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	return count, query.All(v)
}
