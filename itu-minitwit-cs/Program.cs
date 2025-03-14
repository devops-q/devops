using System.Data.Common;
using System.Runtime.CompilerServices;
using System.Security.Cryptography;
using System.Text;
using System.Text.RegularExpressions;
using System.Threading.Tasks;
using Microsoft.AspNetCore;
using Microsoft.AspNetCore.Http.Extensions;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.ModelBinding.Validation;
using Microsoft.AspNetCore.Identity;
using Microsoft.Data.Sqlite;
using Microsoft.Extensions.Options;
using Scriban;
using Scriban.Runtime;
using System.Web;

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

builder.Services.Configure<DatabaseSettings>(builder.Configuration.GetSection("DatabaseSettings"));

var app = builder.Build();

var databaseSettings = app.Services.GetRequiredService<IOptions<DatabaseSettings>>().Value;
string DATABASE = databaseSettings.Database;

Console.WriteLine(DATABASE);

app.UseSession();
app.UseStaticFiles(); // Enable serving static files like CSS


SqliteConnection ConnectDb()
{
  // Returns a new connection to the database
  var connection = new SqliteConnection("Data Source=" + DATABASE);
  connection.Open();

  return connection;
}

void InitDb()
{
  // Initialize database by executing a schema.sql file against it.
  const string schemaPath = "schema.sql";
  var schemaCommandText = File.ReadAllText(schemaPath);

  using var db = ConnectDb();
  var sqliteCommand = db.CreateCommand();

  sqliteCommand.CommandText = schemaCommandText;
  sqliteCommand.ExecuteReader();
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

long? get_user_id(string username, HttpContext context)
{
  var db = (SqliteConnection)context.Items["db"];
  var command = db.CreateCommand();
  command.CommandText = @"select user_id from user where username = @username";
  command.Parameters.AddWithValue("@username", username);

  return (long?)command.ExecuteScalar();
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
  if (context.Session.GetString("flash") == null)
    context.Session.SetString("flash", "");

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

app.MapGet("/", (HttpRequest request, HttpContext context) =>
    timeline(request, context));

IResult timeline(HttpRequest request, HttpContext context)
{
  var db = (SqliteConnection)context.Items["db"];
  Console.WriteLine("We got a visitor from: " + request.HttpContext.Connection.RemoteIpAddress?.ToString());

  if (string.IsNullOrEmpty(context.Session.GetString("user_id")))
    return Results.Redirect("/public");

  var query = @"
        SELECT message.*, user.* FROM message, user
        WHERE message.author_id = user.user_id AND (
            user.user_id = @p0 OR
            user.user_id IN (SELECT whom_id FROM follower WHERE who_id = @p0))
        ORDER BY message.pub_date DESC LIMIT @p1";

  var messages = QueryDb((SqliteConnection)context.Items["db"], query, [context.Session.GetString("user_id"), PER_page]);

  var user = QueryDb((SqliteConnection)context.Items["db"], "SELECT * FROM user WHERE user_id = @p0",
     [context.Session.GetString("user_id"),], one: true);

  var data = new Dictionary<string, object>
  {
    ["title"] = "My Timeline",
    ["messages"] = messages.Select(message => new Dictionary<string, object>
    {
      ["username"] = message["username"],
      ["text"] = message["text"],
      ["pub_date"] = FormatDatetime(Convert.ToInt32(message["pub_date"])),
      ["email"] = message["email"],
      ["image_url"] = GetGravatarUrl(message["username"].ToString())
    }).ToList(),
    ["endpoint"] = request.Path,
    ["userIDFromSession"] = new Dictionary<string, string>
    {
      ["user_id"] = context.Session.GetString("user_id"),
      ["username"] = user[0]["username"].ToString(),
    },
    ["flashes"] = readFlash(context)
  };



  var render = sendToHtml(data, "timeline");

  return Results.Content(render, "text/html; charset=utf-8");
}

app.MapGet("/public", (HttpRequest request, HttpContext context) =>
  public_timeline(request, context));

IResult public_timeline(HttpRequest request, HttpContext context)
{


  var db = (SqliteConnection)context.Items["db"];



  var query = @"
        SELECT message.*, user.* from message, 
        user where message.flagged = 0 and 
        message.author_id = user.user_id
        order by message.pub_date desc limit @p0";


  var messages = QueryDb(db, query, [PER_page]);
  var data = new Dictionary<string, object>();
  if (context.Session.GetString("user_id") != null)
  {
    var user = QueryDb((SqliteConnection)context.Items["db"], "SELECT * FROM user WHERE user_id = @p0",
        [context.Session.GetString("user_id"),], one: true);
    data = new Dictionary<string, object>
    {
      ["title"] = "Public timeline",
      ["messages"] = messages.Select(message => new Dictionary<string, object>
      {
        ["username"] = message["username"],
        ["text"] = message["text"],
        ["pub_date"] = FormatDatetime(Convert.ToInt32(message["pub_date"])),
        ["email"] = message["email"],
        ["image_url"] = GetGravatarUrl(message["username"].ToString()) // Add a generated image URL
      }).ToList(),
      ["endpoint"] = request.Path,
      ["userIDFromSession"] = new Dictionary<string, string>
      {
        ["user_id"] = context.Session.GetString("user_id"),
        ["username"] = user[0]["username"].ToString(),
      },
      ["flashes"] = readFlash(context)

    };
  }
  else
  {
    data = new Dictionary<string, object>
    {
      ["title"] = "Public timeline",
      ["messages"] = messages.Select(message => new Dictionary<string, object>
      {
        ["username"] = message["username"],
        ["text"] = message["text"],
        ["pub_date"] = FormatDatetime(Convert.ToInt32(message["pub_date"])),
        ["email"] = message["email"],
        ["image_url"] = GetGravatarUrl(message["username"].ToString()) // Add a generated image URL
      }).ToList(),
      ["endpoint"] = request.Path,
      ["flashes"] = readFlash(context)

    };
  }
  ;

  var render = sendToHtml(data, "timeline");

  return Results.Content(render, "text/html; charset=utf-8");
}


app.MapGet("/{username}", (string username, HttpRequest request, HttpContext context) =>
  user_timeline(username, request, context)
);

IResult user_timeline(string username, HttpRequest request, HttpContext context)
{
  var db = (SqliteConnection)context.Items["db"];

  var query = @"SELECT * FROM user WHERE username = @p0";
  var profile_user = QueryDb(db, query, [username], true);

  if (profile_user == null)
    return Results.NotFound();

  var followed = false;
  if (context.Session.GetString("user_id") != null)
  {
    followed = QueryDb(db, @"SELECT 1 FROM follower 
                                 WHERE follower.who_id = @p0 
                                 AND follower.whom_id = @p1",
                             [context.Session.GetString("user_id"), profile_user[0]["user_id"]], true) != null;
  }

  var queryThree = @"SELECT message.*, user.* FROM message, user 
                       WHERE user.user_id = message.author_id AND user.user_id = @p0
                       ORDER BY message.pub_date DESC LIMIT @p1";
  var messages = QueryDb(db, queryThree, [profile_user[0]["user_id"].ToString(), PER_page]);

  var data = new Dictionary<string, object>();

  data = new Dictionary<string, object>
  {
    ["title"] = $"{profile_user[0]["username"]}'s Timeline",
    ["messages"] = messages.Select(message => new Dictionary<string, object>
    {
      ["username"] = message["username"],
      ["text"] = message["text"],
      ["pub_date"] = FormatDatetime(Convert.ToInt32(message["pub_date"])),
      ["email"] = message["email"],
      ["image_url"] = GetGravatarUrl(message["username"].ToString())
    }).ToList(),
    ["endpoint"] = request.Path,
    ["followed"] = followed,
    ["profile_user"] = profile_user[0],
    ["flashes"] = readFlash(context)

  };
  if (context.Session.GetString("user_id") != null)
  {
    data.Add("userIDFromSession", new Dictionary<string, string>
    {
      ["user_id"] = context.Session.GetString("user_id"),
      ["username"] = (string)((List<Dictionary<string, object>>)context.Items["user"])[0]["username"]
    });

  }

  // Ensure the URL is the same without additional layers
  string render = sendToHtml(data, "timeline");

  // Return the content with the correct header
  return Results.Content(render, "text/html; charset=utf-8");
}


app.MapGet("/{username}/follow", follow_user);
IResult follow_user(string username, HttpContext context, HttpRequest request)
{
  if (context.Items["user"] == null)
    return Results.Unauthorized();
  var whomID = get_user_id(username, context);
  if (whomID == null)
    return Results.NotFound();

  var db = (SqliteConnection)context.Items["db"];
  var command = db.CreateCommand();
  command.CommandText = @"insert into follower (who_id, whom_id) values (@whoID, @whomID)";
  command.Parameters.AddWithValue("@whoID", context.Session.GetString("user_id"));
  command.Parameters.AddWithValue("@whomID", whomID);
  command.ExecuteScalar();

  flash($"You are now following &#34;{username}&#34;", context);

  return Results.Redirect($"/{username}");
}

app.MapGet("/{username}/unfollow", unfollow_user);
IResult unfollow_user(string username, HttpContext context, HttpRequest request)
{
  if (context.Items["user"] == null)
    return Results.Unauthorized();
  var whomID = get_user_id(username, context);
  if (whomID == null)
    return Results.NotFound();

  var db = (SqliteConnection)context.Items["db"];
  var command = db.CreateCommand();
  command.CommandText = @"delete from follower where who_id=@whoID and whom_id=@whomID";
  command.Parameters.AddWithValue("@whoID", context.Session.GetString("user_id"));
  command.Parameters.AddWithValue("@whomID", whomID);
  command.ExecuteScalar();
  flash($"You are no longer following &#34;{username}&#34;", context);


  return Results.Redirect($"/{username}");

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
    else if (get_user_id(request.Form["username"], context) != null)
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

      flash("You were successfully registered and can login now", context);

      return Results.Redirect("/login");
    }

  }

  string render = sendToHtml(data, "register");
  return Results.Content(render, "text/html; charset=utf-8");
}

