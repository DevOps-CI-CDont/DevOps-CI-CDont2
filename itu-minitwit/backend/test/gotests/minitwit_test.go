package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const base = "http://localhost:8080"
const sim_url = "http://localhost:8081"

func main() {
	// run the test
	// r := main.SetUpRouter
	testing.Main(func(pat, str string) (bool, error) { return true, nil }, []testing.InternalTest{
		{
			Name: "TestTimelines",
			F:    TestTimelines,
		},
		{
			Name: "TestRegister",
			F:    TestRegister,
		},
		{
			Name: "TestLoginLogout",
			F:    TestLoginLogout,
		},
		{
			Name: "TestPostTweet",
			F:    TestPostTweet,
		},
		{
			Name: "SimTestLatest",
			F:    TestSimLatest,
		},
		{
			Name: "SimTestRegister",
			F:    TestSimRegister,
		},
		{
			Name: "SimTestCreateMsg",
			F:    TestSimCreateMsg,
		},
		{
			Name: "SimTestGetLatestUserMsgs",
			F:    TestSimGetLatestUserMsgs,
		},
	}, []testing.InternalBenchmark{}, []testing.InternalExample{})

}

func login(username string, password string) *http.Response {
	res, err := http.PostForm(base+"/login", url.Values{"username": {username}, "password": {password}})
	if err != nil {
		log.Fatal(err)
	}
	return res
}

func TestRegister(t *testing.T) {
	t.Skip("Skipping register test")
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
	// t.Skip("Skipping login/logout test")
	res := login("Silas", "password")
	assert.Equal(t, 200, res.StatusCode)
	url := base + "/logout"
	getHelper(url, t, 200, "Logout")

}

func TestLoginLogout(t *testing.T) {
	// t.Skip("Skipping login/logout test")
	//test login fails with wrong password
	postHelper("/login", t, 401, "Login fails with wrong password", url.Values{"username": {"Silas"}, "password": {"wrongpassword"}})
	// with wrong username
	postHelper("/login", t, 401, "Login fails with wrong username", url.Values{"username": {"wrongusername"}, "password": {"password"}})

}

func TestPostTweet(t *testing.T) {
	// t.Skip("Skipping post tweet test")
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
	// t.Skip("Skipping timelines test")
	// public & user timeline
	endpoint := base + "/public"
	getHelper(endpoint, t, 401, "Public Timeline") // no tweets -> 401

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
	form := url.Values{}
	form.Add("username", "test")
	form.Add("email", "test")
	form.Add("pwd", "foo")
	fmt.Println("pre TestLatest POST")
	fmt.Println("endpoint", endpoint)
}

func TestSimRegister(t *testing.T) {
}
func TestSimCreateMsg(t *testing.T) {
}

func TestSimGetLatestUserMsgs(t *testing.T) {
}
