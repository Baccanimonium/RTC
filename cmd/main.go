package main

import (
	"context"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
	"video-chat-app"
	"video-chat-app/src/Handlers"
	"video-chat-app/src/Repos"
	"video-chat-app/src/Services"
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

	repos := Repos.NewRepo(db)
	services := Services.NewService(repos)
	handlers := Handlers.NewHandler(services)

	srv := new(RTC.Server)

	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRouter()); err != nil {
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

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("server")

	return viper.ReadInConfig()
}
