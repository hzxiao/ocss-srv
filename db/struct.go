package db

type User struct {
	ID       string `bson:"_id" json:"id"`
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
	Icon     string `bson:"icon" json:"icon"`
	Role     int    `bson:"role" json:"role"`
	Status   int    `bson:"status" json:"status"`
	Create   int64  `bson:"create" json:"create"`
	Update   int64  `bson:"update" json:"update"`
}
