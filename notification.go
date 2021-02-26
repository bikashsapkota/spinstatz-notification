package main

import (
	"context"
	"encoding/json"
	"github.com/bikashsapkota/spinstatz-notification/configuration"
	"github.com/bikashsapkota/spinstatz-notification/model"
	"github.com/go-redis/redis/v8"
	"github.com/jasonlvhit/gocron"
	"log"
	"github.com/bikashsapkota/spinstatz-notification/controller"

	db_model "github.com/bikashsapkota/go_db/model"
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
)


func handleMessage(message string){
	MsgObj := model.NotificationMessage{}
	error_unMarshal := json.Unmarshal([]byte(message), &MsgObj)
	if error_unMarshal != nil {
		log.Println("parse error", error_unMarshal)
	}

	log.Println("unprased ", message)
	log.Println("parsed ", MsgObj)
	log.Println("destination ", MsgObj.Destination)

	for _, element := range MsgObj.Destination {
		log.Println("processing "+element+" Notification")
		if element == configuration.Database {
			log.Println("Saving notification to database")
			notification := db_model.Notifications{
				CreatedAt: MsgObj.CreatedAt,
				UserId: MsgObj.UserId,
				Message: MsgObj.Message,
				Href: MsgObj.Href,
				Seen: MsgObj.Seen,
				Type: MsgObj.Type,
				ReferenceId: MsgObj.ReferenceId,
				Subject: controller.Get_subject(MsgObj.Type),
				Image: MsgObj.Image,

			}
			log.Println(notification)
			log.Println(configuration.DBConfig.DBInterface.SaveNotification(notification))

		}else if element == configuration.Email {
			controller.HandleEmailNotification(MsgObj)
		}else if element == configuration.Live {

		}else if element == configuration.Push {
			controller.HandleOnesignalNotification(MsgObj)
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
		`{
			"source_app_id":"web",
			"created_at":"2006-01-02T15:04:05.1+05:45",
			"destination":["push","web"],
			"one_signal_token": "f8e1b64b-8866-4fb8-bf61-16797a2252c8",
			"type": "payment_withdrawl_requested"
		}`)
}

func main()  {
	configuration.LoadConfig()

	StartRadisConsumer()
	gocron.Every(10).Seconds().Do(StartRadisProducer)
	<-gocron.Start()
	log.Println("end")
}


