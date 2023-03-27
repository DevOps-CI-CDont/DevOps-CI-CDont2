package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"minitwit-backend/init/config"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const base = "http://localhost:8080"
const sim_url = "http://localhost:8081"

func clearTestDB() {
	config.Connect_test_db()
	// delete from messages
	config.DB.Exec("DELETE FROM messages")
	// delete from users
	config.DB.Exec("DELETE FROM users")
	// delete from followers
	config.DB.Exec("DELETE FROM followers")
	// alter sequences
	config.DB.Exec("ALTER SEQUENCE messages_id_seq RESTART WITH 1")
	config.DB.Exec("ALTER SEQUENCE users_id_seq RESTART WITH 1")
	config.DB.Exec("ALTER SEQUENCE followers_id_seq RESTART WITH 1")
	fmt.Println("Test database tables & sequences reset")
}

func encodeJsonAndPOST(t *testing.T, endpoint string, data map[string]string) *http.Response { // utility function to encode data to JSON and POST to endpoint
	// encode the data as a JSON string
	jsonData, _ := json.Marshal(data)

	// create the request with the JSON-encoded data
	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	req.Header.Add("Authorization", "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh")

	// create client
	client := &http.Client{}

	// send request
	resp, _ := client.Do(req)

	return resp
}

func checkLatest(t *testing.T, expectedLatest int) { // utility function to check "latest"
	endpoint := fmt.Sprintf("%s/latest", sim_url)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Fatalf("latest check failed")
	}
	// create client
	client := &http.Client{}
	// send request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("latest check failed")
	}
	assert.Equal(t, 200, resp.StatusCode)
	// read response body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("latest check failed")
	}

	// parse JSON
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatalf("failed to parse JSON response")
	}
	log.Println("data parsed: ", data)

	// extract "latest" value
	latest, ok := data["latest"].(float64)
	if !ok {
		log.Fatalf("failed to extract latest value")
	}
	assert.Equal(t, expectedLatest, int(latest))
}

func TestMain(m *testing.M) {
	// put your setup code here
	clearTestDB()

	// run the tests
	exit_code := m.Run()

	fmt.Println("code:", exit_code)
}

func login(username string, password string) *http.Response {
	res, err := http.PostForm(base+"/login", url.Values{"username": {username}, "password": {password}})
	if err != nil {
		log.Fatal(err)
	}
	return res
}

func TestRegister(t *testing.T) {
	username := "testuser"
	password := "testpassword"
	email := "test@example.com"
	// register
	postHelper("/register", t, 200, "Test register works", url.Values{"username": {username}, "password": {password}, "password2": {password}, "email": {email}})

	postHelper("/register", t, 400, "Username already taken", url.Values{"username": {username}, "password": {password}, "password2": {password}, "email": {email}})

	postHelper("/register", t, 400, "Has to have username", url.Values{"username": {""}, "password": {password}, "password2": {password}, "email": {email}})
	postHelper("/register", t, 400, "Has to have password", url.Values{"username": {username}, "password": {""}, "password2": {password}, "email": {email}})
	//passwords must match
	postHelper("/register", t, 400, "Passwords must match", url.Values{"username": {username}, "password": {password}, "password2": {"notpassword"}, "email": {email}})

	//must enter valid email address
	postHelper("/register", t, 400, "Must enter valid email address", url.Values{"username": {username}, "password": {password}, "password2": {password}, "email": {"notanemail"}})

}

func loginLogout(t *testing.T) {
	res := login("Silas", "password")
	assert.Equal(t, 200, res.StatusCode)
	url := base + "/logout"
	getHelper(url, t, 200, "Logout")

}

func TestLoginLogout(t *testing.T) {
	//test login fails with wrong password
	postHelper("/login", t, 401, "Login fails with wrong password", url.Values{"username": {"Silas"}, "password": {"wrongpassword"}})
	// with wrong username
	postHelper("/login", t, 401, "Login fails with wrong username", url.Values{"username": {"wrongusername"}, "password": {"password"}})

}

func TestPostTweet(t *testing.T) {
	//Check if adding messages works
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}
	form := url.Values{}
	form.Add("text", "hello world")
	res, err := http.NewRequest("POST", base+"/add_message", strings.NewReader(form.Encode()))
	res.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	cookie := &http.Cookie{
		Name:  "user_id",
		Value: "1",
	}
	res.AddCookie(cookie)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(res)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, 200, resp.StatusCode)
	log.Println("Testing add message passed")

}
func getHelper(url string, t *testing.T, expected int, name string) {
	log.Println("Testing " + name)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, expected, res.StatusCode)
	log.Println("Testing " + name + " passed")
}

func postHelper(endpoint string, t *testing.T, expected int, name string, data url.Values) {

	res, err := http.PostForm(base+endpoint, data)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, expected, res.StatusCode)
	log.Println("Testing " + name + " passed")
}

func TestTimelines(t *testing.T) {
	// public & user timeline
	endpoint := base + "/public"
	getHelper(endpoint, t, 200, "Public Timeline") // no tweets -> 401

	// personal timeline: needs authentication
	req, err := http.NewRequest("GET", base+"/mytimeline", nil)
	if err != nil {
		log.Fatal(err)
	}
	cookie := &http.Cookie{
		Name:  "user_id",
		Value: "1",
	}
	req.AddCookie(cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//create client
	client := &http.Client{}
	//send request
	resp, errM := client.Do(req)
	if errM != nil {
		t.Fatalf("mytimeline testfailed")
	}
	assert.Equal(t, 200, resp.StatusCode)
}
func TestSimLatest(t *testing.T) {
	// post something to update LATEST
	endpoint := fmt.Sprintf("%s/register?latest=1337", sim_url)
	data := map[string]string{
		"username": "test",
		"email":    "test@test.com",
		"pwd":      "foo",
	}
	resp := encodeJsonAndPOST(t, endpoint, data)
	assert.Equal(t, 204, resp.StatusCode)
	checkLatest(t, 1337)
}

/* SIMULATOR TESTS BELOW */

func TestSimRegister(t *testing.T) {
	endpoint := fmt.Sprintf("%s/register?latest=1", sim_url)
	data := map[string]string{
		"username": "a",
		"email":    "a@a.a",
		"pwd":      "a",
	}
	resp := encodeJsonAndPOST(t, endpoint, data)
	assert.Equal(t, 204, resp.StatusCode)
	checkLatest(t, 1)

	// check if user exists
	endpoint = fmt.Sprintf("%s/user/a", base)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", endpoint, nil)
	resp, _ = client.Do(req)
	assert.Equal(t, 200, resp.StatusCode)

}
func TestSimCreateMsg(t *testing.T) {
	endpoint := fmt.Sprintf("%s/msgs/a?latest=2", sim_url)
	data := map[string]string{
		"content": "Blub!",
	}
	resp := encodeJsonAndPOST(t, endpoint, data)
	assert.Equal(t, 204, resp.StatusCode)
	checkLatest(t, 2)
}

func TestSimGetLatestUserMsgs(t *testing.T) {

}
