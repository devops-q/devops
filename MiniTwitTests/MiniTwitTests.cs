public class MiniTwitTests : IDisposable
{
    private const string BASE_URL = "http://localhost:5000";
    private readonly HttpClient client;

    public MiniTwitTests()
    {
        client = new HttpClient();
    }

    public void Dispose()
    {
        client.Dispose();
    }

    private async Task<HttpResponseMessage> Register(string username, string password, string password2 = null,
        string email = null)
    {
        password2 = password2 ?? password;
        email = email ?? $"{username}@example.com";

        var content = new FormUrlEncodedContent(new[]
        {
            new KeyValuePair<string, string>("username", username),
            new KeyValuePair<string, string>("password", password),
            new KeyValuePair<string, string>("password2", password2),
            new KeyValuePair<string, string>("email", email)
        });

        return await client.PostAsync($"{BASE_URL}/register", content);
    }

    private async Task<(HttpResponseMessage, HttpClient)> Login(string username, string password)
    {
        var sessionClient = new HttpClient();
        var content = new FormUrlEncodedContent(new[]
        {
            new KeyValuePair<string, string>("username", username),
            new KeyValuePair<string, string>("password", password)
        });

        var response = await sessionClient.PostAsync($"{BASE_URL}/login", content);
        return (response, sessionClient);
    }

    private async Task<(HttpResponseMessage, HttpClient)> RegisterAndLogin(string username, string password)
    {
        await Register(username, password);
        return await Login(username, password);
    }

    private async Task<HttpResponseMessage> Logout(HttpClient sessionClient)
    {
        return await sessionClient.GetAsync($"{BASE_URL}/logout");
    }

    private async Task<HttpResponseMessage> AddMessage(HttpClient sessionClient, string text)
    {
        var content = new FormUrlEncodedContent(new[]
        {
            new KeyValuePair<string, string>("text", text)
        });

        var response = await sessionClient.PostAsync($"{BASE_URL}/add_message", content);
        if (!string.IsNullOrEmpty(text))
        {
            var responseString = await response.Content.ReadAsStringAsync();
            Assert.Contains("Your message was recorded", responseString);
        }

        return response;
    }

    [Fact]
    public async Task TestRegister()
    {
        var response = await Register("user1", "default");
        var responseString = await response.Content.ReadAsStringAsync();

        Assert.Contains("You were successfully registered and can login now", responseString);

        response = await Register("user1", "default");
        responseString = await response.Content.ReadAsStringAsync();
        Assert.Contains("The username is already taken", responseString);

        response = await Register("", "default");
        responseString = await response.Content.ReadAsStringAsync();
        Assert.Contains("You have to enter a username", responseString);

        response = await Register("meh", "");
        responseString = await response.Content.ReadAsStringAsync();
        Assert.Contains("You have to enter a password", responseString);

        response = await Register("meh", "x", "y");
        responseString = await response.Content.ReadAsStringAsync();
        Assert.Contains("The two passwords do not match", responseString);

        response = await Register("meh", "foo", email: "broken");
        responseString = await response.Content.ReadAsStringAsync();
        Assert.Contains("You have to enter a valid email address", responseString);
    }

    [Fact]
    public async Task TestLoginLogout()
    {
        var (response, sessionClient) = await RegisterAndLogin("user0", "default");
        var responseString = await response.Content.ReadAsStringAsync();
        Assert.Contains("You were logged in", responseString);

        response = await Logout(sessionClient);
        responseString = await response.Content.ReadAsStringAsync();
        Assert.Contains("You were logged out", responseString);

        (response, _) = await Login("user0", "wrongpassword");
        responseString = await response.Content.ReadAsStringAsync();
        Assert.Contains("Invalid password", responseString);

        (response, _) = await Login("user2", "wrongpassword");
        responseString = await response.Content.ReadAsStringAsync();
        Assert.Contains("Invalid username", responseString);
    }

    [Fact]
    public async Task TestMessageRecording()
    {
        var (_, sessionClient) = await RegisterAndLogin("foo", "default");
        await AddMessage(sessionClient, "test message 1");
        await AddMessage(sessionClient, "<test message 2>");
        var response = await client.GetAsync($"{BASE_URL}/");
        var responseString = await response.Content.ReadAsStringAsync();
        Assert.Contains("test message 1", responseString);
        Assert.Contains("&lt;test message 2&gt;", responseString);
    }

    [Fact]
    public async Task TestTimelines()
    {
        var (_, fooSession) = await RegisterAndLogin("foo", "default");
        await AddMessage(fooSession, "the message by foo");
        await Logout(fooSession);

        var (_, barSession) = await RegisterAndLogin("bar", "default");
        await AddMessage(barSession, "the message by bar");

        var response = await barSession.GetAsync($"{BASE_URL}/public");
        var responseString = await response.Content.ReadAsStringAsync();
        Assert.Contains("the message by foo", responseString);
        Assert.Contains("the message by bar", responseString);

        response = await barSession.GetAsync($"{BASE_URL}/");
        responseString = await response.Content.ReadAsStringAsync();
        Assert.DoesNotContain("the message by foo", responseString);
        Assert.Contains("the message by bar", responseString);

        response = await barSession.GetAsync($"{BASE_URL}/foo/follow");
        responseString = await response.Content.ReadAsStringAsync();
        Assert.Contains("You are now following &#34;foo&#34;", responseString);

        response = await barSession.GetAsync($"{BASE_URL}/");
        responseString = await response.Content.ReadAsStringAsync();
        Assert.Contains("the message by foo", responseString);
        Assert.Contains("the message by bar", responseString);

        response = await barSession.GetAsync($"{BASE_URL}/bar");
        responseString = await response.Content.ReadAsStringAsync();
        Assert.DoesNotContain("the message by foo", responseString);
        Assert.Contains("the message by bar", responseString);

        response = await barSession.GetAsync($"{BASE_URL}/foo");
        responseString = await response.Content.ReadAsStringAsync();
        Assert.Contains("the message by foo", responseString);
        Assert.DoesNotContain("the message by bar", responseString);

        response = await barSession.GetAsync($"{BASE_URL}/foo/unfollow");
        responseString = await response.Content.ReadAsStringAsync();
        Assert.Contains("You are no longer following &#34;foo&#34;", responseString);

        response = await barSession.GetAsync($"{BASE_URL}/");
        responseString = await response.Content.ReadAsStringAsync();
        Assert.DoesNotContain("the message by foo", responseString);
        Assert.Contains("the message by bar", responseString);
    }
}