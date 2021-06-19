package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/heroku/GoCashager/JsonHelper"
	"github.com/heroku/GoCashager/utils"
	_ "github.com/heroku/x/hmetrics/onload"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getUserInfo(client_mongo *mongo.Client, ctx context.Context, uid string) <-chan string {
	finalres := make(chan string)
	go func() {
		cursor, err := client_mongo.Database("Cashager").Collection("user+"+uid).Find(ctx, bson.M{"type": "baseInfo"})
		if err != nil {
			log.Fatalln(err)
		}
		var allItems []bson.M
		if err = cursor.All(ctx, &allItems); err != nil {
			log.Fatalln(err)
		}
		var response []byte = JsonHelper.ProvideUserInfo(allItems[0])
		var res = string(response)
		finalres <- res
	}()
	return finalres
}
func getActivities(client_mongo *mongo.Client, ctx context.Context, uid string) <-chan string {
	finalres := make(chan string)
	go func() {
		var activities = utils.Activities{}
		err := client_mongo.Database("Cashager").Collection("user+"+uid).FindOne(ctx, bson.M{"type": "activities"}).Decode(&activities)
		if err != nil {
			log.Fatalln(err)

			return
		}
		var response string = string(JsonHelper.ProvideAllActivities(activities.Activities))
		var res = string(response)
		finalres <- res
	}()
	return finalres
}

func addActivity(client_mongo *mongo.Client, ctx context.Context, uid string, name string, desc string, amount string, id string) <-chan int {
	finalres := make(chan int)
	go func() {
		var singleActivity = utils.Activity{}
		var activityMap = map[string]map[string]string{"activities": {"name": singleActivity.Name, "desc": singleActivity.Desc, "amount": singleActivity.Amount, "id": singleActivity.Id}}
		_, err := client_mongo.Database("Cashager").Collection("user+"+uid).UpdateOne(ctx, bson.M{"type": "activities"}, bson.M{"$push": activityMap})
		if err != nil {
			log.Fatalln(err)
		}
		finalres <- 200
	}()
	return finalres
}

func main() {
	var client_mongo *mongo.Client
	var ctx = context.TODO()
	client_mongo, _ = mongo.NewClient(options.Client().ApplyURI(utils.MONGO_ACCESS_URL))
	ctx, _ = context.WithTimeout(context.Background(), 1*time.Hour)
	err := client_mongo.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/index.html")

	router.GET(utils.GET_USER_INFO+"/:uid", func(c *gin.Context) {
		var uid = c.Param("uid")
		res := <-getUserInfo(client_mongo, ctx, uid)
		c.String(200, res)
	})
	router.GET(utils.GET_USER_ACTIVITIES+"/:uid", func(c *gin.Context) {
		var uid = c.Param("uid")
		res := <-getActivities(client_mongo, ctx, uid)
		c.String(200, res)
	})

	router.POST(utils.ADD_ACTIVITY_ROUTE+"/:uid", func(c *gin.Context) {
		var uid = c.Param("uid")
		resCode := <-addActivity(client_mongo, ctx, uid, c.PostForm("name"), c.PostForm("desc"), c.PostForm("amount"), c.PostForm("id"))
		c.String(resCode, "Success")
	})

	router.Run(":" + port)
}
