package simulator

import (
	"encoding/json"
	"io"
	"io/ioutil"
	main "minitwit-backend/init/Api"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var LATEST = 0

func update_latest(c *gin.Context) {
	try_latest := c.Param("latest")
	int_val, err := strconv.Atoi(try_latest)
	if err != nil{
		fmt.Println(err)
	}
	if try_latest != "" && err == nil {
		LATEST = int_val
	}
}

func not_req_from_simulator(c *gin.Context) bool {
	from_simulator := c.Request.Header.Get("Authorization")

	return from_simulator != "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh"
}

func SetUpRouter() *gin.Engine {
	r := gin.Default()
	return r
}

func Start() {
	Router := SetUpRouter()

	// router config
	Router.Use(cors.Default()) // cors.Default() should allow all origins
	// it's important to set this before any routes are registered so that the middleware is applied to all routes
	// ALL MY HOMIES HATE CORS

	// endpoints
	Router.GET("/latest", getLatest)
	Router.POST("/register", register)
	Router.GET("/msgs", getMsgs)
	Router.GET("/msgs/:username", msgsPerUser)
	Router.POST("/msgs/:username", msgsPerUser)
	Router.GET("/fllws/:username", follow)
	Router.POST("/fllws/:username", follow)

	Router.Run(":8081")
}

func getLatest(c *gin.Context) {
	c.JSON(200, gin.H{
		"latest": LATEST,
	})
}

func register(c *gin.Context) {
	update_latest(c)

	bytes, _ := io.ReadAll(c.Request.Body)
	body := make(map[string]string)
	json.Unmarshal(bytes, &body)

	// retrieve data from request
	username := body["username"]
	email := body["email"]
	password := body["pwd"]
	password2 := body["pwd"]

	// create a new HTTP client
	client := &http.Client{}

	// create a new request with formData
	req, err := http.NewRequest("POST", "http://138.68.93.147:8080/register", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// add formData to the request body
	data := url.Values{}
	data.Set("username", username)
	data.Set("email", email)
	data.Set("password", password)
	data.Set("password2", password2)
	req.Body = ioutil.NopCloser(strings.NewReader(data.Encode()))

	// set the content-type header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// send the request
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer resp.Body.Close()

	// read the response body
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// convert the response body to a string
	bodyString := string(bodyBytes)
	c.JSON(200, gin.H{
		"status": 200,
		"msg":    bodyString,
	})
}

func getMsgs(c *gin.Context) {
	update_latest(c)

	if not_req_from_simulator(c) {
		c.JSON(403, gin.H{
			"status":    403,
			"error_msg": "You are not authorized to use this resource!",
		})
		return
	}

	// default
	no_msgs := 100
	if c.Param("no") != "" {
		no, _ := strconv.Atoi(c.Param("no"))
		no_msgs = no
	}

	query := `SELECT message.*, user.* FROM message, user WHERE message.flagged = 0
			 AND message.author_id = user.user_id ORDER BY message.pub_date DESC LIMIT ?`

	messages, _ := main.DB.Query(query, no_msgs)

	var filtered_msgs []map[string]string

	for messages.Next() {
		entry := make(map[string]string)
		var msg main.Message
		var user main.User
		messages.Scan(&msg.Message_id, &msg.Author_id, &msg.Text, &msg.Pub_date, &msg.Flagged, &user.User_id, &user.Username, &user.Email, &user.Pw_hash)
		msg.Author = user

		entry["content"] = msg.Text
		entry["pub_date"] = strconv.Itoa(msg.Pub_date)
		entry["user"] = msg.Author.Username
		filtered_msgs = append(filtered_msgs, entry)
	}

	c.JSON(200, filtered_msgs)
}

func msgsPerUser(c *gin.Context) {
	update_latest(c)

	if not_req_from_simulator(c) {
		c.JSON(403, gin.H{
			"status":    403,
			"error_msg": "You are not authorized to use this resource!",
		})
		return
	}

	// default
	no_msgs := 100
	if c.Param("no") != "" {
		no, _ := strconv.Atoi(c.Param("no"))
		no_msgs = no
	}

	if c.Request.Method == "GET" {
		user_id := main.GetUserIdByName(c.Param("username"))

		if user_id == "" {
			c.AbortWithStatus(404)
			return
		}

		query := `SELECT message.*, user.* FROM message, user 
					WHERE message.flagged = 0 AND
					user.user_id = message.author_id AND user.user_id = ?
					ORDER BY message.pub_date DESC LIMIT ?`

		messages, _ := main.DB.Query(query, user_id, no_msgs)

		var filtered_msgs []map[string]string

		for messages.Next() {
			entry := make(map[string]string)
			var msg main.Message
			var user main.User
			messages.Scan(&msg.Message_id, &msg.Author_id, &msg.Text, &msg.Pub_date, &msg.Flagged, &user.User_id, &user.Username, &user.Email, &user.Pw_hash)
			msg.Author = user

			entry["content"] = msg.Text
			entry["pub_date"] = strconv.Itoa(msg.Pub_date)
			entry["user"] = msg.Author.Username
			filtered_msgs = append(filtered_msgs, entry)
		}

		c.JSON(200, filtered_msgs)
	} else if c.Request.Method == "POST" {
		bytes, _ := io.ReadAll(c.Request.Body)
		body := make(map[string]string)
		json.Unmarshal(bytes, &body)

		query := `INSERT INTO message (author_id, text, pub_date, flagged)
					VALUES (?, ?, ?, 0)`

		main.DB.Exec(query, main.GetUserIdByName(c.Param("username")), body["content"], time.Now().Unix())

		c.JSON(204, gin.H{})
	}

}

func follow(c *gin.Context) {
	update_latest(c)

	if not_req_from_simulator(c) {
		c.JSON(403, gin.H{
			"status":    403,
			"error_msg": "You are not authorized to use this resource!",
		})
		return
	}

	user_id := main.GetUserIdByName(c.Param("username"))

	if user_id == "" {
		c.AbortWithStatus(404)
		return
	}

	bytes, _ := io.ReadAll(c.Request.Body)
	body := make(map[string]string)
	json.Unmarshal(bytes, &body)

	if c.Request.Method == "POST" && body["follow"] != "" {
		follows_username := body["follow"]
		follows_user_id := main.GetUserIdByName(follows_username)

		if follows_user_id == "-1" {
			c.AbortWithStatus(404)
			return
		}

		query := "INSERT INTO follower (who_id, whom_id) VALUES (?, ?)"
		main.DB.Exec(query, user_id, follows_user_id)

		c.JSON(204, gin.H{})
	} else if c.Request.Method == "POST" && body["unfollow"] != "" {
		unfollows_username := body["unfollow"]
		unfollows_user_id := main.GetUserIdByName(unfollows_username)

		if unfollows_user_id == "-1" {
			c.AbortWithStatus(404)
			return
		}

		query := "DELETE FROM follower WHERE who_id=? and WHOM_ID=?"
		main.DB.Exec(query, user_id, unfollows_user_id)

		c.JSON(204, gin.H{})
	} else if c.Request.Method == "GET" {
		// default
		no_followers := 100
		if c.Param("no") != "" {
			no, _ := strconv.Atoi(c.Param("no"))
			no_followers = no
		}

		query := `SELECT user.username FROM user
					INNER JOIN follower ON follower.whom_id=user.user_id
					WHERE follower.who_id=?
					LIMIT ?`

		followers, _ := main.DB.Query(query, user_id, no_followers)
		followers_names := []string{}

		for followers.Next() {
			var username string
			followers.Scan(&username)
			followers_names = append(followers_names, username)
		}

		c.JSON(200, followers_names)
	}

}
