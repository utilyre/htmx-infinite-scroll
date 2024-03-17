package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

//go:embed views/*.html
var viewsFS embed.FS

type Article struct {
	Title   string
	Content string
}

var articles []Article

func init() {
	n := 37
	articles = make([]Article, n)
	for i := range n {
		articles[i] = Article{
			Title:   fmt.Sprintf("Article %d", i+1),
			Content: "Nothing",
		}
	}
}

type ArticlesView struct {
	Articles []Article
	NextPage int
}

func main() {
	views, err := template.ParseFS(viewsFS, "views/*.html")
	if err != nil {
		log.Fatalln(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"ArticlesView": ArticlesView{
				Articles: articles[:10],
				NextPage: 2,
			},
		}

		w.Header().Set("Content-Type", "text/html")
		if err := views.ExecuteTemplate(w, "home", resp); err != nil {
			log.Println("failed to run GET /:", err)
		}
	})

	mux.HandleFunc("GET /articles", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)

		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil || page <= 0 {
			page = 1
		}

		numArticles := len(articles)
		start := min(numArticles, 10*(page-1))
		end := min(numArticles, 10*page)
		var nextPage int
		if end < numArticles {
			nextPage = page + 1
		}

		resp := ArticlesView{
			Articles: articles[start:end],
			NextPage: nextPage,
		}

		w.Header().Set("Content-Type", "text/html")
		if err := views.ExecuteTemplate(w, "articles", resp); err != nil {
			log.Println("failed to run GET /articles:", err)
		}
	})

	log.Println("Listening on http://localhost:3000")
	log.Fatalln(http.ListenAndServe(":3000", mux))
}
