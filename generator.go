package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
)

type Company struct {
	ID          int
	MainBranch  string
	Branch      string
	Name        string
	Website     string
	Mail        string
	PhoneNumber string
}

type Config struct {
	Config_string string `json:"connection_string"`
	Driver        string `json:"driver"`
}

func read_config(path string, c *gin.Context) Config {

	var config Config
	jsonFile, err := os.Open(path)

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}
	log.Println("Successfully Opened " + path)

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &config)

	return config
}

func connect_to_db(config Config, c *gin.Context) *sql.DB {

	var db *sql.DB
	db, err := sql.Open(config.Driver, config.Config_string)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}
	pingErr := db.Ping()
	if pingErr != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	return db
}

func execute_query_companies(db *sql.DB, query string, c *gin.Context) {
	var company Company
	var query_result []Company

	if strings.ToLower(query) == "all" {
		query = "SELECT * FROM [dbo].[Companies]"
	}

	rows, err := db.Query(query)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&company.ID, &company.MainBranch, &company.Branch, &company.Name, &company.Website, &company.Mail, &company.PhoneNumber)

		query_result = append(query_result, company)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}
	}

	err = rows.Err()
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	c.IndentedJSON(http.StatusOK, query_result)
}

func execute_query_on_companies(c *gin.Context) {

	log.Println("Connection attempt to execute query on companies table.")
	query := c.Param("query")
	config := read_config("config.json", c)
	db := connect_to_db(config, c)

	log.Println("New Connection!")
	execute_query_companies(db, query, c)
	db.Close()

	log.Println("Query executed and Connection closed.")

}

func webscrape(c *gin.Context) {
	url := "https://www.wunderground.com/history/daily/EPGD/date/2021-11-20"
	resp, err := http.Get(url)
	if err != nil {
		c.String(http.StatusBadRequest, "Couldn't make request.")
	}
	defer resp.Body.Close()
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Couldn't read html body.")
	}
	c.String(http.StatusOK, string(html))
}

func main() {
	router := gin.Default()
	router.GET("/companies/:query", execute_query_on_companies)
	router.GET("/webscraper/", webscrape)
	router.Run("localhost:8080")
}
