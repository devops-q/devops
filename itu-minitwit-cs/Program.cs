using System.Data.Common;
using System.Security.Cryptography;
using System.Text;
using System.Text.RegularExpressions;
using Microsoft.AspNetCore;
using Microsoft.AspNetCore.Http.Extensions;
using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Mvc.ModelBinding.Validation;
using Microsoft.Data.Sqlite;
using Microsoft.Extensions.FileProviders;
using Microsoft.Extensions.Primitives;
using Scriban;


using Scriban.Parsing;
using Scriban.Runtime;

string DATABASE = "../minitwit.db";
int PER_page = 30;
bool DEBUG = true;
string SECRET_KEY = "development key";

var builder = WebApplication.CreateBuilder(args);

// Add session support
builder.Services.AddDistributedMemoryCache();
builder.Services.AddSession(options =>
{
  options.IdleTimeout = TimeSpan.FromMinutes(30); // Adjust as needed
  options.Cookie.HttpOnly = true;
  options.Cookie.IsEssential = true;
});



var app = builder.Build();

app.UseSession();
app.UseStaticFiles(); // Enable serving static files like CSS


SqliteConnection ConnectDb()
{
    // Returns a new connection to the database
    var connection = new SqliteConnection("Data Source=" + DATABASE);
    connection.Open();

    return connection;
}

List<Dictionary<string, object>> QueryDb(SqliteConnection db, string query, object[] args = null, bool one = false)
{
  using (var command = db.CreateCommand())
  {
    command.CommandText = query;
    if (args != null)
    {
      for (int i = 0; i < args.Length; i++)
      {
        command.Parameters.AddWithValue($"@p{i}", args[i]);
      }
    }

    using (var reader = command.ExecuteReader())
    {
      var result = new List<Dictionary<string, object>>();
      while (reader.Read())
      {
        var row = new Dictionary<string, object>();
        for (int i = 0; i < reader.FieldCount; i++)
        {
          row[reader.GetName(i)] = reader.IsDBNull(i) ? null : reader.GetValue(i);
        }

        result.Add(row);
      }

      if (one)
      {
        return result.Count > 0 ? new List<Dictionary<string, object>> { result[0] } : null;
      }

      return result;
    }
  }
}

object? get_user_id(string username)
{
    // TODO implement method (issue #17)
    throw new NotImplementedException();
}

string FormatDatetime(int timestamp)
{
  // Convert a unix timestamp (seconds) to a human-readable date string.
  var datetime = DateTimeOffset.FromUnixTimeSeconds(timestamp);
  return datetime.ToString("yyyy-MM-dd @ HH:mm");
}

void BeforeRequest(HttpContext context)
{
  // Make sure we are connected to the database each request and look
  // up the current user so that we know he's there.
  context.Items["db"] = ConnectDb();
  context.Items["user"] = null;

  if (context.Session.TryGetValue("user_id", out var userIdBytes))
  {
    var userId = Encoding.UTF8.GetString(userIdBytes);
    var user = QueryDb((SqliteConnection)context.Items["db"], "SELECT * FROM user WHERE user_id = @p0",
      new object[] { userId }, one: true);
    context.Items["user"] = user;
  }
}

void AfterRequest(HttpContext context)
{
    // Closes the database again at the end of the request.
    var db = (SqliteConnection)context.Items["db"];
    db?.Close();
}


app.Use(async (context, next) =>
{
    BeforeRequest(context);
    await next.Invoke();
    AfterRequest(context);
});

// Timeline route
app.MapGet("/", (HttpRequest request, HttpContext context) =>
    timeline(request, context));

