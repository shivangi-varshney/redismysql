package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
)

var tpl *template.Template
var f string
var l string

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func foo(w http.ResponseWriter, req *http.Request) {

	f = req.FormValue("Party Name")
	l = req.FormValue("STATE")
	fmt.Println(f, l)
	db, err := sql.Open("mysql", "root:Shivangi18!@tcp(127.0.0.1:3306)/project")

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	fmt.Println("Successfully Connected")
	fmt.Println("Yes it is connected")

	pool := newPool()

	conn := pool.Get()

        defer conn.Close()

	insert, err1 := db.Prepare("INSERT INTO article (title,content) VALUES (?,?)")

	if err1 != nil {
		panic(err1.Error())
	}

	insert.Exec(f, l)
	defer insert.Close()


        selDB,err2:=db.Query("SELECT * FROM article WHERE id=15")
        if err2 != nil {
		panic(err2.Error())
	}
	
       usr := User{}
       res := []User{}
      for selDB.Next() {
        var id int
        var title, content string
        err3 := selDB.Scan(&id, &title, &content)
        if err3 != nil {
            panic(err3.Error())
        }
        usr.Id = id
        usr.Title = title
        usr.Content = content
        res = append(res, usr)
    }
        //defer res.Close()

        err = setStruct(conn,usr.Id,usr.Title,usr.Content)
	if err != nil {
		fmt.Println(err)
	}

	
	/*err2 := tpl.ExecuteTemplate(w, "index.gohtml", person{f, l, s})
	if err2 != nil {
		http.Error(w, err2.Error(), 500)
		log.Fatalln(err2)
	}*/

}
func login(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "nomination_form.gohtml", nil)

}
func newPool() *redis.Pool {
	return &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 80,
		// max number of connections
		MaxActive: 12000,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}  
type User struct {
	Id int `json:"id"`
	Title  string `json:"title"`
	Content  string `json:"content"`
}
func setStruct(c redis.Conn,i int,t string,co string) error {

	const objectPrefix string = "user:"

	usr := User{
                Id: i,
		Title: t ,
		Content: co ,
	}

	// serialize User object to JSON
	json, err := json.Marshal(usr)
	fmt.Println(reflect.TypeOf(json))
	if err != nil {
		return err
	}

	// SET object
	_, err = c.Do("SET", usr.Title, json)
	if err != nil {
		return err
	}

	return nil
}


func main() {

	
	http.HandleFunc("/", login)
	http.HandleFunc("/foo", foo)
	http.ListenAndServe(":8002", nil)

}
