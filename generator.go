package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/anaskhan96/soup"
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

func execute_query(db *sql.DB, query string, c *gin.Context) {
	var company Company
	var query_result []Company

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

func execute_query_by_api(c *gin.Context) {

	log.Println("Connection attempt to execute query on database.")
	query := c.Param("query")
	config := read_config("config.json", c)
	db := connect_to_db(config, c)

	log.Println("New Connection!")
	execute_query(db, query, c)
	db.Close()

	log.Println("Query executed and Connection closed.")

}

func select_table(db *sql.DB, table string, c *gin.Context) {
	var company Company
	var query_result []Company

	query := fmt.Sprintf("SELECT * FROM [dbo].[%s]", table)

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

func select_all_by_api(c *gin.Context) {

	log.Println("Connection attempt to select whole table.")
	table := c.Param("table")
	config := read_config("config.json", c)
	db := connect_to_db(config, c)

	log.Println("New Connection!")
	select_table(db, table, c)
	db.Close()

	log.Println("Query executed and Connection closed.")

}

func get_html_source(url string, c *gin.Context) string {
	url = "https://" + url
	resp, err := http.Get(url)
	if err != nil {
		c.String(http.StatusBadRequest, "Couldn't make request.")
	}
	defer resp.Body.Close()
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Couldn't read html body.")
	}
	return string(html)
}

func get_html(c *gin.Context) {
	url := string(c.Param("url"))
	html := get_html_source(url, c)
	c.String(http.StatusOK, html)
}

func get_html_tag(c *gin.Context) {

	var response []string
	url := string(c.Param("url"))
	//tag := string(c.Param("tag"))
	html := get_html_source(url, c)
	doc := soup.HTMLParse(html)

	occurences := doc.FindAll("button", "type", "button")
	for _, link := range occurences {
		content := link.Text()
		response = append(response, content)
		log.Println(content)

	}

	json, err := json.Marshal(response)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	c.IndentedJSON(http.StatusOK, json)
}

func get_html_id(c *gin.Context) {

	url := string(c.Param("url"))
	// id := string(c.Param("id"))
	html := get_html_source(url, c)
	// doc := soup.HTMLParse(html)

	c.String(http.StatusOK, html)
}

func get_html_class(c *gin.Context) {

	url := string(c.Param("url"))
	// class := string(c.Param("class"))
	html := get_html_source(url, c)
	// doc := soup.HTMLParse(html)

	c.String(http.StatusOK, html)
}

func main() {
	router := gin.Default()
	router.GET("/db/:query", execute_query_by_api)
	router.GET("/db/all/:table", select_all_by_api)
	router.GET("/webscraper/webpage/:url", get_html)
	router.GET("/webscraper/webpage/:url/tag/:tag", get_html_tag)
	router.GET("/webscraper/webpage/:url/id/:id")
	router.GET("/webscraper/webpage/:url/class/:class")
	router.Run("localhost:8080")
}
