package middleware

import (
	"net/http"
)

// CORSMiddleware добавляет необходимые CORS заголовки для работы с cookies
func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем запросы с фронтенда (для разработки)
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")

		// Разрешаем необходимые заголовки
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		// Разрешаем методы
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Разрешаем отправку cookies
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Обрабатываем preflight запросы
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Продолжаем выполнение цепочки middleware
		next(w, r)
	}
}