IResult timeline(HttpRequest request, HttpContext context)
{
  var db = new SqliteConnection("Data source=" + DATABASE);
  db.Open();

  var userIDFromSession = "1";

  var command = db.CreateCommand();
  Console.WriteLine("We got a visitor from: " + request.HttpContext.Connection.RemoteIpAddress?.ToString());

  if (string.IsNullOrEmpty(userIDFromSession))
    return Results.Redirect("/public");

  var query = @"
        SELECT message.*, user.* FROM message, user
        WHERE message.author_id = user.user_id AND (
            user.user_id = @UserId OR
            user.user_id IN (SELECT whom_id FROM follower WHERE who_id = @UserId))
        ORDER BY message.pub_date DESC LIMIT @PerPage";

  command.CommandText = query;
  command.Parameters.Add(new SqliteParameter("@UserId", userIDFromSession));
  command.Parameters.Add(new SqliteParameter("@PerPage", PER_page));

  List<Dictionary<string, string>> messages = new List<Dictionary<string, string>>();

  using (var reader = command.ExecuteReader())
  {
    while (reader.Read())
    {
      Dictionary<string, string> dict = new Dictionary<string, string>();

      for (int i = 0; i < reader.FieldCount; i++)
      {

        string key = reader.GetName(i);

        string value = reader.IsDBNull(i) ? "" : reader.GetString(i);
        if (key.Equals("pub_date"))
        {
          var dateTimeOffset = DateTimeOffset.FromUnixTimeSeconds(long.Parse(value));
          value = dateTimeOffset.DateTime.ToString("yyyy-MM-dd HH:mm:ss");
        }
        dict[key] = value;
      }

      messages.Add(dict);
    }
  }

 var messagesWithImages = messages.Select(message => new Dictionary<string, object>
  {
    ["username"] = message["username"],
    ["text"] = message["text"],
    ["pub_date"] = message["pub_date"],
    ["email"] = message["email"],
    ["image_url"] = GetGravatarUrl(message["username"].ToString()) // Add a generated image URL
  }).ToList();

  // Data dictionary for template
  var data = new Dictionary<string, object>
  {
    ["title"] = "My Timeline",
    ["messages"] = messagesWithImages,
    ["endpoint"] = request.Path,
    ["userIDFromSession"] = new Dictionary<string, string>
    {
      ["user_id"] = userIDFromSession,
      ["username"] = "TestUser"
    },
    ["profile_user"] = new Dictionary<string, string>
    {
      ["user_id"] = "2",
      ["username"] = "AnotherUser"
    },
    ["followed"] = false
  };



  var finalRenderedHTML = sendToHtml(data, "timeline");

  return Results.Content(finalRenderedHTML, "text/html; charset=utf-8");
}

app.MapGet("/public", (HttpRequest request, HttpContext context) =>
  public_timeline(request, context));

IResult public_timeline(HttpRequest request, HttpContext context)
{


  var db = new SqliteConnection("Data source=" + DATABASE);
  db.Open();

  var command = db.CreateCommand();



  var query = @"
        SELECT message.*, user.* from message, 
        user where message.flagged = 0 and 
        message.author_id = user.user_id
        order by message.pub_date desc limit @PerPage";


  command.CommandText = query;
  command.Parameters.Add(new SqliteParameter("@PerPage", PER_page));

  List<Dictionary<string, string>> messages = new List<Dictionary<string, string>>();

  using (var reader = command.ExecuteReader())
  {
    while (reader.Read())
    {
      Dictionary<string, string> dict = new Dictionary<string, string>();

      for (int i = 0; i < reader.FieldCount; i++)
      {

        string key = reader.GetName(i);

        string value = reader.IsDBNull(i) ? "" : reader.GetString(i);
        if (key.Equals("pub_date"))
        {
          var dateTimeOffset = DateTimeOffset.FromUnixTimeSeconds(long.Parse(value));
          value = dateTimeOffset.DateTime.ToString("yyyy-MM-dd HH:mm:ss");
        }
        dict[key] = value;
      }

      messages.Add(dict);
    }
  }

   var messagesWithImages = messages.Select(message => new Dictionary<string, object>
  {
    ["username"] = message["username"],
    ["text"] = message["text"],
    ["pub_date"] = message["pub_date"],
    ["email"] = message["email"],
    ["image_url"] = GetGravatarUrl(message["username"].ToString()) // Add a generated image URL
  }).ToList();


  var data = new Dictionary<string, object>
  {
    ["title"] = "Public timeline",
    ["messages"] = messagesWithImages,
    ["endpoint"] = request.Path,
    ["profile_user"] = new Dictionary<string, string>
    {
      ["user_id"] = "2",
      ["username"] = "AnotherUser"
    },
    ["followed"] = false
  };


  var finalRenderedHTML = sendToHtml(data, "timeline");

  return Results.Content(finalRenderedHTML, "text/html; charset=utf-8");
}
// For cases, where we need the username from the url, use the example below.
app.MapGet("/{username}", (string username, HttpRequest request, HttpContext context) =>
  user_timeline(username, request, context)
);

