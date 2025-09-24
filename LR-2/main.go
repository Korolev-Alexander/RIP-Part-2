package main

import (
	"log"
	"net/http"

	"smartdevices/internal/handlers"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Подключение к PostgreSQL через GORM
	dsn := "host=localhost user=root password=root dbname=RIP port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	// Инициализация handlers с передачей DB
	handlers.Init(db)

	// Статические файлы
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Маршруты по ТЗ
	http.HandleFunc("/devices", handlers.DevicesHandler)              // GET - поиск услуг (ORM)
	http.HandleFunc("/devices/", handlers.DeviceDetailHandler)        // GET - детали услуги (ORM)
	http.HandleFunc("/request", handlers.RequestHandler)              // GET - просмотр заявки (ORM)
	http.HandleFunc("/request/add", handlers.AddToRequestHandler)     // POST - добавить в заявку (ORM)
	http.HandleFunc("/request/delete", handlers.DeleteRequestHandler) // POST - удалить заявку (SQL UPDATE)

	log.Println("🚀 Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
