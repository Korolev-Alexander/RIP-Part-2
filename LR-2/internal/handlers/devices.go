package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"smartdevices/internal/models"

	"gorm.io/gorm"
)

var (
	db               *gorm.DB
	tmplDevices      = template.Must(template.ParseFiles("templates/layout.html", "templates/devices.html"))
	tmplDeviceDetail = template.Must(template.ParseFiles("templates/layout.html", "templates/device_detail.html")) // ДОБАВИЛ
)

func Init(database *gorm.DB) {
	db = database
}

// GET /devices - поиск услуг через GORM
func DevicesHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	var services []models.Service
	query := db.Where("is_active = ?", true)

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	result := query.Find(&services)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	tmplDevices.ExecuteTemplate(w, "layout.html", map[string]interface{}{
		"Devices": services,
		"Search":  search,
	})
}

// ДОБАВИЛ ЭТОТ ФУНКЦИЮ - GET /devices/{id}
func DeviceDetailHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID из URL (например: /devices/1)
	idStr := r.URL.Path[len("/devices/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid device ID", http.StatusBadRequest)
		return
	}

	var service models.Service
	result := db.First(&service, id)
	if result.Error != nil {
		http.NotFound(w, r)
		return
	}

	tmplDeviceDetail.ExecuteTemplate(w, "layout.html", map[string]interface{}{
		"Device": service,
	})
}
