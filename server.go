package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Contact struct {
	Id 		int
	Name 	string
	Fone  	string
}

var db, err = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/db-contacts")

func getById(c *gin.Context) {
	var (
		contact Contact
		result gin.H
	)
	id := c.Param("id")
	row := db.QueryRow("select id, Name, Fone from contact where id = ?;", id)
	err = row.Scan(&contact.Id, &contact.Name, &contact.Fone)
	if err != nil {
		// If no results send null
		result = gin.H{
			"result": nil,
			"count":  0,
		}
	} else {
		result = gin.H{
			"result": contact,
			"count":  1,
		}
	}
	c.JSON(http.StatusOK, result)
}

func getAll(c *gin.Context){
	var (
			contact Contact
			contacts []Contact
		)
		rows, err := db.Query("select id, Name, Fone from contact;")
		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&contact.Id, &contact.Name, &contact.Fone)
			contacts = append(contacts, contact)
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		c.JSON(http.StatusOK, gin.H)
}

func add(c *gin.Context) {
		var buffer bytes.Buffer
		Name := c.PostForm("Name")
		Fone := c.PostForm("Fone")
		stmt, err := db.Prepare("insert into contact (Name, Fone) values(?,?);")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(Name, Fone)

		if err != nil {
			fmt.Print(err.Error())
		}

		// Fastest way to append strings
		buffer.WriteString(Name)
		buffer.WriteString(" ")
		buffer.WriteString(Fone)
		defer stmt.Close()
		name := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf(" %s successfully created", name),
		})
}

func update(c *gin.Context) {
		var buffer bytes.Buffer
		id := c.Query("id")
		Name := c.PostForm("Name")
		Fone := c.PostForm("Fone")
		stmt, err := db.Prepare("update contact set Name= ?, Fone= ? where id= ?;")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(Name, Fone, id)
		if err != nil {
			fmt.Print(err.Error())
		}

		// Fastest way to append strings
		buffer.WriteString(Name)
		buffer.WriteString(" ")
		buffer.WriteString(Fone)
		defer stmt.Close()
		name := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully updated to %s", name),
		})
}

func delete(c *gin.Context) {
		id := c.Query("id")
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
	// make sure connection is available
	err = db.Ping()
	if err != nil {
		fmt.Print(err.Error())
	}
	
	router := gin.Default()
	router.GET("/contact/:id", getById)
	router.GET("/contacts", getAll)
	router.POST("/contact", add)
	router.PUT("/contact", update)
	router.DELETE("/contact", delete)
	router.Run(":8000")
}
