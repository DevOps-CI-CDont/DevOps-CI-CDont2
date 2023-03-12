package Api

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var Router *gin.Engine

func SetUpRouter() *gin.Engine {
	r := gin.Default()
	return r
}

func Start() {
	Router = SetUpRouter()

	Connect_db()

	// router config
	Router.Use(cors.Default()) // cors.Default() should allow all origins
	// it's important to set this before any routes are registered so that the middleware is applied to all routes
	// ALL MY HOMIES HATE CORS :D

	// endpoints
	Router.GET("/mytimeline", getTimeline)
	Router.GET("/public", getPublicTimeline)
	Router.GET("/user/:username", getUsersTweets)
	Router.POST("/user/:username/follow", followUser)
	Router.POST("/user/:username/unfollow", unfollowUser)
	Router.POST("/add_message", postMessage)
	Router.POST("/login", login)
	Router.POST("/register", register)
	Router.GET("/logout", logout)
	Router.GET("/AmIFollowing/:username", amIFollowing)
	Router.GET("/allUsers", getAllUsers)
	Router.GET("AllIAmFollowing", getAllFollowing)
	Router.Run(":8080")
}

// Capitalized names are public, lowercase are privat
type User struct {
	gorm.Model
	Username string `json:"username"`
	Email    string `json:"email"`
	Pw_hash  string `json:"pw_hash"`
}

type follower struct {
	gorm.Model
	Who_id  int `json:who_id"`
	Whom_id int `json:whom_id"`
}

type Message struct {
	gorm.Model
	Author_id   int    `json:"author_id"`
	Text        string `json:"text"`
	Pub_date    int    `json:"pub_date"`
	Flagged     int    `json:"flagged"`
	Author_name string `json:"author_name"`
}

var DB *gorm.DB // global DB variable
var PER_PAGE = 30
var DEBUG = true
var SECRET_KEY = "development key"

func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	hexString := hex.EncodeToString(hash[:])
	return hexString
}

func Connect_db() error {
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s client_encoding=%s",
		"cicdont-do-user-13570987-0.b.db.ondigitalocean.com",
		25060,
		"doadmin",
		"AVNS_FeRFl5bSz6UNMVF6Llx",
		"minitwit",
		"require",
		"Europe/Berlin",
		"UTF8")
	db, err := gorm.Open(postgres.Open(dbinfo), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	//Migrate schema
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Message{})
	db.AutoMigrate(&follower{})

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

func amIFollowing(c *gin.Context) {
	username := c.Param("username")
	userID := getUserIdIfLoggedIn(c)

	var follower follower
	var user User
	err := DB.Table("followers").
		Where("who_id = ? AND whom_id = ?", userID, DB.Table("users").Select("user_id").Where("username = ?", username).Find(&user)).First(&follower).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(200, false)
			return
		}
		c.AbortWithStatusJSON(500, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(200, true)
}

func getTimeline(c *gin.Context) {
	userID := getUserIdIfLoggedIn(c)

	var messages []Message
	result := DB.Table("messages").
		Select("messages.*, users.*").
		Joins("JOIN users ON messages.author_id = users.user_id").
		Where("messages.flagged = ? AND (users.user_id = ? OR users.user_id IN (?))",
			0, userID, DB.Table("followers").Select("whom_id").Where("who_id = ?", userID)).
		Order("messages.pub_date DESC").
		Limit(PER_PAGE).
		Scan(&messages)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving messages"})
		return
	}

	c.JSON(200, gin.H{"tweets": messages})
}

func getPublicTimeline(c *gin.Context) {
	num_msgs := c.Request.URL.Query().Get("num_msgs")
	int_num_msgs, err := strconv.Atoi(num_msgs)
	if num_msgs == "" || err != nil {
		int_num_msgs = 30
	}

	fmt.Println("int_num_msgs", int_num_msgs)

	var messages []Message
	err = DB.
		Table("messages").
		Select("messages.*, users.*").
		Joins("JOIN users ON messages.author_id = users.user_id").
		Where("messages.flagged = ?", 0).
		Order("messages.pub_date desc").
		Limit(int_num_msgs).
		Find(&messages).Error

	if err != nil {
		c.JSON(500, gin.H{"error": "failed to retrieve messages"})
		return
	}

	// if no messages, return 401
	if len(messages) == 0 {
		c.JSON(401, gin.H{"message": "no messages"})
		return
	}

	fmt.Println("messages", messages)

	c.JSON(200, gin.H{"tweets": messages})
}

func getUsersTweets(c *gin.Context) {
	name := c.Param("username")
	user := User{}
	if err := DB.Where("username = ?", name).First(&user).Error; err != nil {
		c.JSON(200, gin.H{"message": "user does not exist"})
		return
	}

	num_msgs := c.Request.URL.Query().Get("num_msgs")
	int_num_msgs, err := strconv.Atoi(num_msgs)
	if num_msgs == "" || err != nil {
		int_num_msgs = 30
	}

	var messages []Message
	if err := DB.Where("author_id = ?", user.ID).Order("pub_date desc").Limit(int_num_msgs).Preload("Author").Find(&messages).Error; err != nil {
		errorCheck(err)
	}

	c.JSON(200, gin.H{"tweets": messages})
}

