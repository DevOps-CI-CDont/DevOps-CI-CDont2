package Api

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"minitwit-backend/init/config"
	"minitwit-backend/init/models"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	"gorm.io/gorm"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var Router *gin.Engine

type metrics struct {
	funcCounter *prometheus.CounterVec
	memoryUsage prometheus.Gauge
}

func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		funcCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "function_calls_total",
			Help: "Number of calls to each function",
		}, []string{"method", "endpoint", "code"}),
		memoryUsage: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "Memory_usage_percentage",
			Help: "Infrastructure monitoring",
		}),
	}
	reg.MustRegister(m.funcCounter)
	reg.MustRegister(m.memoryUsage)
	return m
}

func SetUpRouter() *gin.Engine {
	r := gin.Default()
	return r
}

func Start(mode string) {
	log.SetOutput(os.Stdout)

	Router = SetUpRouter()

	if mode == "test" {
		config.Connect_test_db()
	} else {
		config.Connect_prod_db()
	}

    config := cors.DefaultConfig()
    config.AllowAllOrigins = true
    config.AllowHeaders = []string{"Authorization", "content-type"}
    Router.Use(cors.New(config))

	// router config
	//Router.Use(cors.Default()) // cors.Default() should allow all origins
	// it's important to set this before any routes are registered so that the middleware is applied to all routes
	// ALL MY HOMIES HATE CORS :D

	// metrics
	reg := prometheus.NewRegistry()
	m := NewMetrics(reg)
	go infrastructureGauge(10, m)

	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})

	// endpoints
	Router.GET("/metrics", gin.WrapH(promHandler))
	Router.GET("/mytimeline", getTimeline, incrementCounter(m, "/mytimeline"))
	Router.GET("/public", getPublicTimeline, incrementCounter(m, "/public"))
	Router.GET("/user/:username", getUsersTweets, incrementCounter(m, "/user/:username"))
	Router.POST("/user/:username/follow", followUser, incrementCounter(m, "/user/:username/follow"))
	Router.POST("/user/:username/unfollow", unfollowUser, incrementCounter(m, "/user/:username/unfollow"))
	Router.POST("/add_message", postMessage, incrementCounter(m, "/add_message"))
	Router.POST("/login", login, incrementCounter(m, "/login"))
	Router.POST("/register", register, incrementCounter(m, "/register"))
	Router.GET("/logout", logout, incrementCounter(m, "/logout"))
	Router.GET("/AmIFollowing/:username", amIFollowing, incrementCounter(m, "/AmIFollowing/:username"))
	Router.GET("/allUsers", getAllUsers, incrementCounter(m, "/allUsers"))
	Router.GET("/AllIAmFollowing", getAllFollowing) // is this getting used? @TODO
	Router.GET("/getUserNameById", GetUsernameByIDEndpoint, incrementCounter(m, "/getUserNameById"))
	Router.POST("/flagTweet", flagTweet, incrementCounter(m, "/flagTweet"))

	Router.Run(":8080")
}

func incrementCounter(m *metrics, endpointName string) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		m.funcCounter.WithLabelValues(c.Request.Method, endpointName, strconv.Itoa(c.Writer.Status())).Inc()
	}
	return fn
}

func GetUsernameByIDEndpoint(c *gin.Context) {
	userID := c.Request.URL.Query().Get("id")
	var user models.User
	err := config.DB.Table("users").Where("id = ?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(400, false)
			return
		}
		fmt.Println("error", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(200, user.Username)
}

func infrastructureGauge(intervalInSeconds int, m *metrics) {
	for {
		memory, err := memory.Get()
		if err != nil {
			log.Printf("Error reading memory usage: %s", err)
			return
		}
		memoryUsageInPercent := float64(memory.Used * 100 / memory.Total)

		m.memoryUsage.Set(memoryUsageInPercent)
		time.Sleep(time.Duration(intervalInSeconds) * time.Second)
	}
}

var PER_PAGE = 30
var DEBUG = true

func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	hexString := hex.EncodeToString(hash[:])
	return hexString
}

func errorCheck(err error) {
	if err != nil {
		log.Println(err)
	}
}

// endpoints

func amIFollowing(c *gin.Context) {
	nameTryingToFollow := c.Param("username")
	userID := getUserIdIfLoggedIn(c)

	var follower models.Follower
	var whom models.User

	err := config.DB.Table("users").
		Where("username = ?", nameTryingToFollow).
		First(&whom).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(400, false)
			return
		}
		fmt.Println("error", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "internal server error"})
		return
	}

	err = config.DB.Table("followers").
		Where("who_id = ? AND whom_id = ?", userID, whom.ID).
		First(&follower).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(200, false)
			return
		}
		fmt.Println("error", err)
		c.AbortWithStatusJSON(500, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(200, true)
}

