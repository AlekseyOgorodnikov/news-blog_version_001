package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Name                  string
	Age                   uint16
	Money                 int16
	Avg_grades, Happiness float64
	Hobbies               []string
}

func (u User) getAllInfo() string {
	return fmt.Sprintf("User name is: %s. He is %d and he has "+
		"money equal %d.", u.Name, u.Age, u.Money)
}

func (u *User) setNewName(newName string) {
	u.Name = newName
}

func homePage(w http.ResponseWriter, r *http.Request) {
	bob := User{"Алексей", 31, 100, 4.2, 0.8, []string{"Футбол", "Баскетбол", "Кино"}}
	// bob.setNewName("Bob")
	// fmt.Fprintf(w, bob.getAllInfo())

	templateHTML, err := template.ParseFiles("static/html/homePage.html")
	if err != nil {
		fmt.Println("Error in build templates file!")
		return
	}

	templateHTML.Execute(w, bob)
}

func handleReq() {
	server := http.Server{
		Addr: "127.0.0.1:4040",
	}

	http.HandleFunc("/", homePage)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	log.Println("Сервер запущен на http://localhost:4040/ ...")

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	handleReq()
}
