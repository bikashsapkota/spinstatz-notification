package configuration

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"github.com/Shopify/sarama"
	"github.com/bikashsapkota/go_db"
	"github.com/fatih/structs"
	"gopkg.in/yaml.v2"
)

const (
	file = "configuration/config.yml"
)

var globalConfig GlobalConfig
var Config Conf
var DBConfig DBConf
//var KafkaConfig KafkaConf

type GlobalConfig struct {
	Default    Conf
	Production interface{}
	Staging    interface{}
}

type Conf struct {
	LogFile          string  `yaml:"logfile"`
	BaseUrl          string  `yaml:"base_url"`
	PdfUrl           string  `yaml:"pdf_url"`
	ConsumerGroup    string  `yaml:"consumer_group"`
	Topic            topic   `yaml:"topic"`
	Message          message `yaml:"message"`
	WebhookUrl       string  `yaml:"webhook_url"`
	WebhookSecretKey string  `yaml:"webhook_secret_key"`
	OnesignalAppId 	string	`yaml:"onesignal_app_id"`
	PusherToken		string	`yaml:"pusher_token"`
	RedisHost		string	`yaml:"redis_host"`
	RedisPort		string	`yaml:"redis_port"`
	RedisPassword	string	`yaml:"redis_password"`
	MailHost	string	`yaml:"mail_host"`
	MailPort	int	`yaml:"mail_port"`
	MailUsername	string	`yaml:"mail_username"`
	MailPassword	string	`yaml:"mail_password"`
	MailEncryption	string	`yaml:"mail_encryption"`
	MailFrom	string	`yaml:"mail_from"`
}

type DBConf struct {
	DBInterface go_db.DatabaseService
}


type topic struct {
	EventAdded    string   "event_added"
	EventModified string   "event_modified"
	AudioAdded    string   "audio_added"
	AudioPlayed   string   "audio_played"
	KafkaTopics   []string `yaml:"kafka_topics"`
}

type message struct {
	SourceAppId  string `yaml:"source_app_id"`
	SenderName   string `yaml:"sender_name"`
	SenderEmail  string `yaml:"sender_email"`
	SenderNumber string `yaml:"sender_number"`
}

func SetupLogging(file io.Writer) {
	log.SetOutput(file)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("[notification] ")
	sarama.Logger = log.New(file, "[Sarama] ", log.LstdFlags)
}

func LoadConfig() {
	source, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println(err)
	}
	err = yaml.Unmarshal(source, &globalConfig)
	if err != nil {
		log.Println(err)
	}
	Config = globalConfig.Default

	switch os.Getenv("GO_ENV") {
	case "staging":
		cF := structs.New(&Config)
		fillStruct(cF.Fields(), globalConfig.Staging)

	case "production":
		cF := structs.New(&Config)
		fillStruct(cF.Fields(), globalConfig.Production)
	default:
		log.Panicln("Set environment variable $GO_ENV. It can be either staging or production.")
	}

	f, err := os.OpenFile(Config.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}

	SetupLogging(f)

	DBConfig.DBInterface = NewDB()
	DBConfig.DBInterface.InitDb(f, 5, 5)

	fmt.Println(os.Getenv("GO_ENV"))
	log.Printf("Config file loaded successfully.")
}

func LoadTestConfig() {
	DBConfig.DBInterface = NewTestDB()
	source, err := ioutil.ReadFile(os.Getenv("CONFIG_BOSS"))
	if err != nil {
		log.Println(err)
	}
	err = yaml.Unmarshal(source, &globalConfig)
	if err != nil {
		log.Println(err)
	}
	Config = globalConfig.Default
	switch os.Getenv("GO_ENV") {
	case "staging":
		cF := structs.New(&Config)
		fillStruct(cF.Fields(), globalConfig.Staging)

	case "production":
		cF := structs.New(&Config)
		fillStruct(cF.Fields(), globalConfig.Production)
	default:
		log.Panicln("Set environment variable $GO_ENV. It can be either staging or production.")
	}
	log.Printf("Test Config file loaded successfully.")
}

func NewDB() go_db.DatabaseService {
	return &go_db.PgDatabase{}
}

func NewTestDB() go_db.DatabaseService {
	return &go_db.MockDatabase{}
}


func fillStruct(fields []*structs.Field, mapData interface{}) {
	switch reflect.TypeOf(mapData).Kind() {
	case reflect.Map:
		structMap := mapData.(map[interface{}]interface{})
		for k, v := range structMap {
			for _, field := range fields {
				if field.Tag("yaml") == k.(string) {
					if reflect.TypeOf(v).Kind() == reflect.Map {
						fillStruct(field.Fields(), v)
					} else if reflect.TypeOf(v).Kind() == reflect.Slice {
						sliceMapData := reflect.ValueOf(v)

						if sliceMapData.Len() > 0 {
							switch reflect.TypeOf(sliceMapData.Index(0).Interface()).Kind() {
							case reflect.String:
								sliceVal := make([]string, sliceMapData.Len())

								for i := 0; i < sliceMapData.Len(); i++ {
									sliceVal[i] = sliceMapData.Index(i).Interface().(string)
								}
								field.Set(sliceVal)

							case reflect.Int:
								sliceVal := make([]int, sliceMapData.Len())

								for i := 0; i < sliceMapData.Len(); i++ {
									sliceVal[i] = sliceMapData.Index(i).Interface().(int)
								}
								field.Set(sliceVal)
							case reflect.Bool:
								sliceVal := make([]bool, sliceMapData.Len())

								for i := 0; i < sliceMapData.Len(); i++ {
									sliceVal[i] = sliceMapData.Index(i).Interface().(bool)
								}
								field.Set(sliceVal)
							default:
								log.Println("Unknown data type: ", reflect.TypeOf(sliceMapData.Index(0).Interface()))
								return
							}
						}
					} else {
						field.Set(v)
					}
					break
				}
			}
		}
	default:
		log.Println("Type: ", reflect.TypeOf(mapData).Kind())
	}
}
