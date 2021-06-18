package utils

type Activity struct {
	Name   string `json:"name"`
	Desc   string `json:"desc"`
	Amount string `json:"amount"`
	Id     string `json:"id"`
}

type UserInfo struct {
	Uid          string `json:"uid"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Totalbalance int64  `json:"totalBalance"`
}

type Activities struct {
	Activities []map[string]string `json:"activities"`
}