IResult user_timeline(string username, HttpRequest request, HttpContext context)
{

  var db = new SqliteConnection("Data source=" + DATABASE);
  db.Open();


  var command = db.CreateCommand();


  var query = @"select * from user where username = @userName";

  command.CommandText = query;
  command.Parameters.Add(new SqliteParameter("@userName", username));

  List<Dictionary<string, string>> messages = new List<Dictionary<string, string>>();

  using (var reader = command.ExecuteReader())
  {
    while (reader.Read())
    {
      Dictionary<string, string> dict = new Dictionary<string, string>();

      for (int i = 0; i < reader.FieldCount; i++)
      {

        string key = reader.GetName(i);

        string value = reader.IsDBNull(i) ? "" : reader.GetString(i);
        if (key.Equals("pub_date"))
        {
          var dateTimeOffset = DateTimeOffset.FromUnixTimeSeconds(long.Parse(value));
          value = dateTimeOffset.DateTime.ToString("yyyy-MM-dd HH:mm:ss");
        }
        dict[key] = value;
      }

      messages.Add(dict);
    }
  }


  if (messages.Count == 0)
  {
    return Results.NotFound();
  }

  for (int i = 0; i < messages.Count; i++)
  {
    foreach (KeyValuePair<string, string> k in messages[i])
    {
      Console.WriteLine(k);
    }
  }

  var userIDFromSession = "1";
  var profile_user = messages[0]; // profile user

  var queryFollowed = @"SELECT 1 FROM follower 
                      WHERE follower.who_id = @currentUserID 
                      AND follower.whom_id = @profileUserID";

  command.CommandText = queryFollowed;
  command.Parameters.Clear();
  command.Parameters.Add(new SqliteParameter("@currentUserID", userIDFromSession));
  command.Parameters.Add(new SqliteParameter("@profileUserID", profile_user["user_id"]));

  var followed = false;

  using (var reader = command.ExecuteReader())
  {
    if (reader.Read())
    {
      followed = true;
    }
  }



  var queryThree = @"select message.*, user.* from message, user where
            user.user_id = message.author_id and user.user_id = @profileUser
            order by message.pub_date desc limit @Perpage";

  command.CommandText = queryThree;
  command.Parameters.Add(new SqliteParameter("@profileUser", profile_user["user_id"]));
  command.Parameters.Add(new SqliteParameter("@Perpage", PER_page));


  List<Dictionary<string, string>> messagesThree = new List<Dictionary<string, string>>();

  using (var reader = command.ExecuteReader())
  {
    while (reader.Read())
    {
      Dictionary<string, string> dict = new Dictionary<string, string>();

      for (int i = 0; i < reader.FieldCount; i++)
      {

        string key = reader.GetName(i);

        string value = reader.IsDBNull(i) ? "" : reader.GetString(i);
        if (key.Equals("pub_date"))
        {
          var dateTimeOffset = DateTimeOffset.FromUnixTimeSeconds(long.Parse(value));
          value = dateTimeOffset.DateTime.ToString("yyyy-MM-dd HH:mm:ss");
        }
        dict[key] = value;
      }

      messagesThree.Add(dict);
    }
  }

  var messagesWithImages = messagesThree.Select(message => new Dictionary<string, object>
  {
    ["username"] = message["username"],
    ["text"] = message["text"],
    ["pub_date"] = message["pub_date"],
    ["email"] = message["email"],
    ["image_url"] = GetGravatarUrl(message["username"].ToString()) // Add a generated image URL
  }).ToList();


  var data = new Dictionary<string, object>
  {
    ["title"] = $"{profile_user["username"]}'s Timeline",
    ["messages"] = messagesWithImages,
    ["endpoint"] = request.Path,
    ["followed"] = followed,
    ["profile_user"] = profile_user,
    
  };

  string finalRenderedHTML = sendToHtml(data, "timeline");

  return Results.Content(finalRenderedHTML, "text/html; charset=utf-8");
}

