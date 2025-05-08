package API

import (
	"encoding/json"
	"github.com/SayHelloLexa/searcher/pkg/crawler"
	"github.com/SayHelloLexa/searcher/pkg/index"
	"github.com/gorilla/mux"
	"net/http"
)

type API struct {
	Router  *mux.Router
	storage map[string]crawler.Document
	idx     index.Index
}

// New - конструктор объекта API
func New() *API {
	api := &API{
		Router:  mux.NewRouter(),
		storage: make(map[string]crawler.Document),
		idx:     index.New(),
	}
	api.endpoints()

	return api
}

// endpoints - метод регистрации эндпоинтов API
func (api *API) endpoints() {
	api.Router.Use(jsonHeaderMiddleware, api.scanDocMiddleware)

	api.Router.HandleFunc("/api/v1/docs/{id}", api.getDocById).Methods(http.MethodGet)
	api.Router.HandleFunc("/api/v1/docs", api.getDocStorage).Methods(http.MethodGet)
	api.Router.HandleFunc("/api/v1/docs/{id}", api.deleteDocById).Methods(http.MethodDelete)
	api.Router.HandleFunc("/api/v1/{id}", api.createDoc).Methods(http.MethodPost)
}

func (api *API) getDocStorage(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(api.storage)
	if err != nil {
		http.Error(w, "ошибка при сериализации в JSON формат", http.StatusInternalServerError)
	}
}

func (api *API) getDocById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	v, ok := api.storage[id]
	if !ok {
		http.Error(w, "документ не найден в хранилище storage", http.StatusInternalServerError)
	}

	encoder := json.NewEncoder(w)
	err := encoder.Encode(v)
	if err != nil {
		http.Error(w, "ошибка при сериализации в JSON формат", http.StatusInternalServerError)
	}
}

func (api *API) createDoc(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	decoder := json.NewDecoder(r.Body)
	var doc crawler.Document
	err := decoder.Decode(&doc)
	if err != nil {
		http.Error(w, "ошибка при диссереализации тела запроса", http.StatusInternalServerError)
	}

	api.storage[id] = doc
	w.WriteHeader(http.StatusCreated)
}

func (api *API) deleteDocById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	delete(api.storage, id)
}