app.MapMethods("/login", new[] { "GET", "POST" }, async (HttpRequest request, HttpContext context) =>
{
  // Logs the user in.

  if (context.Items["user"] != null)
  {
    return Results.Redirect("/");
  }

  string error = null;
  string username = null;

  if (request.Method == "POST")
  {
    var form = await request.ReadFormAsync();
    username = form["username"].ToString();
    var password = form["password"].ToString();

    var user = QueryDb((SqliteConnection)context.Items["db"], "SELECT * FROM user WHERE username = @p0",
      new object[] { username }, one: true)?.FirstOrDefault();

    if (user == null)
    {
      error = "Invalid username";
    }
    else if (!CheckPasswordHash(user["pw_hash"].ToString(), password))
    {
      error = "Invalid password";
    }
    else
    {
      context.Session.SetString("user_id", user["user_id"].ToString());
      flash("You were logged in", context);
      return Results.Redirect("/");
    }
  }

  var data = new Dictionary<string, object>
  {
    ["flashes"] = readFlash(context),
    ["error"] = error,
    ["username"] = username

  };


  var finalRenderedHTML = sendToHtml(data, "login");
  return Results.Content(finalRenderedHTML, "text/html; charset=utf-8");
});

bool CheckPasswordHash(string storedHash, string password)
{
  var pwHasher = new PasswordHasher<string>();

  var verificationResult = pwHasher.VerifyHashedPassword((string)"abc", storedHash, password);

  return verificationResult == PasswordVerificationResult.Success;
}

