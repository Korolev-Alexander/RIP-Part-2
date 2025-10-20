package main

import (
	"log"
	"net/http"
	"strings"

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
	smartDeviceAPI := apiHandlers.NewSmartDeviceAPIHandler(db)
	smartOrderAPI := apiHandlers.NewSmartOrderAPIHandler(db)
	orderItemAPI := apiHandlers.NewOrderItemAPIHandler(db)
	clientAPI := apiHandlers.NewClientAPIHandler(db)

	// Статические файлы
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Главная страница - сразу показываем устройства
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			handlers.SmartDevicesHandler(w, r)
			return
		}
		handlers.Show404Page(w, "Страница не найдена")
	})

	// HTML маршруты
	http.HandleFunc("/smart-devices", handlers.SmartDevicesHandler)
	http.HandleFunc("/smart-devices/", handlers.SmartDeviceDetailHandler)
	http.HandleFunc("/smart-cart", handlers.SmartCartHandler)
	http.HandleFunc("/smart-cart/add", handlers.AddToSmartCartHandler)
	http.HandleFunc("/smart-cart/delete", handlers.DeleteSmartCartHandler)
	http.HandleFunc("/smart-cart/count", handlers.GetSmartCartCountHandler)
	http.HandleFunc("/request/", handlers.RequestByIDHandler)

	// API маршруты - Smart Devices
	http.HandleFunc("/api/smart-devices", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			smartDeviceAPI.GetSmartDevices(w, r)
		case "POST":
			smartDeviceAPI.CreateSmartDevice(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Обработка всех /api/smart-devices/... маршрутов
	http.HandleFunc("/api/smart-devices/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		switch {
		case strings.Contains(path, "/image"):
			if r.Method == "POST" {
				smartDeviceAPI.UploadDeviceImage(w, r)
			} else if r.Method == "DELETE" {
				smartDeviceAPI.DeleteDeviceImage(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		default:
			// Обычные CRUD операции
			switch r.Method {
			case "GET":
				smartDeviceAPI.GetSmartDevice(w, r)
			case "PUT":
				smartDeviceAPI.UpdateSmartDevice(w, r)
			case "DELETE":
				smartDeviceAPI.DeleteSmartDevice(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}
	})

	// API маршруты - Smart Orders
	http.HandleFunc("/api/smart-orders/cart", smartOrderAPI.GetCart)
	http.HandleFunc("/api/smart-orders", smartOrderAPI.GetSmartOrders)

	// Обработка всех /api/smart-orders/... маршрутов
	http.HandleFunc("/api/smart-orders/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		switch {
		case strings.Contains(path, "/complete"):
			if r.Method == "PUT" {
				smartOrderAPI.CompleteSmartOrder(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		case strings.Contains(path, "/form"):
			if r.Method == "PUT" {
				smartOrderAPI.FormSmartOrder(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		default:
			// Обычные CRUD операции
			switch r.Method {
			case "GET":
				smartOrderAPI.GetSmartOrder(w, r)
			case "PUT":
				smartOrderAPI.UpdateSmartOrder(w, r)
			case "DELETE":
				smartOrderAPI.DeleteSmartOrder(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}
	})

	// API маршруты - Order Items (ИСПРАВЛЕННАЯ МАРШРУТИЗАЦИЯ)
	http.HandleFunc("/api/order-items/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" {
			orderItemAPI.UpdateOrderItem(w, r)
		} else if r.Method == "DELETE" {
			orderItemAPI.DeleteOrderItem(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// API маршруты - Clients
	http.HandleFunc("/api/clients/login", clientAPI.Login)
	http.HandleFunc("/api/clients/logout", clientAPI.Logout)
	http.HandleFunc("/api/clients/register", clientAPI.CreateClient)
	http.HandleFunc("/api/clients/update", clientAPI.UpdateClient)
	http.HandleFunc("/api/clients/", clientAPI.GetClient)
	http.HandleFunc("/api/clients", clientAPI.GetClients)

	log.Println("🚀 Сервер запущен на http://localhost:8080")
	log.Println("📱 HTML интерфейс доступен")
	log.Println("🔗 API доступно (22 метода)")

	log.Println("📦 Smart Devices API:")
	log.Println("   GET    /api/smart-devices              - список устройств")
	log.Println("   GET    /api/smart-devices/{id}         - устройство по ID")
	log.Println("   POST   /api/smart-devices              - создать устройство")
	log.Println("   PUT    /api/smart-devices/{id}         - обновить устройство")
	log.Println("   DELETE /api/smart-devices/{id}         - удалить устройство")
	log.Println("   POST   /api/smart-devices/{id}/image   - загрузить картинку")
	log.Println("   DELETE /api/smart-devices/{id}/image   - удалить картинку")

	log.Println("📋 Smart Orders API:")
	log.Println("   GET    /api/smart-orders/cart          - корзина")
	log.Println("   GET    /api/smart-orders               - список заявок")
	log.Println("   GET    /api/smart-orders/{id}          - заявка по ID")
	log.Println("   PUT    /api/smart-orders/{id}          - обновить заявку")
	log.Println("   PUT    /api/smart-orders/{id}/form     - сформировать заявку")
	log.Println("   PUT    /api/smart-orders/{id}/complete - завершить заявку")
	log.Println("   DELETE /api/smart-orders/{id}          - удалить заявку")

	log.Println("🛒 Order Items API:")
	log.Println("   PUT    /api/order-items/{deviceId}     - изменить количество")
	log.Println("   DELETE /api/order-items/{deviceId}     - удалить из заявки")

	log.Println("👥 Clients API:")
	log.Println("   GET    /api/clients                    - список клиентов")
	log.Println("   GET    /api/clients/{id}               - клиент по ID")
	log.Println("   POST   /api/clients/register           - регистрация")
	log.Println("   PUT    /api/clients/update             - обновить данные")
	log.Println("   POST   /api/clients/login              - аутентификация")
	log.Println("   POST   /api/clients/logout             - деавторизация")

	log.Println("🎯 Всего методов: 22")

	// ⚠️ ЭТА СТРОЧКА ОБЯЗАТЕЛЬНА! - запускает HTTP сервер
	http.ListenAndServe(":8080", nil)
}
