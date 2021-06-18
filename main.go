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
		var s string
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
			s = res
		}()

		// var allItems []bson.M
		// if err = cursor.All(ctx, &allItems); err != nil {
		// 	log.Fatalln(err)

		// 	return
		// }
		// //fmt.Print(allItems[0])
		// var response []byte = JsonHelper.ProvideUserInfo(allItems[0])
		// var res = string(response)
		// fmt.Print(res)
		c.String(200, s)
	})

	router.Run(":" + port)
}
