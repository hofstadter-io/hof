---
title: "Custom Code"
brief: "and iterative development with code gen"
weight: 40
---

{{<lead>}}
One of `hof`'s core design choices was to
make it possible to work directly in the output.
{{</lead>}}


We need to implement our Resource handlers,
so we add our custom code in the generated files.
Here ar the `UserCreate` and `UserRead` handlers.

_Normally, A generator author would provide
defaults for the handler body. (expand on this)


{{<codeInner title="in output/resources/User.go" lang="go">}}
func UserCreateHandler(c echo.Context) (err error) {

	// process any path and query params

	// default body
	// c.String(http.StatusNotImplemented, "Not Implemented")
	// (you can comment out generated code)

	var u types.User
	if err := c.Bind(&u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := types.UserCreate(&u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, "OK")
}

func UserReadHandler(c echo.Context) (err error) {

	// (or you can delete generated code)
	
	username := c.Param("username")

	u, err := types.UserRead(username)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, u)
}
{{</codeInner>}}

We leave the other handlers as an exercise.


### Using the Resources

Let's create and read a user with curl.

{{<codeInner lang="sh">}}
curl -X POST \
  -H 'Content-Type: application/json' \
  -d '{ "Username": "bob", "Email": "bob@hof.io" }' \
  localhost:8080/user

curl localhost:8080/user/bob
{{</codeInner>}}


### Iterative development

`hof` enables you to both update your designs and edit the output,
while still being able to regenerate the code base and continue customizing.

In the next section we will add relations to our types and data model.
This will require changes to schemas, generators, templates, and custom code.
Since we already have a running server with custom code,
this will be our second iteration with `hof` in the develop cycle.

