package db

type User struct {
	ID       string `bson:"_id" json:"id"`
	Username string `bson:"username" json:"username"`
}
