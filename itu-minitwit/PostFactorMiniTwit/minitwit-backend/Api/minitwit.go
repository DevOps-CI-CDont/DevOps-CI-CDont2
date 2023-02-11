package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	router := gin.Default()

	// endpoints
	router.GET("/", getTimeline)
	router.GET("/public", getPublicTimeline)
	router.GET("/user/:username", getUsername)
	router.POST("/follow", followUser)
	router.POST("/unfollow", unfollowUser)
	router.POST("/post", postMessage)
	router.POST("/login", login)
	router.POST("/register", register)
	router.POST("/logout", logout)

	router.Run("localhost:8080")
}

// Capitalized names are public, lowercase are privat
type User struct {
	User_id  int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Pw_hash  string `json:"pw_hash"`
}

type follower struct {
	Who_id  int `json:"who_id"`
	Whom_id int `json:"whom_id"`
}

type Message struct {
	Message_id int    `json:"message_id"`
	Author_id  int    `json:"author_id"`
	Text       string `json:"text"`
	Pub_date   int    `json:"pub_date"`
	Flagged    int    `json:"flagged"`
	Author     User   `json:"author"`
}

var dbPath = "./../../tmp/minitwit.db"
var DB *sql.DB // global DB variable
var PER_PAGE = 30
var DEBUG = true
var SECRET_KEY = "development key"

func connect_db() error {
	db, err := sql.Open("sqlite3", dbPath)
	errorCheck(err)

	DB = db
	return nil
}

// endpoints
func getTimeline(c *gin.Context) {
	log.Println("getTimeline called")
	//query database
}

func errorCheck(err error) {
	if err != nil {
		log.Println(err)
	}
}

func getPublicTimeline(c *gin.Context) {
	connect_db()
	log.Println("connect_db done")
	rows, err := DB.Query(`select message.*, user.* from message, user
	where message.flagged = 0 and message.author_id = user.user_id
	order by message.pub_date desc limit ?`, PER_PAGE)
	errorCheck(err)

	// make a empty slice of messages
	messages := make([]Message, 0)

	for rows.Next() {
		msg := Message{}
		user := User{}
		err = rows.Scan(&msg.Message_id, &msg.Author_id, &msg.Text, &msg.Pub_date, &msg.Flagged, &user.User_id, &user.Username, &user.Email, &user.Pw_hash)
		log.Println(msg)
		msg.Author = user

		messages = append(messages, msg)
	}
	log.Println("printing messages")
	log.Println(messages)

	errorCheck(err)

	c.JSON(200, gin.H{"tweets": messages})

	defer rows.Close()
}

func getUsername(c *gin.Context) {
	log.Println("getUsername called")
	// display user's tweet
}

func followUser(c *gin.Context) {
	log.Println("followUser called")
}

func unfollowUser(c *gin.Context) {
	log.Println("unfollowUser called")
}

func postMessage(c *gin.Context) {
	log.Println("postMessage called")
}

func login(c *gin.Context) {
	log.Println("login called")
}

func register(c *gin.Context) {
	log.Println("register called")
}

func logout(c *gin.Context) {
	log.Println("logout called")
}
