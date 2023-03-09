using System.Text;
using CSharpier;

var builder = WebApplication.CreateBuilder(args);
var port = Environment.GetEnvironmentVariable("PORT") ?? "3000";

var app = builder.Build();

app.UseExceptionHandler(exceptionHandlerApp 
    => exceptionHandlerApp.Run(async context 
        => await Results.Problem()
                     .ExecuteAsync(context)));

app.MapGet("/", () => {
	app.Logger.LogInformation("Hello World!");
	return "CSharpier settings tbd...";
});

app.MapPost("/", async (HttpContext context, HttpRequest request) => 
{
	var input = await request.ReadFromJsonAsync<Input>();
	var result = await CodeFormatter.FormatAsync(input.Source);

	if (result.CompilationErrors.Any()) {
		var msg = new StringBuilder();
		msg.AppendLine("Errors while formatting:");
		foreach (var err in result.CompilationErrors) {
			msg.AppendLine(err.ToString());
		}
		return Results.BadRequest(msg.ToString());
	}

	return Results.Text(result.Code);
});

app.Run($"http://*:{port}");

record Input(string Source);
