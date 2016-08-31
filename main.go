package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/mattes/migrate/driver/mysql"
)

var dbdsn = "root:@tcp(localhost:3306)/rolex_activity?charset=utf8&parseTime=true&loc=Local"

type Context struct {
	DbClient *gorm.DB
}

var context Context

type Activity struct {
	//gorm.Model
	UniqId   string `json:"UniqId"`
	SourceIp string `json:"SourceIp"`
}

func main() {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	var err error
	context.DbClient, err = gorm.Open("mysql", dbdsn)
	if err != nil {
		panic(err)
	}
	defer context.DbClient.Close()

	context.DbClient.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&Activity{})

	router.GET("/infos", ActivityInfo)
	router.POST("/activities", CreateActivities)

	server := &http.Server{
		Addr:           ":4500",
		Handler:        router,
		MaxHeaderBytes: 1 << 20,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func ActivityInfo(ctx *gin.Context) {
	var activities []Activity
	if err := context.DbClient.
		Find(&activities).Error; err != nil {
		ctx.JSON(500, gin.H{"err": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"data": activities})
}

func CreateActivities(ctx *gin.Context) {
	var activity Activity
	if err := ctx.BindJSON(&activity); err != nil {
		ctx.JSON(500, gin.H{"err": err.Error()})
		return
	}

	activity.SourceIp = ctx.Request.RemoteAddr

	if err := context.DbClient.Save(activity).Error; err != nil {
		ctx.JSON(500, gin.H{"err": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"data": "done"})
}
