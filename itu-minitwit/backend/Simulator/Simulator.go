package simulator

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	main "minitwit-backend/init/Api"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type FilteredMessage struct {
	Content  string  `json:"content"`
	Pub_date float64 `json:"pub_date"`
	User     string  `json:"user"`
}

var LATEST = 0

func update_latest(c *gin.Context) {
	try_latest := c.Request.URL.Query().Get("latest")
	int_val, err := strconv.Atoi(try_latest)
	if err != nil {
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

var api_base_url, is_present = os.LookupEnv("API_BASE_URL")

func Start() {
	Router := SetUpRouter()

	if !is_present {
		api_base_url = "http://0.0.0.0:8080"
	}

	// router config
	Router.Use(cors.Default()) // cors.Default() should allow all origins
	// it's important to set this before any routes are registered so that the middleware is applied to all routes
	// ALL MY HOMIES HATE CORS :D

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
	req, err := http.NewRequest("POST", api_base_url+"/register", nil)
	if err != nil {
		c.JSON(400, gin.H{"error_msg": err.Error()})
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
		c.JSON(400, gin.H{"error_msg": err.Error()})
		return
	}

	if resp.StatusCode != 200 {
		c.JSON(400, gin.H{"error_msg": "Something went wrong with the registration!"})
		return
	}
	defer resp.Body.Close()
	c.JSON(204, gin.H{}) // no content
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
	fmt.Println("getMsgs")
	num_msgs := getNumMsgs(c)
	fmt.Println("num_msgs: ", num_msgs)

	// create a new HTTP client
	client := &http.Client{}

	url := api_base_url + "/public?num_msgs=" + strconv.Itoa(num_msgs)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.JSON(400, gin.H{"error_msg": err.Error()})
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)

	if err != nil {
		c.JSON(400, gin.H{"error_msg": err.Error()})
		return
	}

	// Read the response body as JSON
	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// make filteredMsgs
	filteredMsgs := make([]FilteredMessage, 0)
	for _, msg := range data["tweets"].([]interface{}) {
		msg := msg.(map[string]interface{})
		author := msg["author"].(map[string]interface{})
		filteredMsgs = append(filteredMsgs, FilteredMessage{
			Content:  msg["text"].(string),
			Pub_date: msg["pub_date"].(float64),
			User:     author["username"].(string),
		})
	}

	// Marshal the filteredMsgs slice into a JSON-encoded byte slice
	jsonBytes, err := json.Marshal(filteredMsgs)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Set the response header and write the JSON-encoded byte slice to the response writer
	c.Header("Content-Type", "application/json")
	c.Writer.Write(jsonBytes)
}

func getNumMsgs(c *gin.Context) int {
	// default
	num_msgs := c.Request.URL.Query().Get("no")
	if num_msgs == "" {
		int_num_msgs := 30
		return int_num_msgs
	} else {
		int_num_msgs, err := strconv.Atoi(num_msgs)
		if err != nil {
			int_num_msgs = 30
		}
		return int_num_msgs
	}
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
	num_msgs := getNumMsgs(c)

	if c.Request.Method == "GET" {
		user_name := c.Param("username")

		if user_name == "" {
			c.AbortWithStatus(404)
			return
		}

		//CREATE NEW CLIENT
		client := &http.Client{}

		// create get request
		url := api_base_url + "/user/" + user_name + "?num_msgs=" + strconv.Itoa(num_msgs)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			c.JSON(400, gin.H{"error_msg": err.Error()})
			return
		}

		// send the request
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(400, gin.H{"error_msg": err.Error()})
			return
		}

		// Read the response body as JSON
		var data map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// make filteredMsgs
		filteredMsgs := make([]FilteredMessage, 0)
		for _, msg := range data["tweets"].([]interface{}) {
			msg := msg.(map[string]interface{})
			author := msg["author"].(map[string]interface{})
			filteredMsgs = append(filteredMsgs, FilteredMessage{
				Content:  msg["text"].(string),
				Pub_date: msg["pub_date"].(float64),
				User:     author["username"].(string),
			})
		}

		// Marshal the filteredMsgs slice into a JSON-encoded byte slice
		jsonBytes, err := json.Marshal(filteredMsgs)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// Set the response header and write the JSON-encoded byte slice to the response writer
		c.Header("Content-Type", "application/json")
		c.Writer.Write(jsonBytes)

	} else if c.Request.Method == "POST" {
		bytes, _ := io.ReadAll(c.Request.Body)
		body := make(map[string]string)
		json.Unmarshal(bytes, &body)

		query := `INSERT INTO messages (author_id, text, pub_date, flagged)
					VALUES ($1, $2, $3, 0)`

		author_id := main.GetUserIdByName(c.Param("username"))
		if author_id == "-1" {
			fmt.Println("non-existing user tried to post a message: " + body["content"])
			c.JSON(404, gin.H{})
			return
		}

		main.DB.Exec(query, main.GetUserIdByName(c.Param("username")), body["content"], time.Now().Unix())
		fmt.Println("user " + c.Param("username") + " posted a message: " + body["content"])
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
	user_name := c.Param("username")
	fmt.Println("username " + user_name)
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

		fmt.Println("user " + user_id + " tries to follow " + follows_user_id)

		if follows_user_id == "-1" {
			c.AbortWithStatus(404)
			return
		}

		query := "INSERT INTO followers (who_id, whom_id) VALUES ($1, $2)"
		main.DB.Exec(query, user_id, follows_user_id)

		c.JSON(204, gin.H{})
	} else if c.Request.Method == "POST" && body["unfollow"] != "" {
		unfollows_username := body["unfollow"]
		unfollows_user_id := main.GetUserIdByName(unfollows_username)

		fmt.Println("user " + user_id + " tries to unfollow " + unfollows_user_id)

		if unfollows_user_id == "-1" {
			c.AbortWithStatus(404)
			return
		}

		query := "DELETE FROM followers WHERE who_id=$1 and WHOM_ID=$2"
		main.DB.Exec(query, user_id, unfollows_user_id)

		c.JSON(204, gin.H{})
	} else if c.Request.Method == "GET" {
		// default
		num_followers := getNumMsgs(c)
		query := `SELECT users.username FROM users
					INNER JOIN followers ON followers.whom_id=users.user_id
					WHERE followers.who_id=$1
					LIMIT $2`

		followers, _ := main.DB.Query(query, user_id, num_followers)
		followers_names := []string{}

		for followers.Next() {
			var username string
			followers.Scan(&username)
			followers_names = append(followers_names, username)
		}

		c.JSON(200, gin.H{
			"follows": followers_names,
		})
	}

}
