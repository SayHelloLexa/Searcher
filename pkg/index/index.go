package index

import (
	"slices"
	"strings"
)

// Тип "Обратный индекс"
type Index map[string][]int

// Создание нового экземпляра обратного индекса
func New() Index {
	return make(Index)
}

/*
	Подаем на вход Title и ID документа, разбиваем
	Title на слова, проходимся по массиву слов и добавляем в 
	обратный индекс слово в качестве ключа и ID документа в качестве значения
*/
func (d *Index) Add(s string, id int) {
	words := strings.Fields(strings.ToLower(s))
	
	// Добавляем слово в качестве ключа и ID документа в качестве значения
	for _, v := range words {
		(*d)[v] = append((*d)[v], id) 
	}

	// Cортируем массив документов по ID
	for _, v := range *d {
		slices.Sort(v)
	}
}

/*
	Ищем массивы индексов, которые соответствуют
	флагу	
*/
func (d *Index) Search(f string) []int {
	return (*d)[strings.ToLower(f)]
}
