package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"smartdevices/internal/models"

	"gorm.io/gorm"
)

var (
	db                    *gorm.DB
	tmplSmartDevices      = template.Must(template.ParseFiles("templates/layout.html", "templates/smart_devices.html"))
	tmplSmartDeviceDetail = template.Must(template.ParseFiles("templates/layout.html", "templates/smart_device_detail.html"))
	tmplSmartCart         = template.Must(template.ParseFiles("templates/layout.html", "templates/smart_cart.html"))
	tmpl404               = template.Must(template.ParseFiles("templates/404.html"))
)

func Init(database *gorm.DB) {
	db = database
}

// Вспомогательная функция для получения количества товаров
func getSmartCartCount(clientID uint) int64 {
	var count int64
	db.Model(&models.OrderItem{}).
		Joins("JOIN smart_orders ON smart_orders.id = order_items.order_id").
		Where("smart_orders.client_id = ? AND smart_orders.status = ?", clientID, "draft").
		Count(&count)
	return count
}

// Вспомогательная функция для расчета общего трафика
func calculateTotalTraffic(orderID uint) float64 {
	var total float64

	db.Model(&models.OrderItem{}).
		Select("SUM(smart_devices.data_per_hour * order_items.quantity)").
		Joins("JOIN smart_devices ON smart_devices.id = order_items.device_id").
		Where("order_items.order_id = ?", orderID).
		Scan(&total)

	log.Printf("🔄 Расчет трафика для заявки %d: %.2f Кб/ч", orderID, total)
	return total
}

// Вспомогательная функция для показа 404 страницы
func Show404Page(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusNotFound)
	tmpl404.Execute(w, map[string]string{
		"ErrorMessage": message,
	})
}

// GET /request/{id} - просмотр заявки по ID
func RequestByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/request/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		Show404Page(w, "Неверный ID заявки")
		return
	}

	var order models.SmartOrder
	var items []models.OrderItem

	result := db.Preload("Client").First(&order, id)
	if result.Error != nil {
		Show404Page(w, "Заявка не найдена")
		return
	}

	if order.Status == "deleted" {
		Show404Page(w, "Заявка была удалена")
		return
	}

	db.Preload("Device").Where("order_id = ?", order.ID).Find(&items)

	order.TotalTraffic = calculateTotalTraffic(order.ID)

	err = tmplSmartCart.ExecuteTemplate(w, "layout.html", map[string]interface{}{
		"Request":   order,
		"Items":     items,
		"ShowCart":  false,
		"CartCount": getSmartCartCount(1),
	})

	if err != nil {
		log.Printf("❌ Template error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GET /smart-devices - поиск устройств через GORM
func SmartDevicesHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	var devices []models.SmartDevice
	query := db.Where("is_active = ?", true)

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	result := query.Find(&devices)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	err := tmplSmartDevices.ExecuteTemplate(w, "layout.html", map[string]interface{}{
		"Devices":   devices,
		"Search":    search,
		"ShowCart":  true,
		"CartCount": getSmartCartCount(1),
	})

	if err != nil {
		log.Printf("Template error in SmartDevicesHandler: %v", err)
	}
}

// GET /smart-devices/{id} - детальная страница устройства
func SmartDeviceDetailHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/smart-devices/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid device ID", http.StatusBadRequest)
		return
	}

	var device models.SmartDevice
	result := db.First(&device, id)
	if result.Error != nil {
		http.NotFound(w, r)
		return
	}

	log.Printf("📱 Device Detail - ID: %d, Name: %s, NamespaceURL: %s", device.ID, device.Name, device.NamespaceURL)

	err = tmplSmartDeviceDetail.ExecuteTemplate(w, "layout.html", map[string]interface{}{
		"Device":    device,
		"ShowCart":  false,
		"CartCount": getSmartCartCount(1),
	})

	if err != nil {
		log.Printf("❌ Template error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GET /smart-cart - просмотр корзины
func SmartCartHandler(w http.ResponseWriter, r *http.Request) {
	var order models.SmartOrder
	var items []models.OrderItem

	result := db.Preload("Client").Where("status = ? AND client_id = ?", "draft", 1).First(&order)

	if result.Error != nil {
		Show404Page(w, "Корзина пуста. Добавьте устройства из каталога.")
		return
	}

	db.Preload("Device").Where("order_id = ?", order.ID).Find(&items)

	order.TotalTraffic = calculateTotalTraffic(order.ID)

	log.Printf("📱 Загрузка корзины ID %d: %d товаров, трафик: %.2f Кб/ч",
		order.ID, len(items), order.TotalTraffic)

	err := tmplSmartCart.ExecuteTemplate(w, "layout.html", map[string]interface{}{
		"Request":   order,
		"Items":     items,
		"ShowCart":  false,
		"CartCount": getSmartCartCount(1),
	})

	if err != nil {
		log.Printf("❌ Template error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// POST /smart-cart/add - добавление в корзину
func AddToSmartCartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	deviceID := r.FormValue("device_id")
	if deviceID == "" {
		http.Error(w, "Device ID is required", http.StatusBadRequest)
		return
	}

	dID, err := strconv.Atoi(deviceID)
	if err != nil {
		http.Error(w, "Invalid device ID", http.StatusBadRequest)
		return
	}

	var order models.SmartOrder
	result := db.Where("status = ? AND client_id = ?", "draft", 1).First(&order)

	if result.Error != nil {
		order = models.SmartOrder{
			Status:   "draft",
			ClientID: 1,
			Address:  "ул. Примерная, д. 1, кв. 5",
		}
		db.Create(&order)
		log.Printf("📝 Создана новая корзина ID: %d", order.ID)
	}

	var existingOrderItem models.OrderItem
	findResult := db.Where("order_id = ? AND device_id = ?", order.ID, dID).First(&existingOrderItem)

	if findResult.Error == nil {
		existingOrderItem.Quantity++
		db.Save(&existingOrderItem)
		log.Printf("➕ Увеличено количество устройства %d в корзине %d: %d шт.", dID, order.ID, existingOrderItem.Quantity)
	} else {
		orderItem := models.OrderItem{
			OrderID:  order.ID,
			DeviceID: uint(dID),
			Quantity: 1,
		}
		db.Create(&orderItem)
		log.Printf("🆕 Добавлено устройство %d в корзину %d", dID, order.ID)
	}

	totalTraffic := calculateTotalTraffic(order.ID)
	log.Printf("📊 Общий трафик корзины %d: %.2f Кб/ч", order.ID, totalTraffic)

	http.Redirect(w, r, "/smart-cart", http.StatusSeeOther)
}

// POST /smart-cart/delete - удаление корзины через RAW SQL
func DeleteSmartCartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orderID := r.FormValue("order_id")
	if orderID == "" {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	_, err = sqlDB.Exec("UPDATE smart_orders SET status = 'deleted' WHERE id = $1", orderID)
	if err != nil {
		http.Error(w, "Error deleting order: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("🗑️ Deleted cart: id=%s", orderID)
	http.Redirect(w, r, "/smart-devices", http.StatusSeeOther)
}

// GET /smart-cart/count - количество товаров в корзине
func GetSmartCartCountHandler(w http.ResponseWriter, r *http.Request) {
	var count int64

	db.Model(&models.OrderItem{}).
		Joins("JOIN smart_orders ON smart_orders.id = order_items.order_id").
		Where("smart_orders.client_id = ? AND smart_orders.status = ?", 1, "draft").
		Count(&count)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"count": ` + strconv.FormatInt(count, 10) + `}`))
}
