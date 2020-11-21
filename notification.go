package main

import (
	"context"
	"encoding/json"
	"github.com/bikashsapkota/spinstatz-notification/configuration"
	"github.com/bikashsapkota/spinstatz-notification/model"
	"github.com/go-redis/redis/v8"
	"github.com/jasonlvhit/gocron"
	"log"

	db_model "github.com/bikashsapkota/go_db/model"
	//db_storage "github.com/bikashsapkota/go_db"
)

func list_contains(list[] string, str string) bool {
	for _, element := range list {
		if element == str {
			return true
		}
	}
	return false
}

var (
	ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis-12902.c73.us-east-1-2.ec2.cloud.redislabs.com:12902",
		Password: "RtEwuhB1OOkVF1rZS4foQnDRbgycg04P",
		DB:       0,
		})
	DBHelper = configuration.DBConfig.DBInterface
)


func handleMessage(message string){
	MsgObj := model.NotificationMessage{}
	error_unMarshal := json.Unmarshal([]byte(message), &MsgObj)
	if error_unMarshal != nil {
		log.Println("parse error", error_unMarshal)
	}

	for _, element := range MsgObj.Destination {
		if element == configuration.Database {
			log.Println("Saving notification to database")
			notification := db_model.Notifications{
				CreatedAt: MsgObj.CreatedAt,
			}
			log.Println(notification)
			log.Println(configuration.DBConfig.DBInterface.SaveNotification(notification))
			//_, err := DBHelper.SaveNotification(notification)
			//
			//if err != nil {
			//	log.Println(err)
			//}


		}else if element == configuration.Email {

		}else if element == configuration.Live {

		}else if element == configuration.Push {

		}else {

		}
	}



}

func StartRadisConsumer(){
	sub := rdb.Subscribe(ctx, "mychannel")
	ch := sub.Channel()
	for msg := range ch {
		log.Println(msg.Channel, msg.Payload)
		handleMessage(msg.Payload)
	}
}


func StartRadisProducer()  {
	log.Println("publishing redis")
	rdb.Publish(
		ctx, "mychannel",
		`{"source_app_id":"web","created_at":"2006-01-02T15:04:05.1+05:45","type":"notifi","destination":["database","web"]}`)
}

func main()  {
	configuration.LoadConfig()

	go StartRadisConsumer()
	gocron.Every(10).Seconds().Do(StartRadisProducer)
	<-gocron.Start()
	log.Println("end")
}
