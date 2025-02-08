using System.Data.Common;
using System.Net;
using System.Security.Cryptography;
using System.Text;
using Microsoft.AspNetCore;
using Microsoft.Data.Sqlite;
using Microsoft.Extensions.FileProviders;
using Scriban;
using Scriban.Parsing;
using Scriban.Runtime;
string DATABASE = "/tmp/minitwit.db";
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
app.UseStaticFiles();  // This allows static file serving, such as your CSS




//  timeLine function.
app.MapGet("/", (HttpRequest request, HttpContext context) =>
    timeline(request, context));

IResult timeline(HttpRequest request, HttpContext context)
{
  var db = new SqliteConnection("Data source=" + DATABASE);
  db.Open();


  // var userIDFromSession = context.Session.GetString("user_id");

  var userIDFromSession = "1";

  var command = db.CreateCommand();
  Console.WriteLine("We got a visitor from: " + request.HttpContext.Connection.RemoteIpAddress?.ToString());


  if (string.IsNullOrEmpty(userIDFromSession))
    return Results.Redirect("/public");

  var query = @"SELECT message.*, user.* FROM message, user
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

        dict[key] = value;
      }

      messages.Add(dict);
    }
  }




static string dateFormat(long timestamp) {
      DateTimeOffset dateTimeOffset = DateTimeOffset.FromUnixTimeSeconds(timestamp);
    return dateTimeOffset.DateTime.ToString("yyyy-MM-dd HH:mm:ss");  // returns the date in "yyyy-MM-dd HH:mm:ss" format

}

foreach (var message in messages)
{
    var timestamp = long.Parse(message["pub_date"]);
    var dateTimeOffset = DateTimeOffset.FromUnixTimeSeconds(timestamp);
    message["pub_date"] = dateTimeOffset.DateTime.ToString("yyyy-MM-dd HH:mm:ss");
}





// Data dictionary that will include the datetime formatting function
var data = new Dictionary<string, object>()
{
    ["title"] = userIDFromSession,  // This can be a string or any other value you want to display as the title
    ["messages"] = messages,
    ["endpoint"] = request.Path,
    ["userIDFromSession"] = new Dictionary<string, string>
    {
        ["user_id"] = userIDFromSession,
        ["username"] = "TestUser"  // Replace with actual username from session
    },
    // Pass the datetime formatting function to the template context
    // The datetimeformat function will now be available in the template
    ["profile_user"] = new Dictionary<string, string>
    {
        ["user_id"] = "2",
        ["username"] = "AnotherUser"  // Replace with actual profile username
    },
    ["followed"] = false,  // Modify this based on the actual follow status
};

// Now you can call the "datetimeformat" function in your template when rendering the message's timestamp

  // Ensure template file exists before reading
  var templatePath = Path.Combine(Directory.GetCurrentDirectory(), "templates");

  // Create a single template context and assign the template loader
  var templateContext = new TemplateContext();
  templateContext.TemplateLoader = new MiniTwitTemplateLoader(templatePath);



  // Load and parse the main template
var templateText = File.ReadAllText(Path.Combine(templatePath, "timeline.html"), System.Text.Encoding.UTF8);
  var template = Template.Parse(templateText);

  // Import data into Scriban's ScriptObject
  var scriptObject = new ScriptObject();
  scriptObject.Import(data);
  templateContext.PushGlobal(scriptObject); // Use the same context!

  // Render the template
  var renderedHTML = template.Render(templateContext);

return Results.Content(renderedHTML, "text/html; charset=utf-8");
}



// app.MapGet("/public", () => "");



// IResult public_timeline(HttpRequest request, HttpContext context)
// {

//   var query = @"select message.*, user.* from message, user where message.flagged = 0 
//   and message.author_id = user.user_id order by message.pub_date desc limit @PerPage";

//   var twits = db.Query(query, new { PerPage = PER_page });
//   var data = new Dictionary<string, object>()
//   {
//     ["twits"] = twits
//   };

//   var templateText = File.ReadAllText("templates/timeline.html");
//   var template = Template.Parse(templateText);

//   var renderedHTML = template.Render(data, memberRenamer: member => member.Name);


//   return Results.Content(renderedHTML, "text/html");

// }


// app.MapGet("/<username>", () => "");

// IResult user_timeline(HttpRequest request, HttpContext context, string username)
// {
//   var profile_user = db.QueryFirstOrDefault<string>(@"SELECT * FROM user WHERE username = @Username", new { Username = username });

//   if (string.IsNullOrEmpty(profile_user))
//   {
//     return Results.NotFound();
//   }
//   var followed = false;
//   var userIDFromSession = context.Session.GetString("user_id");

//   if (!string.IsNullOrEmpty(userIDFromSession))
//   {
//     followed = db.QueryFirstOrDefault<int>(@"SELECT 1 FROM follower 
//                             WHERE who_id = @UserId AND whom_id = @ProfileUserId", new { UserId = userIDFromSession, ProfileUserId = profile_user.UserId }) != null;
//   }

//   var query = @"
//         SELECT message.*, user.* FROM message, user
//         WHERE user.user_id = message.author_id 
//         AND user.user_id = @ProfileUserId
//         ORDER BY message.pub_date DESC
//         LIMIT @PerPage";
//   var messages = db.Query(query, new { ProfileUserId = profile_user.UserId, PerPage = PER_page });

//   var templateText = File.ReadAllText("templates/timeline.html");

//   var template = Template.Parse(templateText);

//   var renderedHTML = template.Render(messages, memberRenamer: member => member.Name);

//   return Results.Content(renderedHTML, "text/html");
// }


app.Run();




public class MiniTwitTemplateLoader : ITemplateLoader
{
  private readonly string _basePath;

  public MiniTwitTemplateLoader(string basePath)
  {
    _basePath = basePath;
  }

  public string GetPath(TemplateContext context, SourceSpan callerSpan, string templateName)
  {
    return Path.Combine(_basePath, templateName);
  }

  public string Load(TemplateContext context, SourceSpan callerSpan, string templatePath)
  {
    return File.ReadAllText(templatePath);
  }

  public ValueTask<string> LoadAsync(TemplateContext context, SourceSpan callerSpan, string templatePath)
  {
    throw new NotImplementedException();
  }

  public object LoadInteractive(TemplateContext context, SourceSpan callerSpan, string templatePath)
  {
    return Load(context, callerSpan, templatePath);
  }
}


public static class GravatarHelper
{
    // Method to get the Gravatar URL
    public static string GetGravatarUrl(string email, int size = 48)
    {
        // Clean and lower the email address
        email = email.Trim().ToLower();

        // Get the MD5 hash of the email
        string hash = GetMd5Hash(email);

        // Return the Gravatar URL
        return $"https://www.gravatar.com/avatar/{hash}?s={size}";
    }

    // Method to get MD5 hash of the email address
    private static string GetMd5Hash(string input)
    {
        using (MD5 md5 = MD5.Create())
        {
            byte[] inputBytes = Encoding.UTF8.GetBytes(input);
            byte[] hashBytes = md5.ComputeHash(inputBytes);

            // Convert byte array to hex string
            StringBuilder sb = new StringBuilder();
            foreach (byte b in hashBytes)
            {
                sb.Append(b.ToString("x2"));
            }
            return sb.ToString();
        }
    }
}