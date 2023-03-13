package domain

type KafkaRepository interface {
	TopicConsume() error
}

type User struct {
	FirstName string      `json:"first_name" structs:"first_name"`
	LastName  string      `json:"last_name" structs:"last_name"`
	Email     string      `json:"email" structs:"email"`
	Phone     PhoneStruct `json:"phone" structs:"phone"`
}

type PhoneStruct struct {
	Number      string `json:"number" bson:"number" structs:"number"`
	CountryCode string `json:"country_code" bson:"country_code" structs:"country_code"`
}
