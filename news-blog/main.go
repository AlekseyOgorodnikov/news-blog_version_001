package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Article struct {
	Id                     uint16
	Title, Anons, FullText string
}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("static/html/index.html", "static/html/header.html", "static/html/footer.html")
	if err != nil {
		fmt.Fprintf(w, "Fatal: %s\n", err.Error())
	}

	// Connect bd mysql
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		log.Fatal(err)
	}
	// Closed connect to bd mysql
	defer db.Close()

	res, err := db.Query("SELECT * FROM `ariclea`")
	if err != nil {
		panic(err)
	}

	var post Article
	var posts = []Article{}
	for res.Next() {
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}
		posts = append(posts, post)
		/* output data in terminal
		fmt.Printf("Id: %d\n Title: %s\n Anons: %s\n Text:%s\n.\n", post.Id, post.Title, post.Anons, post.FullText)
		*/
	}

	t.ExecuteTemplate(w, "index", posts)
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("static/html/create.html", "static/html/header.html", "static/html/footer.html")
	if err != nil {
		fmt.Fprintf(w, "Fatal: %s\n", err.Error())
	}

	t.ExecuteTemplate(w, "create", nil)
}

func save_article(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")

	if title == "" || anons == "" || full_text == "" {
		fmt.Fprintf(w, "Не все поля заполнены!")
	} else {
		// Connect bd mysql
		db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
		if err != nil {
			log.Fatal(err)
		}
		// Closed connect to bd mysql
		defer db.Close()

		// Create data in db
		insert, err := db.Prepare("INSERT INTO `ariclea` (`title`,`anons`,`full_text`) VALUES(?,?,?)")
		if err != nil {
			log.Fatal(err)
		}
		insert.Exec(title, anons, full_text)
		log.Println("INSERT: title: " + title + " | anons: " + anons + " | full_text: " + full_text)
		defer insert.Close()

		// Redirect page
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}

func showPost(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("static/html/show.html", "static/html/header.html", "static/html/footer.html")
	if err != nil {
		fmt.Fprintf(w, "Fatal: %s\n", err.Error())
	}

	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	/* // testing output id query in html template
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "ID: %v\n", vars["id"]) */

	// Connect bd mysql
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		log.Fatal(err)
	}
	// Closed connect to bd mysql
	defer db.Close()

	res, err := db.Query("SELECT * FROM `ariclea` WHERE id = ?", vars["id"])
	if err != nil {
		panic(err)
	}

	var showPosts = Article{}
	var post Article
	for res.Next() {
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}
		showPosts = post
	}

	t.ExecuteTemplate(w, "show", showPosts)
}

func handleFunc() {
	server := http.Server{
		Addr: "127.0.0.1:4040",
	}

	// Create router with gorilla/mux
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/create", create).Methods("GET")
	rtr.HandleFunc("/post/{id:[0-9]+}", showPost).Methods("GET")
	rtr.HandleFunc("/save_article", save_article).Methods("POST")

	http.Handle("/", rtr)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	log.Println("Сервер запущен на http://localhost:4040/ ...")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	handleFunc()
}
