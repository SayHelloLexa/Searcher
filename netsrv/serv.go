package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/SayHelloLexa/searcher/pkg/crawler"
	"github.com/SayHelloLexa/searcher/pkg/crawler/spider"
)

// Scan - функция для сканирования сайтов
func scan() (map[string]crawler.Document, error) {
	c := spider.New()
	urls := []string{"https://go.dev", "https://golang.org"}
	storage := make(map[string]crawler.Document)

	for i := range urls {
		d, err := c.Scan(urls[i], 2)
		if err != nil {
			return storage, fmt.Errorf("scaning error: %v", err)
		}

		for _, v := range d {
			storage[v.Title] = v
		}
	}

	return storage, nil
}

func handler(conn net.Conn, storage map[string]crawler.Document) {
	defer conn.Close()

	for {
		// Читаем запрос клиента
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			log.Fatal(err)
		}

		// Обрабатываем запрос клиента
		input, _, _ := strings.Cut(string(buf[:n]), "\n")
		conn.SetWriteDeadline(time.Now().Add(time.Second * 10))

		// Проверяем запрос на пустоту
		if input == "" || input == "\n" {
			_, err = conn.Write([]byte("empty search\n"))
			if err != nil {
				conn.Write([]byte("ERROR: failed to send response\n"))
				log.Printf("Соединение разорвано %v", err)
				break
			}
			continue
		}

		// Обработка разрыва соединения
		if input == "exit" {
			conn.Write([]byte("Connection closed\n"))
			break
		}

		// Отправляем ответ клиенту
		conn.SetWriteDeadline(time.Now().Add(time.Second * 5))
		for _, v := range storage {
			if strings.Contains(strings.ToLower(v.Title), strings.ToLower(input)) {
				_, err := conn.Write([]byte(v.URL))
				if err != nil {
					continue
				}

				_, err = conn.Write([]byte("\n"))
				if err != nil {
					continue
				}
			}
		}
	}
}

func main() {
	// Занимаем порт 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	storage, err := scan()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server is listening...")
	fmt.Println(storage)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handler(conn, storage)
	}
}
