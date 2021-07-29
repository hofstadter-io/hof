package api

import (
	"github.com/hofstadter-io/hof/schema/test"
)

baseReq: test.#HttpRequest & {
	host: "https://postman-echo.com"
}

baseGet: baseReq & {
	method: "GET"
	path: "/get"
}

basePost: baseReq & {
	method: "POST"
	path: "/post"
}

basic: _ @test(suite,api,basic)
basic: {

	get: _ @test(api,basic,get)
	get: {
		req: baseGet & {
			query: {
				cow: "moo"
			}
		}
		resp: {
			status: 200
			body: {
				args: close({
					cow: "moo"
				})
				...
			}
		}
	}

	post: _ @test(api,basic,post)
	post: {
		#args: {
			foo: "bar"
		}
		req: basePost & {
			data: {
				cow: "moo"
			}
		}
		resp: {
			status: 200
			body: {
				args: #args
				json: {
					cow: "moo"
				}
				...
			}
		}
	}

}
