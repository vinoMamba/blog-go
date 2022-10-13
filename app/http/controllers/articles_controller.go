package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/vinoMamba/goblog/pkg/logger"
)

type Article struct {
	ID    int64
	Title string
	Body  string
}

type ArticlesController struct {
}

func (*ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	// 从数据库中读取数据
	sqlStatement := `SELECT * FROM articles WHERE id=$1`
	//TODO 这里需要将操作数据库的代码抽离出来
	row := db.QueryRow(sqlStatement, id)
	var article Article
	err := row.Scan(&article.ID, &article.Title, &article.Body)
	switch err {
	case sql.ErrNoRows:
		fmt.Fprint(w, "没有这篇文章")
	case nil:
		tmpl, err := template.ParseFiles("./resources/views/articles/show.html")
		if err != nil {
			logger.LogError(err)
		}
		tmpl.Execute(w, article)
	default:
		logger.LogError(err)
	}
}
