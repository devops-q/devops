package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os/exec"
	"testing"
)

const host = "http://127.0.0.1"
const port = "8181"
const baseURL = host + ":" + port

type MinitwitTestSuite struct {
	suite.Suite
	Cmd exec.Cmd
}

func TestMinitwitTestSuite(t *testing.T) {
	suite.Run(t, new(MinitwitTestSuite))
}

// helper functions
func register(t *testing.T, username, password, password2, email string) *http.Response {
	jar, err := cookiejar.New(nil)
	assert.NoError(t, err)
	client := &http.Client{Jar: jar}
	if password2 == "" {
		password2 = password
	}
	if email == "" {
		email = username + "@example.com"
	}

	form := url.Values{}
	form.Add("username", username)
	form.Add("password", password)
	form.Add("password2", password2)
	form.Add("email", email)

	resp, err := client.PostForm(baseURL+"/register", form)
	assert.NoError(t, err)
	return resp
}

func login(t *testing.T, username string, password string) (*http.Response, *http.Client) {
	jar, err := cookiejar.New(nil)
	assert.NoError(t, err)
	client := &http.Client{Jar: jar}

	form := url.Values{}

	form.Add("username", username)
	form.Add("password", password)

	resp, err := client.PostForm(baseURL+"/login", form)
	assert.NoError(t, err)
	return resp, client
}

func registerAndLogin(t *testing.T, username, password string) (*http.Response, *http.Client) {
	register(t, username, password, "", "")
	return login(t, username, password)
}

func logout(t *testing.T, client *http.Client) *http.Response {
	resp, err := client.Get(baseURL + "/logout")
	assert.NoError(t, err)
	return resp
}

func addMessage(t *testing.T, client *http.Client, text string) *http.Response {
	form := url.Values{}
	form.Add("text", text)

	resp, err := client.PostForm(baseURL+"/add_message", form)
	assert.NoError(t, err)

	if text != "" {
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Contains(t, string(body), "Your message was recorded")
	}
	return resp
}

func (suite *MinitwitTestSuite) TestRegister() {
	resp := register(suite.T(), "user01", "default", "", "")
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(suite.T(), string(body), "You were successfully registered and can login now")

	resp = register(suite.T(), "user01", "default", "", "")
	body, _ = io.ReadAll(resp.Body)
	assert.Contains(suite.T(), string(body), "The username is already taken")

	resp = register(suite.T(), "", "default", "", "")
	body, _ = io.ReadAll(resp.Body)
	assert.Contains(suite.T(), string(body), "You have to enter a username")

	resp = register(suite.T(), "meh", "", "", "")
	body, _ = io.ReadAll(resp.Body)
	assert.Contains(suite.T(), string(body), "You have to enter a password")

	resp = register(suite.T(), "meh", "x", "y", "")
	body, _ = io.ReadAll(resp.Body)
	assert.Contains(suite.T(), string(body), "The two passwords do not match")

	resp = register(suite.T(), "meh", "foo", "", "broken")
	body, _ = io.ReadAll(resp.Body)
	assert.Contains(suite.T(), string(body), "You have to enter a valid email address")
}

func (suite *MinitwitTestSuite) TestLoginLogout() {
	resp, client := registerAndLogin(suite.T(), "user1", "default")
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(suite.T(), string(body), "You were logged in")

	resp = logout(suite.T(), client)
	body, _ = io.ReadAll(resp.Body)
	assert.Contains(suite.T(), string(body), "You were logged out")

	resp, _ = login(suite.T(), "user1", "wrongpassword")
	body, _ = io.ReadAll(resp.Body)
	assert.Contains(suite.T(), string(body), "Invalid password")

	resp, _ = login(suite.T(), "user2", "wrongpassword")
	body, _ = io.ReadAll(resp.Body)
	assert.Contains(suite.T(), string(body), "Invalid username")
}

func (suite *MinitwitTestSuite) TestMessageRecording() {
	_, client := registerAndLogin(suite.T(), "foo", "default")
	addMessage(suite.T(), client, "test message 1")
	addMessage(suite.T(), client, "<test message 2>")

	resp, err := http.Get(baseURL + "/")
	assert.NoError(suite.T(), err)
	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)
	assert.Contains(suite.T(), bodyStr, "test message 1")
	assert.Contains(suite.T(), bodyStr, "&lt;test message 2&gt;")
}

func (suite *MinitwitTestSuite) TestTimelines() {
	_, client1 := registerAndLogin(suite.T(), "foo", "default")
	addMessage(suite.T(), client1, "the message by foo")
	logout(suite.T(), client1)

	_, client2 := registerAndLogin(suite.T(), "bar", "default")
	addMessage(suite.T(), client2, "the message by bar")

	resp, err := client2.Get(baseURL + "/public")
	assert.NoError(suite.T(), err)
	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)
	assert.Contains(suite.T(), bodyStr, "the message by foo")
	assert.Contains(suite.T(), bodyStr, "the message by bar")

	resp, err = client2.Get(baseURL + "/")
	assert.NoError(suite.T(), err)
	body, _ = io.ReadAll(resp.Body)
	bodyStr = string(body)
	assert.NotContains(suite.T(), bodyStr, "the message by foo")
	assert.Contains(suite.T(), bodyStr, "the message by bar")

	resp, err = client2.Get(baseURL + "/foo/follow")
	assert.NoError(suite.T(), err)
	body, _ = io.ReadAll(resp.Body)
	assert.Contains(suite.T(), string(body), "You are now following &#34;foo&#34;")

	resp, err = client2.Get(baseURL + "/")
	assert.NoError(suite.T(), err)
	body, _ = io.ReadAll(resp.Body)
	bodyStr = string(body)
	assert.Contains(suite.T(), bodyStr, "the message by foo")
	assert.Contains(suite.T(), bodyStr, "the message by bar")

	resp, err = client2.Get(baseURL + "/bar")
	assert.NoError(suite.T(), err)
	body, _ = io.ReadAll(resp.Body)
	bodyStr = string(body)
	assert.NotContains(suite.T(), bodyStr, "the message by foo")
	assert.Contains(suite.T(), bodyStr, "the message by bar")

	resp, err = client2.Get(baseURL + "/foo")
	assert.NoError(suite.T(), err)
	body, _ = io.ReadAll(resp.Body)
	bodyStr = string(body)
	assert.Contains(suite.T(), bodyStr, "the message by foo")
	assert.NotContains(suite.T(), bodyStr, "the message by bar")

	resp, err = client2.Get(baseURL + "/foo/unfollow")
	assert.NoError(suite.T(), err)
	body, _ = io.ReadAll(resp.Body)
	assert.Contains(suite.T(), string(body), "You are no longer following &#34;foo&#34;")

	resp, err = client2.Get(baseURL + "/")
	assert.NoError(suite.T(), err)
	body, _ = io.ReadAll(resp.Body)
	bodyStr = string(body)
	assert.NotContains(suite.T(), bodyStr, "the message by foo")
	assert.Contains(suite.T(), bodyStr, "the message by bar")
}
