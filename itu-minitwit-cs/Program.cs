using System.Data.Common;
using System.Security.Cryptography;
using System.Text;
using System.Text.RegularExpressions;
using Microsoft.AspNetCore;
using Microsoft.AspNetCore.Http.Extensions;
using Microsoft.AspNetCore.Mvc.ModelBinding.Validation;
using Microsoft.Data.Sqlite;
using Microsoft.Extensions.FileProviders;
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


  // Data dictionary for template
  var data = new Dictionary<string, object>
  {
    ["title"] = "My Timeline",
    ["messages"] = messages,
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

  var templatePath = Path.Combine(Directory.GetCurrentDirectory(), "templates");
  var templateContext = new TemplateContext
  {
    TemplateLoader = new MiniTwitTemplateLoader(templatePath)
  };

  var layoutText = File.ReadAllText(Path.Combine(templatePath, "layout.html"), Encoding.UTF8);
  var timelineText = File.ReadAllText(Path.Combine(templatePath, "timeline.html"), Encoding.UTF8);

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

  var data = new Dictionary<string, object>
  {
    ["title"] = "Public timeline",
    ["messages"] = messages,
    ["endpoint"] = request.Path,
    ["profile_user"] = new Dictionary<string, string>
    {
      ["user_id"] = "2",
      ["username"] = "AnotherUser"
    },
    ["followed"] = false
  };

  var templatePath = Path.Combine(Directory.GetCurrentDirectory(), "templates");
  var templateContext = new TemplateContext
  {
    TemplateLoader = new MiniTwitTemplateLoader(templatePath)
  };

  var layoutText = File.ReadAllText(Path.Combine(templatePath, "layout.html"), Encoding.UTF8);
  var timelineText = File.ReadAllText(Path.Combine(templatePath, "timeline.html"), Encoding.UTF8);

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




  var userIDFromSession = "1";  
  var profile_user = messages[0]["username"]; 

  var queryFollowed = @"SELECT 1 FROM follower 
                      WHERE follower.who_id = @currentUserID 
                      AND follower.whom_id = @profileUserID";

  command.CommandText = queryFollowed;
  command.Parameters.Clear();
  command.Parameters.Add(new SqliteParameter("@currentUserID", userIDFromSession));
  command.Parameters.Add(new SqliteParameter("@profileUserID", profile_user));

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
  command.Parameters.Add(new SqliteParameter("@profileUser", profile_user));
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


  var data = new Dictionary<string, object>
  {
    ["title"] = $"{profile_user}'s Timeline",
    ["messages"] = messagesThree,
    ["endpoint"] = request.Path,
    ["followed"] = followed,
    ["profile_user"] = profile_user
  };

  var templatePath = Path.Combine(Directory.GetCurrentDirectory(), "templates");
  var templateContext = new TemplateContext
  {
    TemplateLoader = new MiniTwitTemplateLoader(templatePath)
  };

  var layoutText = File.ReadAllText(Path.Combine(templatePath, "layout.html"), Encoding.UTF8);
  var timelineText = File.ReadAllText(Path.Combine(templatePath, "timeline.html"), Encoding.UTF8);

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

  return Results.Content(finalRenderedHTML, "text/html; charset=utf-8");
}



app.Run();


