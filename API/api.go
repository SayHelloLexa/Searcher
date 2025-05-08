package API

import (
	"encoding/json"
	"github.com/SayHelloLexa/searcher/pkg/crawler"
	"github.com/SayHelloLexa/searcher/pkg/index"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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
	api.Router.HandleFunc("/api/v1/docs/{id}", api.createDoc).Methods(http.MethodPost)
	api.Router.HandleFunc("/api/v1/docs/{id}", api.updDocById).Methods(http.MethodPut, http.MethodPatch)
}

// getDocStorage - получить все документы
func (api *API) getDocStorage(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(api.storage)
	if err != nil {
		http.Error(w, "ошибка при сериализации в JSON формат", http.StatusInternalServerError)
	}
}

// getDocById - получить документ по id
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

// createDoc - создает документ и добавляет в хранилище
func (api *API) createDoc(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if _, ok := api.storage[id]; ok {
		http.Error(w, "в хранилище уже присутствует документ с данным id, попробуйте присвоить другое", http.StatusInternalServerError)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var doc crawler.Document
	err := decoder.Decode(&doc)
	if err != nil {
		http.Error(w, "ошибка при диссереализации тела запроса", http.StatusInternalServerError)
	}

	api.storage[id] = doc
	w.WriteHeader(http.StatusCreated)
}

// updDocById - обновляет документ в хранилище по id
func (api *API) updDocById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if _, ok := api.storage[id]; !ok {
		http.Error(w, "документ не найден в хранилище storage", http.StatusInternalServerError)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var doc crawler.Document
	err := decoder.Decode(&doc)
	if err != nil {
		http.Error(w, "ошибка при диссереализации тела запроса", http.StatusInternalServerError)
	}

	if strconv.Itoa(doc.ID) != id {
		http.Error(w, "id в теле запроса не совпадает с id в хранилище storage и URL", http.StatusInternalServerError)
	}

	api.storage[id] = doc
	w.WriteHeader(http.StatusCreated)
}

// deleteDocById - удаляет документ из хранилища по id
func (api *API) deleteDocById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	delete(api.storage, id)
}
