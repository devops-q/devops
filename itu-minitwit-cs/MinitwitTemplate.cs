using Scriban;
using Scriban.Parsing;
using Scriban.Runtime;

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
}