func followUser(c *gin.Context) {
	userID := getUserIdIfLoggedIn(c)
	whomName := c.Param("username")
	whomID := GetUserIdByName(whomName)

	if doesUsersFollow(userID, whomID) {
		c.JSON(200, gin.H{"message": "user already followed"})
		return
	}

	if whomID == "-1" {
		c.JSON(200, gin.H{"message": "user does not exist"})
		return
	}

	//convert userid and whomid to int
	userIDInt, err := strconv.Atoi(userID)
	errorCheck(err)
	whomIDInt, err := strconv.Atoi(whomID)
	errorCheck(err)

	follower := follower{
		Who_id:  userIDInt,
		Whom_id: whomIDInt,
	}

	result := DB.Create(&follower)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "failed to follow user"})
		return
	}

	c.JSON(200, gin.H{"message": "followed user"})
}

func doesUsersFollow(whoID string, whomID string) bool {
	var follower follower
	if err := DB.Where("who_id = ? and whom_id = ?", whoID, whomID).First(&follower).Error; err != nil {
		return false
	}
	return true
}

func unfollowUser(c *gin.Context) {
	userID := getUserIdIfLoggedIn(c)
	whomName := c.Param("username")
	whomID := GetUserIdByName(whomName)
	if !doesUsersFollow(userID, whomID) {
		c.JSON(200, gin.H{"message": "user doesn't follow the target"})
		return
	}
	log.Println(whomID)
	if whomID == "-1" {
		c.JSON(200, gin.H{"message": "user you are trying to follow does not exist"})
		return
	}

	result := DB.Where("who_id = ? AND whom_id = ?", userID, whomID).Delete(&follower{})
	if result.Error != nil {
		errorCheck(result.Error)
	}

	c.JSON(200, gin.H{"message": "unfollowed user"})
}

func postMessage(c *gin.Context) {
	userID := getUserIdIfLoggedIn(c)
	if userID == "-1" {
		c.JSON(401, gin.H{"message": "user not logged in"})
		return
	}

	text := c.PostForm("text")
	//Convert userid to int
	authorID, err := strconv.Atoi(userID)
	errorCheck(err)

	flagged := 0
	log.Println("text:" + text)
	//convert time.Now().Unix() to int
	pubDate := int(time.Now().Unix())

	message := Message{
		Author_id:   authorID,
		Author_name: GetUsernameByID(userID),
		Text:        text,
		Pub_date:    pubDate,
		Flagged:     flagged,
	}

	result := DB.Create(message)
	if result.Error != nil {
		c.JSON(500, gin.H{"message": "error creating message"})
		return
	}

	c.JSON(200, gin.H{"message": "message posted"})
}

func GetUsernameByID(id string) string {
	var user User
	if err := DB.Where("id = ?", id).First(&user).Error; err != nil {
		return "-1"
	}
	return user.Username
}

func login(c *gin.Context) {

	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" || password == "" {
		c.JSON(400, gin.H{"error": "username or password is empty"})
		return
	}

	var user User
	err := DB.Where("username = ? AND pw_hash = ?", username, HashPassword(password)).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.SetCookie("user_id", "", -1, "/", "localhost", false, false)
			c.JSON(401, gin.H{"error": "username or password is incorrect"})
			return
		} else {
			errorCheck(err)
		}
	}

	//convert user.ID as uint to string
	userID := strconv.Itoa(int(user.ID))

	c.SetCookie("user_id", userID, 3600, "/", "localhost", false, false)
	c.JSON(200, gin.H{"user_id": user.ID})
}

func register(c *gin.Context) {
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

	passwordHashString := HashPassword(password)
	log.Println(passwordHashString)

	user := User{
		Username: username,
		Email:    email,
		Pw_hash:  passwordHashString,
	}

	if err := DB.Create(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": "unable to create user"})
		return
	}

	c.JSON(200, gin.H{"message": "user registered"})
}

func getUserByName(userName string) *User {
	user := &User{}
	result := DB.Where("username = ?", userName).First(user)
	if result.Error != nil {
		return nil
	}
	return user
}

func getUserIdIfLoggedIn(c *gin.Context) string {
	userid, err := c.Cookie("user_id")
	log.Println("cookie user_id: " + userid)
	errorCheck(err)
	if userid == "" || userid == "-1" {
		c.JSON(401, gin.H{"error": "not logged in"})
		return "-1"
	}
	return userid

}

func GetUserIdByName(username string) string {
	var user User
	err := DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return "-1"
	}
	return strconv.Itoa(int(user.ID))
}

func getAllFollowing(c *gin.Context) {
	num_followers := c.Request.URL.Query().Get("num_followers")
	int_followers, err := strconv.Atoi(num_followers)
	if num_followers == "" || err != nil {
		int_followers = 30
	}
	userID := getUserIdIfLoggedIn(c)
	if userID == "-1" {
		c.JSON(401, gin.H{"error": "user not logged in"})
		return
	}

	following := []User{}
	err = DB.Table("users").
		Select("users.*").
		Joins("JOIN followers ON users.user_id = followers.whom_id").
		Where("followers.who_id = ?", userID).
		Limit(int_followers).
		Scan(&following).
		Error

	if err != nil {
		c.JSON(500, gin.H{"error": "unable to retrieve following"})
		return
	}
	c.JSON(200, gin.H{"following": following})

}

func logout(c *gin.Context) {
	c.SetCookie("user_id", "", -1, "/", "localhost", false, false)
}

func getAllUsers(c *gin.Context) {
	users := []User{}
	err := DB.Find(&users).Error
	if err != nil {
		c.JSON(500, gin.H{"error": "unable to retrieve users"})
		return
	}
	c.JSON(200, gin.H{"users": users})
}
