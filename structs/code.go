package structs

type RegStruct struct {
	Id       string `bson:"_id"`
	Name     string
	Surname  string
	Login    string
	Password string
	Balance  int
}

type LoginStruct struct {
	Login    string
	Password string
}

type Books struct {
	Name   string
	Author string
	Year   string
}
