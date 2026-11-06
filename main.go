package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"sync"
)

type Feedback struct {
	ID              int
	Author, Message string
}
type PageData struct {
	Feedbacks []Feedback
	Total     int
}

var feedbacks []Feedback
var nextID int
var mutex sync.Mutex
var page = template.Must(template.ParseFiles("templates/index.html"))

// Показываем отзывы и их количество.
func homeHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		http.NotFound(responseWriter, request)
		return
	}
	if request.Method != http.MethodGet {
		http.Error(responseWriter, "Method not allowed", 405)
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	page.Execute(responseWriter, PageData{feedbacks, len(feedbacks)})
}

// Проверяем форму и сохраняем новый отзыв.
func feedbackHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(responseWriter, "Method not allowed", 405)
		return
	}
	author := strings.TrimSpace(request.FormValue("author"))
	message := strings.TrimSpace(request.FormValue("message"))
	if author == "" || message == "" {
		http.Error(responseWriter, "Имя и текст отзыва обязательны", 400)
		return
	}
	mutex.Lock()
	nextID = nextID + 1
	feedbacks = append(feedbacks, Feedback{nextID, author, message})
	mutex.Unlock()
	http.Redirect(responseWriter, request, "/", http.StatusSeeOther)
}

func aboutHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(responseWriter, "Method not allowed", 405)
		return
	}
	template.Must(template.ParseFiles("templates/about.html")).Execute(responseWriter, nil)
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/feedback", feedbackHandler)
	http.HandleFunc("/about", aboutHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
