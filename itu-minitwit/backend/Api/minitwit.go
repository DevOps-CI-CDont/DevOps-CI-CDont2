package Api

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var Router *gin.Engine

type metrics struct {
	funcCounter *prometheus.CounterVec
}

func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		funcCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "function_calls_total",
			Help: "Number of calls to each function",
		}, []string{"method", "endpont", "code"}),
	}
	reg.MustRegister(m.funcCounter)
	return m
}

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

	// metrics
	reg := prometheus.NewRegistry()
	m := NewMetrics(reg)

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
	Router.GET("/RESET", init_db, incrementCounter(m, "/RESET"))
	Router.GET("/AmIFollowing/:username", amIFollowing, incrementCounter(m, "AmIFollowing/:username"))
	Router.GET("/allUsers", getAllUsers, incrementCounter(m, "/allUsers"))
	Router.Run(":8080")
}

func incrementCounter(m *metrics, endpointName string) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		m.funcCounter.WithLabelValues(c.Request.Method, endpointName, strconv.Itoa(c.Writer.Status())).Inc()
	}
	return fn
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

var DB *sql.DB // global DB variable
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
	db, err := sql.Open("postgres", dbinfo)
	errorCheck(err)

	DB = db

	return nil
}

func init_db(c *gin.Context) {
	const Benjapass = "12345678"
	const Oliverpass = "1234"
	const Silaspass = "password"
	const Januspass = "Janus"

	passwordHashString := HashPassword(Benjapass)
	passwordHashStringO := HashPassword(Oliverpass)
	passwordHashStringS := HashPassword(Silaspass)
	passwordHashStringJ := HashPassword(Januspass)

	// create tables
	_, err := DB.Exec(`
		drop table if exists users;
	`)
	errorCheck(err)

	_, err = DB.Exec(`
		drop table if exists messages;
	`)
	errorCheck(err)

	_, err = DB.Exec(`
		drop table if exists followers;
	`)
	errorCheck(err)

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			user_id serial,
			username text,
			email text,
			pw_hash text,
			PRIMARY KEY (user_id)
		);
	`)
	errorCheck(err)
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
			message_id serial,
			author_id integer,
			text text,
			pub_date integer,
			flagged integer,
			PRIMARY KEY (message_id)
		);
	`)
	errorCheck(err)

	_, err = DB.Exec(`
		create table if not exists followers (who_id integer, whom_id integer);
	`)
	errorCheck(err)

	_, err = DB.Exec(`
		INSERT INTO users (username, email, pw_hash)
		VALUES ('Benjamin', 'bekj@itu.dk', $1);
	`, passwordHashString)
	errorCheck(err)

	_, err = DB.Exec(`
		INSERT INTO users (username, email, pw_hash)
		VALUES ('Oliver', 'ojoe@itu.dk', $1);
	`, passwordHashStringO)
	errorCheck(err)

	_, err = DB.Exec(`
		INSERT INTO users (username, email, pw_hash)
		VALUES ('Silas', 'sipn@itu.dk', $1);
	`, passwordHashStringS)
	errorCheck(err)

	_, err = DB.Exec(`
		INSERT INTO users (username, email, pw_hash)
		VALUES ('Janus', 'januh@itu.dk', $1);
	`, passwordHashStringJ)
	errorCheck(err)

	_, err = DB.Exec(`
		INSERT INTO messages (author_id, text, pub_date, flagged)
		VALUES (1, 'I like apples', 123456789, 0);
	`)
	errorCheck(err)

	_, err = DB.Exec(`
		INSERT INTO messages (author_id, text, pub_date, flagged)
		VALUES (2, 'I like tarteletter', 123456789, 0);
	`)
	errorCheck(err)

	_, err = DB.Exec(`
		INSERT INTO messages (author_id, text, pub_date, flagged)
		VALUES (3, 'I like Pizza', 123456789, 0);
	`)
	errorCheck(err)

	_, err = DB.Exec(`
		INSERT INTO messages (author_id, text, pub_date, flagged)
		VALUES (4, 'I like bananas ', 123456789, 0);
	`)
	errorCheck(err)

	sqlStmt9 := `
		INSERT INTO followers (who_id, whom_id)
		VALUES
		(1, 2),
		(1, 3),
		(1, 4),
		(2, 1),
		(2, 3),
		(2, 4),
		(3, 1),
		(3, 2),
		(3, 4),
		(4, 1),
		(4, 2);
	`
	_, err = DB.Exec(sqlStmt9)
	errorCheck(err)

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
	rows, err := DB.Query(`select * from followers
		where who_id = $1 and whom_id = (select user_id from users where username = $2)`, userID, username)
	errorCheck(err)
	defer rows.Close()
	following := false
	for rows.Next() {
		following = true
	}
	c.JSON(200, following)
}

