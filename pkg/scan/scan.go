package scan

import (
	"fmt"
	"github.com/SayHelloLexa/searcher/pkg/crawler"
	"github.com/SayHelloLexa/searcher/pkg/crawler/spider"
	"github.com/SayHelloLexa/searcher/pkg/index"
)

// Scan - функция для сканирования сайтов + добавление документов
func Scan(idx *index.Index) (map[string]crawler.Document, error) {
	c := spider.New()
	urls := []string{"https://go.dev", "https://golang.org"}
	storage := make(map[string]crawler.Document)

	for i := range urls {
		d, err := c.Scan(urls[i], 2)
		if err != nil {
			return storage, fmt.Errorf("scaning error: %v", err)
		}

		addScanDocuments(d, storage, idx)
	}

	return storage, nil
}

// addScanDocuments - функция для добавления документов в обратный индекс и хранилище
// используется как компонент для Scan
func addScanDocuments(d []crawler.Document, storage map[string]crawler.Document, idx *index.Index) {
	for i, v := range d {
		v.ID = i
		idx.Add(v.Title, v.ID)
		storage[v.Title] = v
	}
}
