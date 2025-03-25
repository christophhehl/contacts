package main

import (
	"database/sql"
	"embed"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed build/*
var staticFS embed.FS

var db *sql.DB

func main() {
	r := gin.Default()

	var err error
	db, err = sql.Open("sqlite3", "contacts.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	c := cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Accept", "Content-Type", "X-Requested-With", "Authorization"},
		AllowCredentials: true,
	})
	r.Use(c)

	r.Use(static.Serve("/", static.LocalFile("build", false)))

	r.GET("/contacts", getContacts)
	r.GET("/csvExport", getCsvExport)

	r.POST("/newContact", createNewContact)

	r.PUT("/updateContact", updateContact)

	r.DELETE("/deleteContact", deleteContact)

	err = r.Run(":8080")
	if err != nil {
		return
	} // Listen and serve on 0.0.0.0:8080
}

type contact struct {
	ID          int
	FirstName   string
	LastName    string
	Street      string
	HouseNumber string
	ZipCode     string
	City        string
	Partner     string
	Children    string
	Email       string
}

func getContacts(c *gin.Context) {
	var contacts []contact

	res, err := db.Query("SELECT id, FirstName, LastName, Street, HouseNumber, ZipCode, City, Partner, Children, Email FROM Contacts")
	if err != nil {
		log.Println("Error occurred fetching data: contacts; ", err)
	}

	defer res.Close()
	for res.Next() {
		var id int
		var FirstName string
		var LastName string
		var Street string
		var HouseNumber string
		var ZipCode string
		var City string
		var Partner string
		var Children string
		var Email string
		err := res.Scan(&id, &FirstName, &LastName, &Street, &HouseNumber, &ZipCode, &City, &Partner, &Children, &Email)
		if err != nil {
			log.Println("Error occurred reading data: contacts; ", err)
		}
		contacts = append(contacts, contact{
			ID:          id,
			FirstName:   FirstName,
			LastName:    LastName,
			Street:      Street,
			HouseNumber: HouseNumber,
			ZipCode:     ZipCode,
			City:        City,
			Partner:     Partner,
			Children:    Children,
			Email:       Email,
		})
	}

	c.JSON(http.StatusOK, contacts)
}

func getCsvExport(c *gin.Context) {

}

func createNewContact(c *gin.Context) {
	var con contact

	err := c.Bind(&con)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	_, err = db.Exec("INSERT INTO Contacts (FirstName, LastName, Street, HouseNumber, ZipCode, City, Partner, Children, Email) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", con.FirstName, con.LastName, con.Street, con.HouseNumber, con.ZipCode, con.City, con.Partner, con.Children, con.Email)
	if err != nil {
		log.Println("Couldn't insert contact: ", con.FirstName, "; err", err)
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func updateContact(c *gin.Context) {
	var con contact

	err := c.Bind(&con)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	_, err = db.Exec("UPDATE Contacts SET FirstName = ?, LastName = ?, Street = ?, HouseNumber = ?, ZipCode = ?, City = ?, Partner = ?, Children = ?, Email = ? WHERE id = ?", con.FirstName, con.LastName, con.Street, con.HouseNumber, con.ZipCode, con.City, con.Partner, con.Children, con.Email, con.ID)
	if err != nil {
		log.Println("Couldn't update contact; ", err)
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func deleteContact(c *gin.Context) {
	var con contact

	err := c.Bind(&con)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	_, err = db.Exec("DELETE FROM Contacts WHERE id = ?", con.ID)
	if err != nil {
		log.Println("Couldn't delete contact; ", err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