func getTimeline(c *gin.Context) {
	//query database
	// check cookie for session,
	userID := getUserIdIfLoggedIn(c)

	rows, err := DB.Query(`select messages.*, users.* from messages, users
		where messages.flagged = 0 and messages.author_id = users.user_id and (
		users.user_id = $1 or
		users.user_id in (select whom_id from followers
		where who_id = $2))
		order by messages.pub_date desc limit $3`, userID, userID, PER_PAGE)
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
	num_msgs := c.Request.URL.Query().Get("num_msgs")
	int_num_msgs, err := strconv.Atoi(num_msgs)
	if num_msgs == "" || err != nil {
		int_num_msgs = 30
	}

	fmt.Println("int_num_msgs", int_num_msgs)

	rows, err := DB.Query(`select messages.*, users.* from messages, users
	where messages.flagged = 0 and messages.author_id = users.user_id
	order by messages.pub_date desc limit $1`, int_num_msgs)
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

	errorCheck(err)

	// if no messages, return 401
	if len(messages) == 0 {
		c.JSON(401, gin.H{"message": "no messages"})
	}

	fmt.Println("messages", messages)

	c.JSON(200, gin.H{"tweets": messages})

	defer rows.Close()
}

func getUsersTweets(c *gin.Context) {
	name := c.Param("username")
	userID := GetUserIdByName(name)
	if userID == "-1" {
		c.JSON(200, gin.H{"message": "user does not exist"})
		return
	}

	num_msgs := c.Request.URL.Query().Get("num_msgs")
	int_num_msgs, err := strconv.Atoi(num_msgs)
	if num_msgs == "" || err != nil {
		int_num_msgs = 30
	}

	rows, err := DB.Query(`select messages.*, users.* from messages, users where messages.author_id = $1 and messages.author_id = users.user_id order by messages.pub_date desc limit $2`, userID, int_num_msgs)
	errorCheck(err)
	defer rows.Close()
	messages := make([]Message, 0)
	for rows.Next() {
		msg := Message{}
		user := User{}
		err = rows.Scan(&msg.Message_id, &msg.Author_id, &msg.Text, &msg.Pub_date, &msg.Flagged, &user.User_id, &user.Username, &user.Email, &user.Pw_hash)
		errorCheck(err)
		msg.Author = user

		messages = append(messages, msg)
	}
	c.JSON(200, gin.H{"tweets": messages})

}

func followUser(c *gin.Context) {

	userid := getUserIdIfLoggedIn(c)
	whom_name := c.Param("username")
	whom_id := GetUserIdByName(whom_name)
	if doesUsersFollow(userid, whom_id) {
		c.JSON(200, gin.H{"message": "user already followed"})
		return
	}

	if whom_id == "-1" {
		c.JSON(200, gin.H{"message": "user does not exist"})
		return
	}
	stmt, err := DB.Prepare(`insert into followers (who_id, whom_id) values ($1, $2)`)
	errorCheck(err)
	defer stmt.Close()

	_, err = stmt.Exec(userid, whom_id)
	errorCheck(err)

	c.JSON(200, gin.H{"message": "followed user"})
}

func doesUsersFollow(who_id string, whom_id string) bool {
	row := DB.QueryRow(`select * from followers where who_id = $1 and whom_id = $2`, who_id, whom_id)

	follower := follower{}
	err := row.Scan(&follower.Who_id, &follower.Whom_id)
	return err == nil

}

