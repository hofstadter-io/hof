func setupRouter(e *echo.Echo) error {

	// Internal routes and Prometheus

	// Static content
	e.Static("/", "client")

	// API routes
	g := e.Group("/api")

	// routes and resources

	return nil
}
