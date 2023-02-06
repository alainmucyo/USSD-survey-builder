package main

import (
	"context"
	"github.com/alainmucyo/ussd-go/sessionstores"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"survey-ussd/core/environment"
	"survey-ussd/core/service"
	"survey-ussd/core/topics"
	"survey-ussd/core/ussd"
	ussd_handler "survey-ussd/handlers/ussd-handler"
	"survey-ussd/store/kafka/consumer"
	"survey-ussd/store/kafka/producer"
	"survey-ussd/store/redis"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx := context.Background()
	envs := environment.New(
		os.Getenv("KAFKA_BROKER_URL"),
		os.Getenv("PORT"),
		os.Getenv("KAFKA_GROUP_ID"),
		os.Getenv("REDIS_URL"),
		os.Getenv("REDIS_PASSWORD"),
		os.Getenv("INCOMING_REQUEST"),
		os.Getenv("OUTGOING_RESPONSE"),
	)
	s := sessionstores.NewRedis(envs.RedisURL, envs.RedisPassword)
	println("Connecting to redis...")
	err = s.Connect()
	if err != nil {
		println("Unable to connect to Redis")
		panic(err.Error())
	}
	println("Connected to redis successfully")
	defer s.Close()
	cache := redis.New(envs, ctx)

	airtimeUssd := ussd.New(s)

	kafkaProducer := producer.New(envs)
	surveyService := service.New(cache, kafkaProducer)

	ussdHandler := ussd_handler.New(airtimeUssd, kafkaProducer, envs, surveyService)

	kafkaTopics := topics.New(envs, ussdHandler)
	kafkaConsumer := consumer.New(envs, ctx, kafkaTopics)
	go kafkaConsumer.Consume()
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong pong")
	})

	r.POST("/requests", ussdHandler.HandleUSSDRequests)

	err = r.Run(":" + envs.Port)
	if err != nil {
		log.Fatal(err)
	}
}
