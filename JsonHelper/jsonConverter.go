package JsonHelper

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/heroku/GoCashager/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func ProvideUserInfo(mongoUser bson.M) []byte {
	var x string = fmt.Sprintf("%v", mongoUser["totalBalance"])
	totbal, _ := strconv.ParseInt(x, 10, 64)
	user := &utils.UserInfo{
		Uid:          fmt.Sprintf("%v", mongoUser["uid"]),
		FirstName:    fmt.Sprintf("%v", mongoUser["firstName"]),
		LastName:     fmt.Sprintf("%v", mongoUser["lastName"]),
		Totalbalance: totbal,
	}
	e, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
	}
	return e
}

func ProvideAllActivities(mongoUserActivities []map[string]string) []byte {
	var acts = utils.Activities{
		Activities: mongoUserActivities,
	}
	e, _ := json.Marshal(&acts)
	return e
}
