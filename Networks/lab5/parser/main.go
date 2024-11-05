package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/net/html/charset"
)

type RSS struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Items []Item `xml:"item"`
}

type Item struct {
	Title string `xml:"title"`
	Body  string `xml:"description"`
}

func transliterate(text string) string {
	var translitMap = map[rune]string{
		'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "yo", 'ж': "zh", 'з': "z",
		'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m", 'н': "n", 'о': "o", 'п': "p", 'р': "r",
		'с': "s", 'т': "t", 'у': "u", 'ф': "f", 'х': "kh", 'ц': "ts", 'ч': "ch", 'ш': "sh", 'щ': "shch",
		'ъ': "", 'ы': "y", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya",
		'А': "A", 'Б': "B", 'В': "V", 'Г': "G", 'Д': "D", 'Е': "E", 'Ё': "Yo", 'Ж': "Zh", 'З': "Z",
		'И': "I", 'Й': "Y", 'К': "K", 'Л': "L", 'М': "M", 'Н': "N", 'О': "O", 'П': "P", 'Р': "R",
		'С': "S", 'Т': "T", 'У': "U", 'Ф': "F", 'Х': "Kh", 'Ц': "Ts", 'Ч': "Ch", 'Ш': "Sh", 'Щ': "Shch",
		'Ъ': "", 'Ы': "Y", 'Ь': "", 'Э': "E", 'Ю': "Yu", 'Я': "Ya",
	}

	result := ""
	for _, char := range text {
		if translitChar, ok := translitMap[char]; ok {
			result += translitChar
		} else {
			result += string(char)
		}
	}
	return result
}

func main() {
	resp, err := http.Get("https://news.rambler.ru/rss/Ivanovo/")
	if err != nil {
		log.Fatalf("Ошибка запроса к RSS: %v", err)
	}
	defer resp.Body.Close()

	r, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		log.Fatalf("Ошибка создания декодера: %v", err)
	}

	data, err := io.ReadAll(r)
	if err != nil {
		log.Fatalf("Ошибка чтения RSS данных: %v", err)
	}

	var rss RSS
	if err := xml.Unmarshal(data, &rss); err != nil {
		log.Fatalf("Ошибка парсинга RSS: %v", err)
	}

	dsn := "iu9networkslabs:Je2dTYr6@tcp(students.yss.su)/iu9networkslabs?charset=utf8mb4"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Ошибка проверки подключения к базе данных: %v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Ошибка начала транзакции: %v", err)
	}

	_, err = tx.Exec("TRUNCATE TABLE iu9Krasnobaev")
	if err != nil {
		tx.Rollback()
		log.Fatalf("Ошибка очистки таблицы: %v", err)
	}

	for _, item := range rss.Channel.Items {
		titleTranslit := transliterate(item.Title)
		bodyTranslit := transliterate(item.Body)

		_, err := tx.Exec("INSERT INTO iu9Krasnobaev (title, body) VALUES (?, ?)", titleTranslit, bodyTranslit)
		if err != nil {
			tx.Rollback()
			log.Printf("Ошибка вставки записи: %v", err)
			return
		} else {
			fmt.Printf("Успешно добавлена запись: %s\n", titleTranslit)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("Ошибка коммита транзакции: %v", err)
	}

	fmt.Println("Все данные успешно обновлены.")

}
