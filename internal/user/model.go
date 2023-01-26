package user

type User struct {
	Id           string         `json:"id" bson:"_id,omitempty"`
	TgId         int            `json:"tg_id" bson:"tg_id"`
	Intaractions []Intaractions `json:"intaractions" bson:"intaractions"`
}

type Intaractions struct {
	Date  string         `json:"date" bson:"date"`
	Views map[string]int `json:"views" bson:"views"`
}

type UserDTO struct {
	TgId  int            `json:"tg_id"`
	Views map[string]int `json:"views"`
}

type GetStatUserDTO struct {
	TgId int `json:"tg_id"`
}
