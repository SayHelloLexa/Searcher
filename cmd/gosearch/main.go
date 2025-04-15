package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}

		_, err = conn.Write([]byte(input))
		if err != nil {
			log.Fatal(err)
		}

		var buf bytes.Buffer
		tmp := make([]byte, 1024)

		conn.SetReadDeadline(time.Now().Add(time.Second))
		for {
			n, err := conn.Read(tmp)
			if err != nil {
				var netErr net.Error
				if errors.As(err, &netErr) && netErr.Timeout() {
					break
				}
				if err == io.EOF {
					break
				}
				log.Fatal("Ошибка чтения:", err)
			}
			buf.Write(tmp[:n])
		}

		if buf.String() == "Connection closed\n" {
			fmt.Println("Сервер закрыл соединение")
			break
		}

		if buf.Len() == 0 {
			fmt.Println("Сервер не завершил сканирование, или ответа не существует")
		} else {
			fmt.Println("Ответ сервера:\n", buf.String())
		}
	}
}
