package api

Call: {
	method: *"GET" | "POST" | "PUT" | "DELETE" | "OPTIONS" | "HEAD" | "CONNECT" | "TRACE" | "PATCH"
	host:   string
	path:   string | *""
	auth?:  string
	headers?: [string]: string
	query?: [string]:   string
	form?: [string]:   string
	data?:    string | {...}
	timeout?: string
	retry?: {
		count: int
		timer: string
		codes: [...int]
	}
}

Response: {
	status?: int
	headers?: [string]: string
	body?: string
	value: _
}
