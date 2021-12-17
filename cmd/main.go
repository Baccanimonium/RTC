package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"os/signal"
	"syscall"
	"time"
	"video-chat-app"
	"video-chat-app/src/Handlers"
	RTC2 "video-chat-app/src/RTC"
	"video-chat-app/src/Repos"
	"video-chat-app/src/Services"
	"video-chat-app/src/SocketHandlers"
	"video-chat-app/src/Tasks"
)

func main() {
	//server.AllRooms.Init()
	//
	//http.HandleFunc("/create", server.CreateRoomRequestHandler)
	//http.HandleFunc("/join", server.JoinRoomRequestHandler)
	//
	//log.Println("Starting Server on Port 8000")
	//err := http.ListenAndServe(":8000", nil)
	//
	//if err != nil {
	//	log.Fatal(err)
	//}

	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:rootPassXXX@127.0.0.1:27017"))

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	err = client.Ping(ctx, readpref.Primary())

	mongoDB := client.Database("drug-addicted")

	if err != nil {
		logrus.Errorf("error occured on mongodb connection close: %s", err.Error())
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"), // no password set
		DB:       0,                                 // use default DB
	})

	_, err = rdb.Ping(ctx).Result()

	if err != nil {
		logrus.Fatalf("failed to initialize redis: %s", err.Error())
	}

	db, err := RTC.NewPostgresDB(RTC.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}
	broadcastChan := make(chan RTC.BroadcastingMessage)
	repos := Repos.NewRepo(db, mongoDB, rdb)
	// инициализируем стейт сокет хаба
	socketHub := SocketHandlers.NewHub(repos, broadcastChan)
	services := Services.NewService(repos, broadcastChan)
	socketClientFactory := SocketHandlers.NewSocketClientFactory(services, socketHub)
	handlers := Handlers.NewHandler(services)

	// запускаем цикл отправки видео между RTC клиенами
	RTC2.PolingRTCClientsLoop()

	//запускаем горутину с прослушиванием каналов сокет хаба
	go socketHub.Run()

	backgroundTasks := Tasks.NewTaskManager(services, socketHub, rdb)

	backgroundTasks.Run()

	srv := new(RTC.Server)

	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRouter(socketClientFactory)); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Print("Start")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	logrus.Print("Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
	if err := rdb.Close(); err != nil {
		logrus.Errorf("error occured on redis connection close: %s", err.Error())
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("server")

	return viper.ReadInConfig()
}
