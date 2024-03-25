package main

import (
	"github.com/larinkir/command_api/internal/cash"
	"github.com/larinkir/command_api/internal/config"
	"github.com/larinkir/command_api/internal/http-server/handlers"
	"github.com/larinkir/command_api/internal/http-server/server"
	"github.com/larinkir/command_api/internal/storage/psql"
	"log"
	"net/http"
)

func main() {

	//Получение config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	//Подключение к ДБ
	storage, err := psql.Connect(cfg.StorageData)
	if err != nil {
		log.Fatal(err)
	}

	//Инициализация cash
	cash.NewCommandProcess()

	//Инициализация роутеров
	http.HandleFunc("/", handlers.FirstPage())
	http.HandleFunc("/save", handlers.HandleSaver(storage))               //POST. Сохранение команды.
	http.HandleFunc("/getall", handlers.HandleGetterAllCommands(storage)) //GET. Получение списка команда.
	http.HandleFunc("/get/", handlers.HandleGetterCommand(storage))       //POST. Получение команды по Id + запуск.
	http.HandleFunc("/delete/", handlers.HandleDeleter(storage))          //DELETE. Удаление команды по Id.
	http.HandleFunc("/stop/", handlers.HandleStopperCommand(storage))     //GET. Остановка команды по Id.

	//Запуск сервера
	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}

}
