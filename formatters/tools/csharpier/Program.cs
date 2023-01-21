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

app.MapPost("/", async (HttpRequest request) => 
{
	var input = await request.ReadFromJsonAsync<Input>();
	var formatted = await CodeFormatter.FormatAsync(input.Source);

	return formatted;
});

app.Run($"http://*:{port}");

record Input(string Source);
