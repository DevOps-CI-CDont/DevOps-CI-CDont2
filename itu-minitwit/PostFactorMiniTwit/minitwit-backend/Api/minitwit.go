package main

import (
	"crypto/sha256"
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	router := gin.Default()

	// endpoints
	router.GET("/", getTimeline)
	router.GET("/public", getPublicTimeline)
	router.GET("/user/:username", getUsername)
	router.POST("/user/:username/follow", followUser)
	router.POST("/user/:username/unfollow", unfollowUser)
	router.POST("/add_message", postMessage)
	router.POST("/login", login)
	router.POST("/register", register)
	router.GET("/logout", logout)

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

func errorCheck(err error) {
	if err != nil {
		log.Println(err)
	}
}

// endpoints
func getTimeline(c *gin.Context) {
	connect_db()
	//query database
	// check cookie for session,
	userID := getUserIdIfLoggedIn(c)

	rows, err := DB.Query(`select message.*, user.* from message, user
		where message.flagged = 0 and message.author_id = user.user_id and (
		user.user_id = ? or
		user.user_id in (select whom_id from follower
		where who_id = ?))
		order by message.pub_date desc limit ?`, userID, userID, PER_PAGE)
	errorCheck(err)
	defer rows.Close()
	messages := make([]Message, 0)
	for rows.Next() {
		msg := Message{}
		user := User{}
		err = rows.Scan(&msg.Message_id, &msg.Author_id, &msg.Text, &msg.Pub_date, &msg.Flagged, &user.User_id, &user.Username, &user.Email, &user.Pw_hash)
		errorCheck(err)
		log.Println(msg)
		msg.Author = user

		messages = append(messages, msg)
	}
	c.JSON(200, gin.H{"tweets": messages})
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
	name := c.Param("username")
	c.String(http.StatusOK, "Hello %s", name)
}

func followUser(c *gin.Context) {
	connect_db()

	userid := getUserIdIfLoggedIn(c)
	whom_name := c.Param("username")
	whom_id := getUserIdByName(whom_name)
	if doesUsersFollow(userid, whom_id) {
		c.JSON(200, gin.H{"message": "user already followed"})
		return
	}

	if whom_id == "-1" {
		c.JSON(200, gin.H{"message": "user does not exist"})
		return
	}
	stmt, err := DB.Prepare(`insert into follower (who_id, whom_id) values (?, ?)`)
	errorCheck(err)
	defer stmt.Close()

	_, err = stmt.Exec(userid, whom_id)
	errorCheck(err)

	c.JSON(200, gin.H{"message": "followed user"})
}

func doesUsersFollow(who_id string, whom_id string) bool {
	connect_db()

	row := DB.QueryRow(`select * from follower where who_id = ? and whom_id = ?`, who_id, whom_id)

	follower := follower{}
	err := row.Scan(&follower.Who_id, &follower.Whom_id)
	return err == nil

}

func unfollowUser(c *gin.Context) {
	connect_db()

	userid := getUserIdIfLoggedIn(c)
	whom_name := c.Param("username")
	whom_id := getUserIdByName(whom_name)
	if !doesUsersFollow(userid, whom_id) {
		c.JSON(200, gin.H{"message": "user dosent follow the target"})
		return
	}
	log.Println(whom_id)
	if whom_id == "-1" {
		c.JSON(200, gin.H{"message": "user does not exist"})
		return
	}

	stmt, err := DB.Prepare(`Delete FROM follower WHERE who_id = ? and whom_id = ?`)
	errorCheck(err)
	defer stmt.Close()

	_, err = stmt.Exec(userid, whom_id)
	errorCheck(err)

	c.JSON(200, gin.H{"message": "unfollowed user"})

}

func postMessage(c *gin.Context) {
	connect_db()

	userid := getUserIdIfLoggedIn(c)

	text := c.PostForm("text")
	authorid := userid
	pub_date := time.Now().Unix()
	flagged := 0

	stmt, err := DB.Prepare(`insert into message (author_id, text, pub_date, flagged) values (?, ?, ?, ?)`)
	errorCheck(err)
	defer stmt.Close()

	_, err = stmt.Exec(authorid, text, pub_date, flagged)
	errorCheck(err)

	c.JSON(200, gin.H{"message": "message posted"})

}

func login(c *gin.Context) {
	connect_db()

	username := c.PostForm("username")
	password := c.PostForm("password")
	//convert password to byte[]

	passwordHash := sha256.Sum256([]byte(password)) //hash password1
	passwordHashString := string(passwordHash[:])
	//check if username and password are empty
	if username == "" || password == "" {
		c.JSON(400, gin.H{"error": "username or password is empty"})
		return
	}

	userid := DB.QueryRow(`select user_id from user where user.Username = ? and user.pw_hash = ?`, username, passwordHashString)

	var userIdAsInt int
	err := userid.Scan(&userIdAsInt)
	if err != nil {
		c.SetCookie("user_id", "", -1, "/", "localhost", false, false)
		return
	}
	// succes: set cookie
	c.SetCookie("user_id", strconv.Itoa(userIdAsInt), 3600, "/", "localhost", false, false)
}

func register(c *gin.Context) {
	connect_db()

	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")
	password2 := c.PostForm("password2")

	//check if username and password are empty
	if username == "" || password == "" || password2 == "" {
		c.JSON(400, gin.H{"error": "username or password is empty"})
		return
	} else if email == "" || !strings.Contains(email, "@") {
		c.JSON(400, gin.H{"error": "email is empty or invalid"})
		return
	} else if password != password2 {
		c.JSON(400, gin.H{"error": "passwords don't match"})
		return
	} else if getUserByName(username) != nil {
		c.JSON(400, gin.H{"error": "username already exists"})
		return
	}

	passwordHash := sha256.Sum256([]byte(password))
	//convert back to string
	passwordHashString := string(passwordHash[:])
	log.Println(passwordHashString)

	stmt, err := DB.Prepare(`insert into user (username, email, pw_hash) values (?, ?, ?)`)
	errorCheck(err)
	defer stmt.Close()
	_, err = stmt.Exec(username, email, passwordHashString)
	errorCheck(err)

	c.JSON(200, gin.H{"message": "user registered"})
}

func getUserByName(userName string) *sql.Row {
	connect_db()
	row := DB.QueryRow(`select * from user where user.username = ?`, userName)
	user := User{}
	err := row.Scan(&user.User_id, &user.Username, &user.Email, &user.Pw_hash)
	if err != nil {
		return nil
	}
	return row
	// if user exists, return user, else return nil

}

func getUserIdIfLoggedIn(c *gin.Context) string {
	userid, err := c.Cookie("user_id")
	errorCheck(err)
	if userid == "" || userid == "-1" {
		c.JSON(401, gin.H{"error": "not logged in"})
		return "-1"
	}
	return userid

}

func getUserIdByName(username string) string {
	connect_db()
	row := DB.QueryRow(`select * from user where user.username = ?`, username)

	user := User{}
	err := row.Scan(&user.User_id, &user.Username, &user.Email, &user.Pw_hash)
	if err != nil {
		return "-1"
	}

	return strconv.Itoa(user.User_id)
}

func logout(c *gin.Context) {
	c.SetCookie("user_id", "", -1, "/", "localhost", false, false)
}
