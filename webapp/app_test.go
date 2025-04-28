package main

import (
	"encoding/json"
	"github.com/SayHelloLexa/searcher/pkg/crawler"
	"github.com/SayHelloLexa/searcher/pkg/index"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIdxHandler(t *testing.T) {
	testIdx := index.New()
	testIdx.Add("test", 1)

	// Генерим запрос
	req := httptest.NewRequest(http.MethodGet, "/index", nil)

	rr := httptest.NewRecorder()

	// Вызываем обработчик
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(testIdx)
	})
	handler.ServeHTTP(rr, req)

	// Проверяем статус
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Получили: %d, ожидали: %d", status, http.StatusOK)
	}

	var response index.Index
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Errorf("Ошибка при декодировании: %v", err)
	}

	// map[string][]int - только 1 элемент у нас там
	if len(response) != 1 {
		t.Errorf("Ожидали 1 элемент, получили: %d", len(response))
	}
}

func TestDocsHandler(t *testing.T) {
	testStorage := map[string]crawler.Document{
		"test document": {
			ID:    1,
			Title: "Test document",
			URL:   "http://test.com",
			Body:  "Lorem ipsum dolores",
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/docs", nil)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(testStorage)
	})
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Получили: %d, ожидали: %d", status, http.StatusOK)
	}

	var response map[string]crawler.Document
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Errorf("Ошибка при декодировании: %v", err)
	}

	if len(response) != 1 {
		t.Errorf("Ожидали 1 элемент, получили: %d", len(response))
	}
}