func unfollowUser(c *gin.Context) {

	userid := getUserIdIfLoggedIn(c)
	whom_name := c.Param("username")
	whom_id := GetUserIdByName(whom_name)
	if !doesUsersFollow(userid, whom_id) {
		c.JSON(200, gin.H{"message": "user doesn't follow the target"})
		return
	}
	log.Println(whom_id)
	if whom_id == "-1" {
		c.JSON(200, gin.H{"message": "user you are trying to follow does not exist"})
		return
	}

	stmt, err := DB.Prepare(`Delete FROM followers WHERE who_id = $1 and whom_id = $2`)
	errorCheck(err)
	defer stmt.Close()

	_, err = stmt.Exec(userid, whom_id)
	errorCheck(err)

	c.JSON(200, gin.H{"message": "unfollowed user"})

}

func postMessage(c *gin.Context) {

	userid := getUserIdIfLoggedIn(c)

	if userid == "-1" {
		c.JSON(401, gin.H{"message": "user not logged in"})
		return
	}

	text := c.PostForm("text")
	authorid := userid
	flagged := 0
	log.Println("text:" + text)

	stmt, err := DB.Prepare(`insert into messages (author_id, text, pub_date, flagged) values ($1, $2, $3, $4)`)
	errorCheck(err)
	defer stmt.Close()

	_, err = stmt.Exec(authorid, text, time.Now().Unix(), flagged)
	errorCheck(err)

	c.JSON(200, gin.H{"message": "message posted"})

}

func login(c *gin.Context) {

	username := c.PostForm("username")
	password := c.PostForm("password")
	//convert password to byte[]

	passwordHash := HashPassword(password)

	//check if username and password are pty
	if username == "" || password == "" {
		c.JSON(400, gin.H{"error": "username or password is empty"})
		return
	}

	userid := DB.QueryRow(`select user_id from users where users.Username = $1 and users.pw_hash = $2`, username, passwordHash)

	var userIdAsInt int
	err := userid.Scan(&userIdAsInt)
	if userIdAsInt == 0 {
		c.SetCookie("user_id", "", -1, "/", "localhost", false, false)
		c.JSON(401, gin.H{"error": "username or password is incorrect"})
		return
	}
	if err != nil {
		c.SetCookie("user_id", "", -1, "/", "localhost", false, false)
		c.JSON(401, gin.H{"error": "username or password is incorrect"})
		return
	}
	// succes: set cookie
	c.SetCookie("user_id", strconv.Itoa(userIdAsInt), 3600, "/", "localhost", false, false)
	c.JSON(200, gin.H{"user_id": userIdAsInt})
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

	stmt, err := DB.Prepare(`insert into users (username, email, pw_hash) values ($1, $2, $3)`)
	errorCheck(err)
	defer stmt.Close()
	_, err = stmt.Exec(username, email, passwordHashString)
	errorCheck(err)

	c.JSON(200, gin.H{"message": "user registered"})
}

func getUserByName(userName string) *sql.Row {
	row := DB.QueryRow(`select * from users where users.username = $1`, userName)
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
	log.Println("cookie user_id: " + userid)
	errorCheck(err)
	if userid == "" || userid == "-1" {
		c.JSON(401, gin.H{"error": "not logged in"})
		return "-1"
	}
	return userid

}

func GetUserIdByName(username string) string {
	stmt, err := DB.Prepare("SELECT user_id FROM users WHERE username = $1")
	errorCheck(err)
	defer stmt.Close()
	var userId string
	err = stmt.QueryRow(username).Scan(&userId)
	errorCheck(err)
	fmt.Println("userId: " + userId)
	return userId
}

func logout(c *gin.Context) {
	c.SetCookie("user_id", "", -1, "/", "localhost", false, false)
}

func getAllUsers(c *gin.Context) {
	rows, err := DB.Query(`select * from users`)
	errorCheck(err)
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.User_id, &user.Username, &user.Email, &user.Pw_hash)
		errorCheck(err)
		users = append(users, user)
	}
	c.JSON(200, gin.H{
		"users": users,
	})
}
