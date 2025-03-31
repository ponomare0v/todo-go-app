package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Импорт драйвера PostgreSQL
	"github.com/ponomare0v/todo-go-app/pkg/handler"
	"github.com/ponomare0v/todo-go-app/pkg/repository"
	"github.com/ponomare0v/todo-go-app/pkg/server"
	"github.com/ponomare0v/todo-go-app/pkg/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// @title Todo App API
// @version 1.0
// @description API Server for TodoList Application

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("ошибка инициализации конфигурации: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("ошибка загрузки переменных окружения: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("не удалось инициализоровать БД: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(server.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Print("TodoApp Started")

	quit := make(chan os.Signal, 1)                      //канал типа os.Signal
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT) // Запись в канал будет происходить когда процесс, в котором выполняется наше приложение, получит сигнал типа sigterm или sigint

	<-quit //строка для чтения из канала, которая будет блокировать выполнение главной горутины main

	logrus.Print("TodoApp Shutting Down")

	// вызовем два метода остановки сервера и закрытия всех соединений с БД. Это гарантирует нам, что мы закончим выполнение всех текущих операций перед выходом из приложения.
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig() //функция считывает значение конфига и записывает во внутренний объект вайпера, а возвращает ошибку
}
