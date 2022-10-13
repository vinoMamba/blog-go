package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/vinoMamba/goblog/pkg/database"
	"github.com/vinoMamba/goblog/pkg/logger"
	"github.com/vinoMamba/goblog/pkg/route"
)

var router *mux.Router
var db *sql.DB

type Article struct {
	ID    int64
	Title string
	Body  string
}

func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	// 从数据库中读取数据
	sqlStatement := `SELECT * FROM articles WHERE id=$1`
	row := db.QueryRow(sqlStatement, id)
	var article Article
	err := row.Scan(&article.ID, &article.Title, &article.Body)
	switch err {
	case sql.ErrNoRows:
		fmt.Fprint(w, "没有这篇文章")
	case nil:
		tmpl, err := template.ParseFiles("./resources/views/articles/show.html")
		if err != nil {
			log.Fatal(err)
		}
		tmpl.Execute(w, article)
	default:
		logger.LogError(err)
	}
}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	sqlStatement := `SELECT * FROM articles`
	rows, err := db.Query(sqlStatement)
	logger.LogError(err)
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var article Article
		err = rows.Scan(&article.ID, &article.Title, &article.Body)
		logger.LogError(err)
		articles = append(articles, article)
	}
	err = rows.Err()
	logger.LogError(err)
	tmpl, err := template.ParseFiles("./resources/views/articles/index.html")
	if err != nil {
		panic(err)
	}
	tmpl.Execute(w, articles)
}

type ArticlesFormData struct {
	Title  string
	Body   string
	URL    *url.URL
	Errors map[string]string
}

func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")
	errors := make(map[string]string)
	if title == "" {
		errors["title"] = "标题不可为空"
	} else if len(title) < 3 || len(title) > 40 {
		errors["title"] = "标题长度需介于 3-40"
	}
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if len(body) < 10 {
		errors["body"] = "内容长度需大于或等于 10 个字节"
	}
	if len(errors) == 0 {
		lastInsertID, err := saveArticleToDB(title, body)
		if lastInsertID > 0 {
			fmt.Fprint(w, "插入成功，ID 为"+strconv.FormatInt(lastInsertID, 10))
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		storeURL, _ := router.Get("articles.store").URL()
		data := ArticlesFormData{
			Title:  title,
			Body:   body,
			URL:    storeURL,
			Errors: errors,
		}
		tmpl, err := template.ParseFiles("./resources/views/articles/create.html")
		if err != nil {
			panic(err)
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			panic(err)
		}
	}
}

func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {
	storeURL, _ := router.Get("articles.store").URL()
	data := ArticlesFormData{
		Title:  "",
		Body:   "",
		URL:    storeURL,
		Errors: nil,
	}
	tmpl, err := template.ParseFiles("./resources/views/articles/create.html")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func forceHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

func removeTrilingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		next.ServeHTTP(w, r)
	})
}

func saveArticleToDB(title string, body string) (int64, error) {
	sqlStatement := `INSERT INTO articles (title, body) VALUES ($1, $2) RETURNING id`
	var id int64
	err := db.QueryRow(sqlStatement, title, body).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
func main() {
	database.Initialize()
	db = database.DB

	route.Initialize()
	router = route.Router

	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")

	router.Use(forceHTMLMiddleware)

	http.ListenAndServe(":3000", removeTrilingSlash(router))
}
