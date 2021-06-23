package JsonHelper

import (
	"fmt"
	"strconv"

	"github.com/heroku/GoCashager/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func ProvideUserInfo(mongoUser bson.M) interface{} {
	var x string = fmt.Sprintf("%v", mongoUser["totalBalance"])
	totbal, _ := strconv.ParseInt(x, 10, 64)
	user := &utils.UserInfo{
		Uid:          fmt.Sprintf("%v", mongoUser["uid"]),
		FirstName:    fmt.Sprintf("%v", mongoUser["firstName"]),
		LastName:     fmt.Sprintf("%v", mongoUser["lastName"]),
		Totalbalance: totbal,
	}
	return user
}

func ProvideAllActivities(mongoUserActivities []map[string]string) interface{} {
	var acts = utils.Activities{
		Activities: mongoUserActivities,
	}
	return acts
}