func getTimeline(c *gin.Context) {
	userID := getUserIdIfLoggedIn(c)

	fmt.Println(userID)

	var messages []models.Message
	result := config.DB.Table("messages").
		Select("messages.*, users.username as user_name").
		Joins("JOIN users ON messages.author_id = users.id").
		Where("messages.flagged = ? AND (users.id = ? OR users.id IN (?)) AND (messages.created_at is not null)",
			0, userID, config.DB.Table("followers").Select("whom_id").Where("who_id = ?", userID)).
		Order("messages.created_at DESC").
		Limit(PER_PAGE).
		Scan(&messages)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving messages"})
		return
	}

	if len(messages) == 0 {
		c.JSON(200, gin.H{"tweets": []models.Message{}})
		return
	}
	c.JSON(200, gin.H{"tweets": messages})
}

func getPublicTimeline(c *gin.Context) {
	// dummy log to test ELK functionality
	log.Info("getPublicTimeline called!")

	num_msgs := c.Request.URL.Query().Get("num_msgs")
	int_num_msgs, err := strconv.Atoi(num_msgs)
	if num_msgs == "" || err != nil {
		int_num_msgs = 30
	}

	fmt.Println("int_num_msgs", int_num_msgs)

	var messages []models.Message
	err = config.DB.
		Table("messages").
		Select("messages.*, users.username as user_name").
		Joins("JOIN users ON messages.author_id = users.id").
		Where("messages.flagged = ? AND messages.created_at IS NOT NULL", 0).
		Order("messages.created_at desc").
		Limit(int_num_msgs).
		Find(&messages).Error

	if err != nil {
		c.JSON(500, gin.H{"error": "failed to retrieve messages"})
		return
	}

	if len(messages) == 0 {
		c.JSON(200, gin.H{"tweets": []models.Message{}})
		return
	}

	log.Println("messages", messages)

	c.JSON(200, gin.H{"tweets": messages})
}

func getUsersTweets(c *gin.Context) {
	name := c.Param("username")
	user := models.User{}
	if err := config.DB.Where("username = ?", name).First(&user).Error; err != nil {
		c.JSON(200, gin.H{"message": "user does not exist"})
		return
	}

	num_msgs := c.Request.URL.Query().Get("num_msgs")
	int_num_msgs, err := strconv.Atoi(num_msgs)
	if num_msgs == "" || err != nil {
		int_num_msgs = 30
	}

	var messages []models.Message
	if err := config.DB.Where("author_id = ?", user.ID).Order("pub_date desc").Limit(int_num_msgs).Find(&messages).Error; err != nil {
		errorCheck(err)
	}

	c.JSON(200, gin.H{"tweets": messages})
}

func followUser(c *gin.Context) {
	userID := getUserIdIfLoggedIn(c)
	if userID == "-1" {
		c.JSON(401, gin.H{"message": "user not logged in"})
		return
	}
	whomName := c.Param("username")
	whomID := GetUserIdByName(whomName)

	if doesUsersFollow(userID, whomID) {
		c.JSON(200, gin.H{"message": "user already followed"})
		return
	}

	if whomID == "-1" {
		c.JSON(400, gin.H{"message": "user does not exist"})
		return
	}

	//convert userid and whomid to int
	userIDInt, err := strconv.Atoi(userID)
	errorCheck(err)
	whomIDInt, err := strconv.Atoi(whomID)
	errorCheck(err)

	Follower := models.Follower{
		Who_id:  userIDInt,
		Whom_id: whomIDInt,
	}

	if err := config.DB.Create(&Follower).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed to follow user"})
		return
	}

	c.JSON(200, gin.H{"message": "followed user"})
}

func doesUsersFollow(whoID string, whomID string) bool {
	var follower models.Follower
	if err := config.DB.Where("who_id = ? and whom_id = ?", whoID, whomID).First(&follower).Error; err != nil {
		return false
	}
	return true
}

