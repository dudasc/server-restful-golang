package main

import (
	"time"
	"database/sql"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func delaySecond(n time.Duration) {
	time.Sleep(n * time.Second)
}

type Contact struct {
	Id 		int
	Name 	string
	Fone  	string
}

var db, err = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/db-contacts")

func getById(c *gin.Context) {
	var contact Contact

	id := c.Param("id")
	row := db.QueryRow("select id, Name, Fone from contact where id = ?;", id)
	
	err = row.Scan(&contact.Id, &contact.Name, &contact.Fone)
	if err != nil {
		c.JSON(http.StatusOK, nil)
	} else {
		c.JSON(http.StatusOK, contact)
	}	
}

func getAll(c *gin.Context){
	var (
		contact Contact
		contacts []Contact
	)

	rows, err := db.Query("select id, name, fone from contact;")
	if err != nil {
		fmt.Print(err.Error())
	}
	for rows.Next() {
		rows.Scan(&contact.Id, &contact.Name, &contact.Fone)
		contacts = append(contacts, contact)		
	}
	defer rows.Close()

	c.JSON(http.StatusOK, contacts)
}

func add(c *gin.Context) {
	Name := c.PostForm("name")
	Fone := c.PostForm("fone")

	stmt, err := db.Prepare("insert into contact (Name, Fone) values(?,?);")
	if err != nil {
		fmt.Print(err.Error())
	}

	_, err = stmt.Exec(Name, Fone)
	if err != nil {
		fmt.Print(err.Error())
	}

	defer stmt.Close()

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("successfully created"),
	})
}

func update(c *gin.Context) {
	Id := c.Param("id")
	Name := c.PostForm("name")
	Fone := c.PostForm("fone")

	stmt, err := db.Prepare("update contact set Name= ?, Fone= ? where Id= ?;")
	if err != nil {
		fmt.Print(err.Error())
	}

	_, err = stmt.Exec(Name, Fone, Id)
	if err != nil {
		fmt.Print(err.Error())
	}

	defer stmt.Close()

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully updated"),
	})
}

func delete(c *gin.Context) {
	id := c.Param("id")
	stmt, err := db.Prepare("delete from contact where id= ?;")
	if err != nil {
		fmt.Print(err.Error())
	}
	_, err = stmt.Exec(id)
	if err != nil {
		fmt.Print(err.Error())
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully deleted: %s", id),
	})
}

//Criação da tabela de contato
func createTable(c *gin.Context){
	stmt, err := db.Prepare("CREATE TABLE contact (id int NOT NULL AUTO_INCREMENT, Name varchar(40), Fone varchar(40), PRIMARY KEY (id));")
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Contact Table successfully migrated....")
	}
}


func main() {
	if err != nil {
		fmt.Print(err.Error())
	}
	defer db.Close()
	
	err = db.Ping()
	if err != nil {
		fmt.Print(err.Error())
	}
	
	router := gin.Default()
	router.GET("/api/contact/:id", getById)
	router.GET("/api/contacts", getAll)
	router.POST("/api/contact", add)
	router.PUT("/api/contact/:id", update)
	router.DELETE("/api/contact/:id", delete)
	router.Run(":8000")
}
