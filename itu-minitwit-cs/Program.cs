using System.Data.Common;
using System.Security.Cryptography;
using System.Text;
using Microsoft.AspNetCore;
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

    // Simulate a session user ID
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

                dict[key] = value;
            }

            messages.Add(dict);
        }
    }

    // Convert Unix timestamps to human-readable dates
    foreach (var message in messages)
    {
        var timestamp = long.Parse(message["pub_date"]);
        var dateTimeOffset = DateTimeOffset.FromUnixTimeSeconds(timestamp);
        message["pub_date"] = dateTimeOffset.DateTime.ToString("yyyy-MM-dd HH:mm:ss");
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

    // Template file path
    var templatePath = Path.Combine(Directory.GetCurrentDirectory(), "templates");
    var templateContext = new TemplateContext
    {
        TemplateLoader = new MiniTwitTemplateLoader(templatePath)
    };

    // Load and parse layout and timeline templates
    var layoutText = File.ReadAllText(Path.Combine(templatePath, "layout.html"), Encoding.UTF8);
    var timelineText = File.ReadAllText(Path.Combine(templatePath, "timeline.html"), Encoding.UTF8);

    var layoutTemplate = Template.Parse(layoutText);
    var timelineTemplate = Template.Parse(timelineText);

    // Render timeline content separately
    var timelineContent = timelineTemplate.Render(data);

    // Inject timeline content into the layout
    var finalData = new Dictionary<string, object>(data)
    {
        ["body"] = timelineContent
    };

    var scriptObject = new ScriptObject();
    scriptObject.Import(finalData);
    templateContext.PushGlobal(scriptObject);

    // Render the final layout with timeline inside it
    var finalRenderedHTML = layoutTemplate.Render(templateContext);

    return Results.Content(finalRenderedHTML, "text/html; charset=utf-8");
}

app.Run();

// Custom Scriban Template Loader
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