func unfollowUser(c *gin.Context) {
	userID := getUserIdIfLoggedIn(c)
	if userID == "-1" {
		c.JSON(401, gin.H{"message": "user not logged in"})
		return
	}
	whomName := c.Param("username")
	whomID := GetUserIdByName(whomName)
	if !doesUsersFollow(userID, whomID) {
		c.JSON(200, gin.H{"message": "user doesn't follow the target"})
		return
	}
	log.Println("user id" + userID + " trying to unfollow " + whomID)
	if whomID == "-1" {
		c.JSON(400, gin.H{"message": "user you are trying to unfollow does not exist"})
		return
	}

	result := config.DB.Where("who_id = ? AND whom_id = ?", userID, whomID).Delete(&models.Follower{})
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
	log.Println("tweet attempting to be posted:" + text)
	//convert time.Now().Unix() to int
	pubDate := int(time.Now().Unix())

	message := models.Message{
		Author_id:   authorID,
		Author_name: GetUsernameByID(userID),
		Text:        text,
		Pub_date:    pubDate,
		Flagged:     flagged,
	}

	/* result := config.DB.Create(message)
	if result.Error != nil {
		c.JSON(500, gin.H{"message": "error creating message"})
		return
	} */

	if err := config.DB.Create(&message).Error; err != nil {
		c.JSON(400, gin.H{"error": "unable to create message"})
		return
	}

	c.JSON(200, gin.H{"message": "message posted"})
}

func GetUsernameByID(id string) string {
	var user models.User
	if err := config.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return "-1"
	}
	return user.Username
}

func login(c *gin.Context) {

	username := c.PostForm("username")
	password := c.PostForm("password")

	log.Println("username: " + username)
	log.Println("password: " + password)

	if username == "" || password == "" {
		c.JSON(400, gin.H{"error": "username or password is empty"})
		return
	}

	var user models.User
	err := config.DB.Where("username = ? AND pw_hash = ?", username, HashPassword(password)).First(&user).Error

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

	user := models.User{
		Username: username,
		Email:    email,
		Pw_hash:  passwordHashString,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": "unable to create user"})
		return
	}

	c.JSON(200, gin.H{"message": "user registered"})
}

func getUserByName(userName string) *models.User {
	user := &models.User{}
	result := config.DB.Where("username = ?", userName).First(user)
	log.Println("getUserByName result: ", result)
	log.Println("error: ", result.Error)
	if result.Error != nil {
		fmt.Println("didn't find user with username: " + userName + ": this is expected for new users")
		return nil
	}
	return user
}

func getUserIdIfLoggedIn(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")

	// if userid == "" || userid == "-1" {
	// 	return "-1"
	// }
	return authHeader

}

func GetUserIdByName(username string) string {
	var user models.User
	err := config.DB.Where("username = ?", username).First(&user).Error
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

	following := []models.User{}
	err = config.DB.Table("users").
		Select("users.*").
		Joins("JOIN followers ON users.id = followers.whom_id").
		Where("followers.who_id = ? AND followers.deleted_at IS NULL", userID).
		Limit(int_followers).
		Scan(&following).
		Error

	if err != nil {
		c.JSON(500, gin.H{"error": "unable to retrieve following"})
		return
	}
	c.JSON(200, following)

}

func logout(c *gin.Context) {
	c.SetCookie("user_id", "", -1, "/", "localhost", false, false)
}

func getAllUsers(c *gin.Context) {
	users := []models.User{}
	err := config.DB.Find(&users).Error
	if err != nil {
		c.JSON(500, gin.H{"error": "unable to retrieve users"})
		return
	}
	c.JSON(200, gin.H{"users": users})
}

func flagTweet(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	// check if .env file exists
	_, err := os.Stat("../.env")
	if os.IsNotExist(err) {
		fmt.Println("no .env file found")
	} else { // load .env file if it exists
		err := godotenv.Load("../.env")
		if err != nil {
			fmt.Println("err = ", err)
			log.Fatal("Error loading .env file")
		}
	}
	authShouldBe := os.Getenv("FLAG_AUTH")
	if auth == "" {
		c.JSON(401, gin.H{"error": "no Authorization header"})
		return
	} else if auth != authShouldBe {
		c.JSON(401, gin.H{"error": "invalid Authorization header"})
		return
	} else {
		log.Println("flagging tweet")
		messageID := c.Request.URL.Query().Get("message_id")
		flagValue := c.Request.URL.Query().Get("flag_value")
		// conv flagval to int
		int_flagval, err := strconv.Atoi(flagValue)
		if err != nil {
			c.JSON(400, gin.H{"error": "flag value is not an integer"})
			return
		}
		message := &models.Message{}
		err = config.DB.Where("id = ?", messageID).First(message).Error
		if err != nil {
			c.JSON(500, gin.H{"error": "unable to retrieve message"})
			return
		}
		message.Flagged = int_flagval
		err = config.DB.Save(message).Error
		if err != nil {
			c.JSON(500, gin.H{"error": "unable to flag message"})
			return
		}
		if int_flagval == 1 {
			log.Println("flagged message")
			c.JSON(200, gin.H{"message": "message flagged"})
		} else if int_flagval == 0 {
			log.Println("unflagged message")
			c.JSON(200, gin.H{"message": "message unflagged"})
		}
	}
}
