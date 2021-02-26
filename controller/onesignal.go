package controller

import (
	"github.com/bikashsapkota/spinstatz-notification/configuration"
	"github.com/bikashsapkota/spinstatz-notification/model"
	"github.com/tbalthazar/onesignal-go"
	"log"
	gomail "gopkg.in/mail.v2"
	"strings"
	"crypto/tls"
)

var (
	onesig_client = onesignal.NewClient(nil)
)

func Get_subject(str string) string {
	if str == configuration.NewCampaignAdded {
		return "Campaign added"
	}else if str == configuration.PaymentWithdrawlRequested {
		return "Payment Request Successful"
	}
	return strings.Replace(str, "_"," ", 2)
}

func get_notification_request_by_type(msgObj model.NotificationMessage) (*onesignal.NotificationRequest, error) {
	heading := Get_subject(msgObj.Type)
	log.Println("heading: "+ heading)

	notificationReq := &onesignal.NotificationRequest{
		AppID:            configuration.Config.OnesignalAppId,
		Contents:         map[string]string{"en": msgObj.Message},
		Headings: 			map[string]string{"en": heading},
		IsIOS:            true,
		IncludePlayerIDs: []string{msgObj.OneSignalToken},
		Data: msgObj.Data,
	}

	return notificationReq, nil
}

func HandleOnesignalNotification(msgObj model.NotificationMessage)  {

	if msgObj.OneSignalToken != "" {
		log.Println("Sending notification to "+ msgObj.OneSignalToken +" from app "+ configuration.Config.OnesignalAppId)
		notificationReq, _ := get_notification_request_by_type(msgObj)
		_, _, err := onesig_client.Notifications.Create(notificationReq)
		if err != nil {
			log.Println("OneSignal Error Occured", err)
		}
	}else {
		log.Println("Onesignal user token not found")
	}
}

func HandleEmailNotification(msgObg model.NotificationMessage){
	log.Println(msgObg.EmailData["to"], msgObg.EmailData["msg"], configuration.Config.MailHost, configuration.Config.MailPort, configuration.Config.MailUsername ,configuration.Config.MailPassword)
	to := msgObg.EmailData["to"]
	message := msgObg.EmailData["message"]

	m := gomail.NewMessage()
	m.SetHeader("From", configuration.Config.MailFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject","Withdraw Requested")
	m.SetBody("text/plain", message)

	d := gomail.NewDialer(configuration.Config.MailHost, configuration.Config.MailPort, configuration.Config.MailUsername, configuration.Config.MailPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		log.Println(err)
		panic(err)
	}
	log.Println("Email Sent Successfully!")
}
