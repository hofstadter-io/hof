package veg

envs: [n=string]: { name: n }
envs: {
  veg: {
    description: "official debian container"
    spec: {
        from: "debian:13-slim"
    }
  }
  golang: {
    description: "official golang container"
    spec: {
        from: "golang:1.25-trixie"
    }
  }
  node: {
    description: "official node container"
    spec: {
        from: "node:25-trixie"
    }
  }
}