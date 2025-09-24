package handlers

import (
	"html/template"
	"net/http"
	"smartdevices/internal/models"
	"strconv"
)

var (
	tmplRequest = template.Must(template.ParseFiles("templates/layout.html", "templates/request.html"))
)

// GET /request - просмотр заявки через GORM
func RequestHandler(w http.ResponseWriter, r *http.Request) {
	// Временная реализация - вернем пустую страницу
	tmplRequest.ExecuteTemplate(w, "layout.html", map[string]interface{}{
		"Request": nil,
		"Items":   []interface{}{},
	})
}

// POST /request/add - добавление в заявку через GORM
func AddToRequestHandler(w http.ResponseWriter, r *http.Request) {
	// Временная заглушка
	http.Error(w, "Not implemented yet", http.StatusNotImplemented)
}

// POST /request/delete - логическое удаление через RAW SQL
func DeleteRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	requestID := r.FormValue("request_id")
	if requestID == "" {
		http.Error(w, "Request ID is required", http.StatusBadRequest)
		return
	}

	// ВЫПОЛНЯЕМ ТРЕБОВАНИЕ ТЗ: RAW SQL UPDATE
	sqlDB, err := db.DB() // ДОБАВИЛ получение sql.DB
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	_, err = sqlDB.Exec("UPDATE requests SET status = 'deleted' WHERE id = $1", requestID)
	if err != nil {
		http.Error(w, "Error deleting request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/devices", http.StatusFound)
}

// GET /request/count - количество товаров в корзине
func GetCartCountHandler(w http.ResponseWriter, r *http.Request) {
	var count int64
	// Логика подсчета товаров в корзине
	db.Model(&models.RequestService{}).Count(&count)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"count": ` + strconv.FormatInt(count, 10) + `}`))
}
