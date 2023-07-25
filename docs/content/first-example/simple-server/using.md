---
title: "Using"
brief: "creating with our generator"
weight: 25
---

{{<lead>}}
Users _design_ by filling in the schema for a __Hof Generator__.
{{</lead>}}

### Server, @gen(server)

This is the block which defines an entrypoint to `hof gen`,
using `@gen(tags...)` and unifying with our generator `gen.Generator`.

{{<codePane title="examples/gen.cue" file="code/first-example/simple-server/examples/gen.html">}}

### ServerDesign, with the schema

This is the user created design for our server generator.

{{<codePane title="examples/server.cue" file="code/first-example/simple-server/examples/server.html">}}

### Generate the Code

From the root of our module, run

{{<codeInner lang="sh">}}
hof gen ./examples
{{</codeInner>}}

You should now have an `./output` directory with the generated code.

### Running the Server

With our code in place, we can build and run the server

{{<codeInner lang="sh">}}
cd ./output
go mod tidy
go build -o server
./server
{{</codeInner>}}

Call the endpoints with curl

{{<codeInner lang="sh">}}
curl localhost:8080/hello
curl localhost:8080/echo/moo
curl localhost:8080/internal/metrics
{{</codeInner>}}
