package tests

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

const baseURL = "http://localhost:8181"

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

func TestRegister(t *testing.T) {
	resp := register(t, "user01", "default", "", "")
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "You were successfully registered and can login now")

	resp = register(t, "user01", "default", "", "")
	body, _ = io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "The username is already taken")

	resp = register(t, "", "default", "", "")
	body, _ = io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "You have to enter a username")

	resp = register(t, "meh", "", "", "")
	body, _ = io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "You have to enter a password")

	resp = register(t, "meh", "x", "y", "")
	body, _ = io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "The two passwords do not match")

	resp = register(t, "meh", "foo", "", "broken")
	body, _ = io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "You have to enter a valid email address")
}

func TestLoginLogout(t *testing.T) {
	resp, client := registerAndLogin(t, "user1", "default")
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "You were logged in")

	resp = logout(t, client)
	body, _ = io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "You were logged out")

	resp, client = login(t, "user1", "wrongpassword")
	body, _ = io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Invalid password")

	resp, client = login(t, "user2", "wrongpassword")
	body, _ = io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Invalid username")
}

func TestMessageRecording(t *testing.T) {
	_, client := registerAndLogin(t, "foo", "default")
	addMessage(t, client, "test message 1")
	addMessage(t, client, "<test message 2>")

	resp, err := http.Get(baseURL + "/")
	assert.NoError(t, err)
	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)
	assert.Contains(t, bodyStr, "test message 1")
	assert.Contains(t, bodyStr, "&lt;test message 2&gt;")
}

func TestTimelines(t *testing.T) {
	_, client1 := registerAndLogin(t, "foo", "default")
	addMessage(t, client1, "the message by foo")
	logout(t, client1)

	_, client2 := registerAndLogin(t, "bar", "default")
	addMessage(t, client2, "the message by bar")

	resp, err := client2.Get(baseURL + "/public")
	assert.NoError(t, err)
	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)
	assert.Contains(t, bodyStr, "the message by foo")
	assert.Contains(t, bodyStr, "the message by bar")

	resp, err = client2.Get(baseURL + "/")
	assert.NoError(t, err)
	body, _ = io.ReadAll(resp.Body)
	bodyStr = string(body)
	assert.NotContains(t, bodyStr, "the message by foo")
	assert.Contains(t, bodyStr, "the message by bar")

	resp, err = client2.Get(baseURL + "/foo/follow")
	assert.NoError(t, err)
	body, _ = io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "You are now following &#34;foo&#34;")

	resp, err = client2.Get(baseURL + "/")
	assert.NoError(t, err)
	body, _ = io.ReadAll(resp.Body)
	bodyStr = string(body)
	assert.Contains(t, bodyStr, "the message by foo")
	assert.Contains(t, bodyStr, "the message by bar")

	resp, err = client2.Get(baseURL + "/bar")
	assert.NoError(t, err)
	body, _ = io.ReadAll(resp.Body)
	bodyStr = string(body)
	assert.NotContains(t, bodyStr, "the message by foo")
	assert.Contains(t, bodyStr, "the message by bar")

	resp, err = client2.Get(baseURL + "/foo")
	assert.NoError(t, err)
	body, _ = io.ReadAll(resp.Body)
	bodyStr = string(body)
	assert.Contains(t, bodyStr, "the message by foo")
	assert.NotContains(t, bodyStr, "the message by bar")

	resp, err = client2.Get(baseURL + "/foo/unfollow")
	assert.NoError(t, err)
	body, _ = io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "You are no longer following &#34;foo&#34;")

	resp, err = client2.Get(baseURL + "/")
	assert.NoError(t, err)
	body, _ = io.ReadAll(resp.Body)
	bodyStr = string(body)
	assert.NotContains(t, bodyStr, "the message by foo")
	assert.Contains(t, bodyStr, "the message by bar")
}
