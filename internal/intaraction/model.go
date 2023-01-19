package intaraction

type Intaraction struct {
	Id    string         `json:"id" bson:"_id,omitempty"`
	Date  string         `json:"date" bson:"date"`
	Views map[string]int `json:"views" bson:"views"`
}

type IntaractionDTO struct {
	Views map[string]int `json:"views"`
}
