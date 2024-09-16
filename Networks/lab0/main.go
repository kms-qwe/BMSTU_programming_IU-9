package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/SlyMarbo/rss"
)

func HomeRouterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //анализ аргументов,

	rssObject, err := rss.Fetch("https://vmo24.ru/rss")
	if err != nil {
		fmt.Println(err)
		panic("failed")
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<!DOCTYPE html><html><head><title>RSS Feed</title></head><body>")

	fmt.Fprintf(w, "<h1>Title:</h1><p><strong>%s</strong></p>", rssObject.Title)
	fmt.Fprintf(w, "<h1>Description:</h1><p><strong>%s</strong></p>", rssObject.Description)
	fmt.Fprintf(w, "<h1>Number of Items:</h1><p><strong>%d</strong></p>", len(rssObject.Items))

	for v := range rssObject.Items {
		item := rssObject.Items[v]
		fmt.Fprintf(w, "<hr>")
		fmt.Fprintf(w, "<h2>Item Number:</h2><p><strong>%d</strong></p>", v)
		fmt.Fprintf(w, "<h3>Title:</h3><p><strong>%s</strong></p>", item.Title)
		fmt.Fprintf(w, "<h3>Link:</h3><p><a href='%s'>%s</a></p>", item.Link, item.Link)
		fmt.Fprintf(w, "<h3>Time:</h3><p><strong>%s</strong></p>", item.Date)
	}

	fmt.Fprintf(w, "</body></html>")
}

func main() {
	http.HandleFunc("/", HomeRouterHandler)  // установим роутер
	err := http.ListenAndServe("185.102.139.168:9000", nil) // задаем слушать порт
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
