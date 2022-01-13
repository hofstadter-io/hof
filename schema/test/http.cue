package test

#HttpTester: {
	req:   #HttpRequest
	resp?: #HttpResponse
}

#HttpRequest: {
	method: *"GET" | "POST" | "PUT" | "DELETE" | "OPTIONS" | "HEAD" | "CONNECT" | "TRACE" | "PATCH"
	host:   string
	path:   string | *""
	auth?:  string
	headers?: [string]: string
	query?: [string]:   string
	data?:    string | {...}
	timeout?: string
	retry?: {
		count?: int
		timer?: string
		codes: [...int]
	}
}

#HttpResponse: {
	status?: int
	headers?: [string]: string
	body?: string
	// latency?: float
}
