package API

import (
	"github.com/SayHelloLexa/searcher/pkg/crawler"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetDocStorage(t *testing.T) {
	api := New()
	api.storage["1"] = crawler.Document{
		ID:    1,
		URL:   "example1.com",
		Title: "Example1",
		Body:  "Lorem ipsum dolor apsem",
	}
	api.storage["2"] = crawler.Document{
		ID:    2,
		URL:   "example2.com",
		Title: "Example2",
		Body:  "Lorem ipsum dolor apsem",
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/docs", nil)
	rr := httptest.NewRecorder()

	api.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("ожидали статус: %v, получили: %v", http.StatusOK, rr.Code)
	}

	if len(api.storage) != 2 {
		t.Errorf("ожидали: %v, получили: %v", 2, len(api.storage))
	}
}

func TestGetDocById(t *testing.T) {
	api := New()
	api.storage["1"] = crawler.Document{
		ID:    1,
		URL:   "example1.com",
		Title: "Example1",
		Body:  "Lorem ipsum dolor apsem",
	}
	api.storage["2"] = crawler.Document{
		ID:    2,
		URL:   "example2.com",
		Title: "Example2",
		Body:  "Lorem ipsum dolor apsem",
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/docs/1", nil)
	rr := httptest.NewRecorder()

	api.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("ожидали статус: %v, получили: %v", http.StatusOK, rr.Code)
	}

	expected := `{"ID":1,"URL":"example1.com","Title":"Example1","Body":"Lorem ipsum dolor apsem"}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("тело ответа не соответствует ожиданиям")
		t.Logf("ожидали: %s, получили: %s", expected, rr.Body.String())
	}
}

func TestDeleteDocById(t *testing.T) {
	api := New()
	api.storage["1"] = crawler.Document{
		ID:    1,
		URL:   "example1.com",
		Title: "Example1",
		Body:  "Lorem ipsum dolor apsem",
	}
	api.storage["2"] = crawler.Document{
		ID:    2,
		URL:   "example2.com",
		Title: "Example2",
		Body:  "Lorem ipsum dolor apsem",
	}

	initLen := len(api.storage)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/docs/2", nil)
	rr := httptest.NewRecorder()

	api.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("ожидали статус: %v, получили: %v", http.StatusOK, rr.Code)
	}

	if len(api.storage) != initLen-1 {
		t.Errorf("ожидали длину хранилища: %v, получили: %v", initLen-1, len(api.storage))
	}
}
