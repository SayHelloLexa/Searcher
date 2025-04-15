package jsonutils

import (
	"fmt"
	"os"
	"strings"
)

// Преобразует URL в строку для создания имени файла
func UrlMap(url string) string {
	url = strings.Map(func(r rune) rune {
		switch r {
		case 'h', 't', 'p', 's', ':', '/':
				return -1 // Удалить
		case '.':
				return '-' // Заменить на '-'
		default:
				return r // Оставить как есть
		}
	}, url)

	return url
}

// Функция создает директорию для хранения
// JSON результатов сканирования страниц
func CreateDir(url string) (string, error) {
	url = UrlMap(url)

	dir := "../../JSON/"
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return "", fmt.Errorf("creating directory error: %v", err)
	}

	filepath := dir + url + ".JSON"
	file, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("creating file error: %v", err)
	}
	defer file.Close()

	return filepath, nil
}

// Функция проверяет существование файла
func IsExist(url string) bool {
    url = "../../JSON/" + UrlMap(url) + ".JSON"

    _, err := os.Stat(url)
    if err == nil {
        return true
    }
    if os.IsNotExist(err) {
        return false
    }
    return false 
}