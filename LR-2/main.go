package main

import (
	"log"
	"net/http"

	apiHandlers "smartdevices/internal/api/handlers"
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

	// Инициализация HTML handlers с передачей DB
	handlers.Init(db)

	// Инициализация API handlers
	serviceAPI := apiHandlers.NewServiceAPIHandler(db)
	requestAPI := apiHandlers.NewRequestAPIHandler(db)
	requestServiceAPI := apiHandlers.NewRequestServiceAPIHandler(db)
	userAPI := apiHandlers.NewUserAPIHandler(db) // ← ДОБАВИТЬ ЭТУ СТРОКУ

	// Статические файлы
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Редирект с корневого пути на страницу устройств
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/smart-devices", http.StatusSeeOther)
			return
		}
		handlers.Show404Page(w, "Страница не найдена")
	})

	// HTML маршруты по ТЗ
	http.HandleFunc("/smart-devices", handlers.SmartDevicesHandler)
	http.HandleFunc("/smart-devices/", handlers.SmartDeviceDetailHandler)
	http.HandleFunc("/smart-cart", handlers.SmartCartHandler)
	http.HandleFunc("/smart-cart/add", handlers.AddToSmartCartHandler)
	http.HandleFunc("/smart-cart/delete", handlers.DeleteSmartCartHandler)
	http.HandleFunc("/smart-cart/count", handlers.GetSmartCartCountHandler)
	http.HandleFunc("/request/", handlers.RequestByIDHandler)

	// API маршруты
	http.HandleFunc("/api/services", serviceAPI.GetServices)
	http.HandleFunc("/api/services/", serviceAPI.GetService)
	http.HandleFunc("/api/requests/cart", requestAPI.GetCart)
	http.HandleFunc("/api/requests", requestAPI.GetRequests)
	http.HandleFunc("/api/requests/", requestAPI.GetRequest)
	http.HandleFunc("/api/request-services/", requestServiceAPI.UpdateRequestService)

	// ДОБАВИТЬ ЭТИ МАРШРУТЫ ДЛЯ ПОЛЬЗОВАТЕЛЕЙ:
	http.HandleFunc("/api/users", userAPI.GetUsers)
	http.HandleFunc("/api/users/", userAPI.GetUser)

	log.Println("🚀 Сервер запущен на http://localhost:8080")
	log.Println("📱 API доступно на http://localhost:8080/api/services")
	http.ListenAndServe(":8080", nil)
}
