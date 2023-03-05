---
title: Relations
brief: "between types and how to handle them"

weight: 50
---

To add relations to our types and code, we need to do a few things.

- add relations to the schemas and example
- update type template to add fields and helpers
- update our custom code in the resource handlers


### Updating Schemas

Recall that the core `hof/schema/dm.#Datamodel` has a `#Relation` with two strings, `Reln` and `Type`.
To show you how we can calculate extra fields from existing ones,
we will modify our example `#Datamodel` schema,
which already extends `hof`'s core `#Datamodel`.

Here, we add a `GoType` field for template rendering

{{<codePane title="schema/dm.cue" file="code/first-example/data-layer/content/schema/dm-reln.html">}}


### Updating Design

We of course need to add the relations between our types
to the datamodel we are using for our server.
Note how we do not have to add `GoType`.
It will still be available in our template input though.

We specify that a User has many Todos
and that a Todo belongs to a User.

{{<codePane title="example/dm.cue" file="code/first-example/data-layer/content/examples/dm-reln.html">}}


### Type Template

We now need to implement type relations in our templates.
We show the whole file here as much was added.
Note

- We stored the top-level TYPE in `$T` so we can reference it in nested template scopes.
- some edge cases in the `*With*` helpers are omitted and left to the user.

{{< codePane title="templates/type.go" file="code/first-example/data-layer/templates/type.go" lang="go" >}}


### Custom Code

We change

- `UserReadHandler` to use the new library function `UserReadWithTodos`
- `TodoCreateHandler` to have a username query parameter and assign the user. Note, we would normally determine the User with auth and context.

{{<codeInner lang="go">}}
func UserReadHandler(c echo.Context) (err error) {

	username := c.Param("username")

	// note the changed function here
	u, err := types.UserReadWithTodos(username)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, u)
}

func TodoCreateHandler(c echo.Context) (err error) {

	var t types.Todo
	if err := c.Bind(&t); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	username := c.QueryParam("username")
	if username == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing query param 'username'")
	}

	u, err := types.UserReadWithTodos(username)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	t.Author = u

	if err := types.TodoCreate(&t); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, "OK")
}
{{</codeInner>}}


### Query Param

We also need to add a query parameter to an automatically generated route from a resource from our datamodel.

For now, we added it manually to the handler, as in the code just above.
We leave addition of this query parameter as a thought experiment.

- Can you conditionally add it in the `#DatamodelToResource` helper? (hint, this is probably the complex, but better way, check for `Reln == "BelongsTo"`)
- Can you unconditionally add it through the example design? (hint, you will need to unify with the generator `In`)


### Regen, Rebuild, Rerun

We can now `hof gen ./examples/server.cue`.
If you see any merge conflicts, you will need to fix those first.

Try rebuilding, rerunning, and making calls to the server to create users and todos.

{{<codeInner lang="sh">}}
// make a user and a todo
curl -H 'Content-Type: application/json' -X POST -d '{ "Username": "bob", "Email": "bob@hof.io" }' localhost:8080/user
curl -H 'Content-Type: application/json' -X POST -d '{ "Title": "hello", "Content": "world" }' localhost:8080/todo?username=bob

// read a user and a todo
curl localhost:8080/user/bob
curl localhost:8080/todo/hello
{{</codeInner>}}


### Commentary

Some comments on the difficulties of data modeling, templating, and the spectrum of languages and technologies.

- type and relation generation is going to language and technology specific
- hof aims to maintain an abstract representation
- users can extend to add the specifics needed
- we expect that middle layers will capture common groups like SQL, OOP, and visibility (public/private)
- best practice is to use labels which are not likely to conflict, can always namespace and nest

[Link to a more in depth discussion]
of the issue with creating a schema for data models
and the complexities with multiple

We'll also see better ways to construct the library later,
when we introduce `Views`.