app.MapGet("/logout", logout);

IResult logout(HttpContext context)
{
  context.Session.Remove("user_id");
  flash("You were logged out", context);

  return Results.Redirect("/");
}
app.MapPost("/add_message", (HttpRequest request, HttpContext context) =>
  add_message(request, context)
);

async Task<IResult> add_message(HttpRequest request, HttpContext context)
{
  var db = context.Items["db"] as SqliteConnection;

  var userIDFromSession = context.Session.GetString("user_id");

  if (string.IsNullOrEmpty(userIDFromSession))
  {
    Console.WriteLine("No id found.");
    return Results.NotFound();
  }

  var form = await request.ReadFormAsync();
  var messageText = HttpUtility.HtmlEncode(form["text"].ToString());

  if (!string.IsNullOrEmpty(messageText))
  {
    var query = @"insert into message (author_id, text, pub_date, flagged) values (@p0, @p1, @p2, 0)";

    QueryDb(db, query, new object[] { userIDFromSession, messageText, DateTimeOffset.UtcNow.ToUnixTimeSeconds() });

    Console.WriteLine("Message added");
  }

  flash("Your message was recorded", context);

  return Results.Redirect("/");
}


void RecreateDatabase()
{
    Console.WriteLine("Recreating database");
    if (File.Exists(DATABASE))
    {
        Console.WriteLine("Deleting old database");
        File.Delete(DATABASE);
    }

    Console.WriteLine("Creating new database");
    InitDb();
}

if (databaseSettings.RecreateDatabase)
{
    RecreateDatabase();
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

static void flash(string message, HttpContext context)
{
    context.Session.SetString("flash", message);
}

static string readFlash(HttpContext context)
{
    var message = context.Session.GetString("flash");
    context.Session.Remove("flash");
    return message;
}