app.MapGet("/register", (HttpRequest request, HttpContext context) =>
    register("GET", request, context));

app.MapPost("/register", (HttpRequest request, HttpContext context) =>
    register("POST", request, context));

IResult register(string method, HttpRequest request, HttpContext context)
{
    if (context.Items["user"] != null)
        return Results.Redirect("/");
    Dictionary<string, object> data = new Dictionary<string, object> {
        { "error", null },
    };
    if (method == "POST")
    {
        data["username"] = (string)request.Form["username"];
        data["email"] = (string)request.Form["email"];
        if (request.Form["username"] == "")
            data["error"] = "You have to enter a username";
        else if (request.Form["email"] == "" || !((string)request.Form["email"]).Contains('@'))
            data["error"] = "You have to enter a valid email address";
        else if (request.Form["password"] == "")
            data["error"] = "You have to enter a password";
        else if ((string)request.Form["password"] != (string)request.Form["password2"])
            data["error"] = "The two passwords do not match";
        else if (get_user_id(request.Form["username"]) != null)
            data["error"] = "The username is already taken";
        else
        {
            var db = (SqliteConnection)context.Items["db"];
            var command = db.CreateCommand();
            command.CommandText = @"
                insert into user 
                (username, email, pw_hash) values (@username, @email, @pw_hash)
            ";
            command.Parameters.AddWithValue("@username", (string)request.Form["username"]);
            command.Parameters.AddWithValue("@email", (string)request.Form["email"]);
            var pwHasher = new PasswordHasher<string>();
            command.Parameters.AddWithValue("@pw_hash", pwHasher.HashPassword((string)request.Form["username"], (string)request.Form["password"]));
            command.ExecuteScalar();
            return Results.Redirect("/login");
        }

    }

    string render = sendToHtml(data, "register");
    return Results.Content(render, "text/html; charset=utf-8");
}


app.Run();

static string sendToHtml(Dictionary<string, object> data, string filename)
{
  var templatePath = Path.Combine(Directory.GetCurrentDirectory(), "templates");
  var templateContext = new TemplateContext
  {
    TemplateLoader = new MiniTwitTemplateLoader(templatePath)
  };

  var layoutText = File.ReadAllText(Path.Combine(templatePath, "layout.html"), Encoding.UTF8);
  var timelineText = File.ReadAllText(Path.Combine(templatePath, filename + ".html"), Encoding.UTF8);

  var layoutTemplate = Template.Parse(layoutText);
  var timelineTemplate = Template.Parse(timelineText);

  var timelineContent = timelineTemplate.Render(data);

  var finalData = new Dictionary<string, object>(data)
  {
    ["body"] = timelineContent
  };

  var scriptObject = new ScriptObject();
  scriptObject.Import(finalData);
  templateContext.PushGlobal(scriptObject);

  var finalRenderedHTML = layoutTemplate.Render(templateContext);
  return finalRenderedHTML;
}

static string GetGravatarUrl(string email, int size = 48)
{
  var hash = BitConverter.ToString(MD5.HashData(Encoding.UTF8.GetBytes(email.Trim().ToLower())))
      .Replace("-", "").ToLower();
  return $"https://www.gravatar.com/avatar/{hash}?s={size}&d=retro";
}


