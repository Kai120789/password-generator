package app

import (
	"fmt"
	"net/http"
	"os"
	"password_generator/internal/service"
	"password_generator/internal/storage"
	"password_generator/internal/transport/http/handler"
	"password_generator/internal/transport/http/router"
	"password_generator/pkg/logger"

	"go.uber.org/zap"
)

func StartApp() {
	// logger
	zapLog, err := logger.New("debug")
	if err != nil {
		zap.S().Fatalf("Error logger init", zap.Error(err))
	}

	log := zapLog.ZapLogger

	filePath := "./passwords.json"

	// check is file exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer file.Close()
		fmt.Println("Файл успешно создан:", filePath)

	} else {
		fmt.Println("Файл уже существует:", filePath)
	}

	stor := storage.New(log, filePath)

	serv := service.New(stor)

	hand := handler.New(serv, log)

	r := router.New(&hand)

	log.Info("starting server", zap.String("address", "localhost:8080"))

	srv := &http.Server{
		Addr:    "localhost:8080",
		Handler: r,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", zap.Error(err))
	}
}
