package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/go-chi/chi"
	"golang.org/x/net/html"
)

func main() {
	routeAndRun()
}

var url string = "https://finance.rambler.ru/currencies/"

func routeAndRun() {
	r := chi.NewRouter()
	r.Get("/", home)
	r.Get("/exchange", getExchangeRate)
	http.ListenAndServe("localhost:8080", r)
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "you can go to /exchange")
}

type Currency struct {
	Code, Nominal, Name, Rate, Change, ChangePercent string
}

func getExchangeRate(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("can't get html page %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Printf("can't parse html page %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	currencyList := []Currency{}
	processAllCurrency(doc, &currencyList)
	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		log.Printf("can't parse template %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, currencyList)
	if err != nil {
		log.Printf("execution template error %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Println("exhange rate was given succesfully")
}

func processAllCurrency(n *html.Node, cucurrencyList *[]Currency) {
	if n.Type == html.ElementNode && n.Data == "tr" {
		processNode(n, cucurrencyList)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		processAllCurrency(c, cucurrencyList)
	}
}

func processNode(n *html.Node, cucurrencyList *[]Currency) {
	newCurrency := Currency{}
	tdIndex := 0
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "td" {
			switch tdIndex {
			case 0:
				newCurrency.Code = extractTextFromNode(c)
			case 1:
				newCurrency.Nominal = extractTextFromNode(c)
			case 2:
				newCurrency.Name = extractTextFromNode(c)
			case 3:
				newCurrency.Rate = extractTextFromNode(c)
			case 4:
				newCurrency.Change = extractTextFromNode(c)
			case 5:
				newCurrency.ChangePercent = extractTextFromNode(c)
			}
			tdIndex += 1
		}
	}
	*cucurrencyList = append(*cucurrencyList, newCurrency)
}

func extractTextFromNode(n *html.Node) string {
	var text string
	var extractText func(*html.Node)
	extractText = func(n *html.Node) {
		if n.Type == html.TextNode {
			text += n.Data
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractText(c)
		}
	}
	extractText(n)
	return text
}
