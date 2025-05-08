package API

import (
	"github.com/SayHelloLexa/searcher/pkg/scan"
	"net/http"
)

// jsonHeaderMiddleware - устанавливает тип контента для всех ответов в формате JSON
func jsonHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// scanDocMiddleware - сканирует документы и сохраняет их в хранилище, а также добавляет в индекс
func (api *API) scanDocMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(api.storage) != 0 {
			next.ServeHTTP(w, r)
			return
		}

		var err error
		api.storage, err = scan.Scan(&api.idx)
		if err != nil {
			http.Error(w, "ошибка при сканировании в scanDocMiddleware", http.StatusInternalServerError)
		}
		next.ServeHTTP(w, r)
	})
}
