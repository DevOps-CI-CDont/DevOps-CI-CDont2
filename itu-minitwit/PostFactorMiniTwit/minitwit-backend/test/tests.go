package main

import (
	"log"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

const base = "http://localhost:8080"

func main() {
	// run the test
	// r := main.SetUpRouter
	testing.Main(func(pat, str string) (bool, error) { return true, nil }, []testing.InternalTest{
		{
			Name: "TestTimeline",
			F:    TestTimeline,
		},
		{
			Name: "TestRegister",
			F:    TestRegister,
		},
		{
			Name: "TestLoginLogout",
			F:    TestLoginLogout,
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
	http.Get(base + "/RESET")

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
	//make sure logg in in works and out wworks
	res := login("Silas", "password")
	assert.Equal(t, 200, res.StatusCode)
	getHelper("/logout", t, 200, "Logout")

}

func TestLoginLogout(t *testing.T) {
	loginLogout(t)

	//test login fails with wrong password
	postHelper("/login", t, 401, "Login fails with wrong password", url.Values{"username": {"Silas"}, "password": {"wrongpassword"}})
	// with wrong username
	postHelper("/login", t, 401, "Login fails with wrong username", url.Values{"username": {"wrongusername"}, "password": {"password"}})

}

func TestPostTweet(t *testing.T) {
	//Check if adding messages works

}
func getHelper(endpoint string, t *testing.T, expected int, name string) {
	log.Println("Testing " + name)
	res, err := http.Get(base + endpoint)
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

func TestTimeline(t *testing.T) {
	// public & user timeline
	getHelper("/public", t, 200, "Public Timeline")
	// needs authentication
	//getHelper("/", t, 200, "User Timeline")

}
