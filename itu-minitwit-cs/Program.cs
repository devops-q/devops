using System.Data.Common;
using System.Net;
using Microsoft.AspNetCore;
using Microsoft.AspNetCore.Components.RenderTree;
using Microsoft.Data.Sqlite;
using Scriban;
string DATABASE = "/tmp/minitwit.db";
int PER_page = 30;
bool DEBUG = true;
string SECRET_KEY = "development key";

var builder = WebApplication.CreateBuilder(args);
var app = builder.Build();




app.Run();

//  timeLine function.
app.MapGet("/", () =>
  timeline(request, context)
);
/// <summary> 
/// Retrieves userID from session; if null or empty return /public (frontpage)
/// if non-null, find messages from db using query, and add these to data dict.
/// finally, return the data onto the timeline.html side 
/// <summary>
IResult timeline(HttpRequest request, HttpContext context)
{
  Console.WriteLine("We got a visitor from: " + request.HttpContext.Connection.RemoteIpAddress?.ToString());
  var userIDFromSession = context.Session.GetString("user_id");
  if (string.IsNullOrEmpty(userIDFromSession))
    return Results.Redirect("/public");
  var query = @"select message.*, user.* from message, user
        where message.author_id = user.user_id and (
            user.user_id = @UserId or
            user.user_id in (select whom_id from follower
                                    where who_id = @UserId))
        order by message.pub_date desc limit @PerPage";


  var twits = db.Query(query, new { UserId = userIDFromSession, PerPage = PER_page });
  var data = new Dictionary<string, object>()
  {
    ["twits"] = twits

  };

  var templateText = File.ReadAllText("templates/timeline.html");
  var template = Template.Parse(templateText);

  var renderedHTML = template.Render(data, memberRenamer: member => member.Name);

  return Results.Content(renderedHTML, "text/html");

}

app.MapGet("/public", () => "");



IResult public_timeline(HttpRequest request, HttpContext context)
{

  var query = @"select message.*, user.* from message, user where message.flagged = 0 
  and message.author_id = user.user_id order by message.pub_date desc limit @PerPage";

  var twits = db.Query(query, new { PerPage = PER_page });
  var data = new Dictionary<string, object>()
  {
    ["twits"] = twits
  };

  var templateText = File.ReadAllText("templates/timeline.html");
  var template = Template.Parse(templateText);

  var renderedHTML = template.Render(data, memberRenamer: member => member.Name);


  return Results.Content(renderedHTML, "text/html");

}


app.MapGet("/<username>", () => "");

IResult user_timeline(HttpRequest request, HttpContext context, string username)
{
  var profile_user = db.QueryFirstOrDefault<string>(@"SELECT * FROM user WHERE username = @Username", new { Username = username });

  if (string.IsNullOrEmpty(profile_user))
  {
    return Results.NotFound();
  }
  var followed = false;
  var userIDFromSession = context.Session.GetString("user_id");

  if (!string.IsNullOrEmpty(userIDFromSession))
  {
    followed = db.QueryFirstOrDefault<int>(@"SELECT 1 FROM follower 
                            WHERE who_id = @UserId AND whom_id = @ProfileUserId", new { UserId = userIDFromSession, ProfileUserId = profile_user.UserId }) != null;
  }

  var query = @"
        SELECT message.*, user.* FROM message, user
        WHERE user.user_id = message.author_id 
        AND user.user_id = @ProfileUserId
        ORDER BY message.pub_date DESC
        LIMIT @PerPage";
  var messages = db.Query(query, new { ProfileUserId = profile_user.UserId, PerPage = PER_page });

  var templateText = File.ReadAllText("templates/timeline.html");

  var template = Template.Parse(templateText);

  var renderedHTML = template.Render(messages, memberRenamer: member => member.Name);

  return Results.Content(renderedHTML, "text/html");
}


