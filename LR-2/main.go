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

	// Инициализация HTML handlers с передачей DB
	handlers.Init(db)

	// Инициализация API handlers (пока закомментировано)
	/*
	   serviceAPI := apiHandlers.NewServiceAPIHandler(db)
	   requestAPI := apiHandlers.NewRequestAPIHandler(db)
	   requestServiceAPI := apiHandlers.NewRequestServiceAPIHandler(db)
	   userAPI := apiHandlers.NewUserAPIHandler(db)
	*/

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

	// API маршруты - ЗАКОММЕНТИРОВАНО для обновления
	/*
	   http.HandleFunc("/api/services", func(w http.ResponseWriter, r *http.Request) {
	       switch r.Method {
	       case "GET":
	           serviceAPI.GetServices(w, r)
	       case "POST":
	           serviceAPI.CreateService(w, r)
	       default:
	           http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	       }
	   })
	   http.HandleFunc("/api/services/", func(w http.ResponseWriter, r *http.Request) {
	       switch r.Method {
	       case "GET":
	           serviceAPI.GetService(w, r)
	       case "PUT":
	           serviceAPI.UpdateService(w, r)
	       case "DELETE":
	           serviceAPI.DeleteService(w, r)
	       default:
	           http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	       }
	   })
	   http.HandleFunc("/api/requests/cart", requestAPI.GetCart)
	   http.HandleFunc("/api/requests", requestAPI.GetRequests)
	   http.HandleFunc("/api/requests/", requestAPI.GetRequest)
	   http.HandleFunc("/api/request-services/", requestServiceAPI.UpdateRequestService)

	   // Маршруты для пользователей
	   http.HandleFunc("/api/users", userAPI.GetUsers)
	   http.HandleFunc("/api/users/", userAPI.GetUser)
	*/

	log.Println("🚀 Сервер запущен на http://localhost:8080")
	log.Println("📱 HTML интерфейс доступен")
	log.Println("⚡ API временно отключено для обновления")
	http.ListenAndServe(":8080", nil)
}
