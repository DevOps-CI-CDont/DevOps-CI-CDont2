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

type User struct {
	user_id  int    `json:"user_id"`
	username string `json:"username"`
	email    string `json:"email"`
	pw_hash  string `json:"pw_hash"`
}

type follower struct {
	who_id  int `json:"who_id"`
	whom_id int `json:"whom_id"`
}

type Message struct {
	message_id int    `json:"message_id"`
	author_id  int    `json:"author_id"`
	text       string `json:"text"`
	pub_date   int    `json:"pub_date"`
	flagged    int    `json:"flagged"`
}

var dbPath = "./../../tmp/minitwit.db"
var DB *sql.DB // global DB variable

func connect_db() error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	DB = db
	return nil
}

// endpoints
func getTimeline(c *gin.Context) {
	log.Println("getTimeline called")
	//query database
}

func getPublicTimeline(c *gin.Context) {
	connect_db()
	log.Println("connect_db done")
	rows, err := DB.Query(`select * from message`, 30)
	//print length of rows
	if err != nil {
		log.Println("đßðßđð")
		log.Println(err)
		log.Fatal(err)
	}

	log.Println("printing rows")
	log.Println(rows)
	// make a dummy message
	messages := make([]Message, 0)

	// for rows.Next() {
	// 	msgtest := Message{}
	// 	messages = append(messages, msgtest)
	// }

	for rows.Next() {
		msg := Message{}
		err = rows.Scan(&msg.message_id, &msg.author_id, &msg.text, &msg.pub_date, &msg.flagged)
		log.Println(&msg.message_id)
		messages = append(messages, msg)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(200, gin.H{"data": messages})

	/* 	for rows.Next() {
	   		err := rows.Scan(&messages)
	   		if err != nil {
	   			log.Fatal(err)
	   		}
	   	}

	   	c.JSON(200, gin.H{"msg": "u called public timeline"}) */

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
