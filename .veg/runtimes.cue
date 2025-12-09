package veg

// https://github.com/qdm12/basedevcontainer
// https://github.com/qdm12/binpot

environs: [n=string]: { name: n }
environs: {
  debian: {
    description: "official debian container"
    spec: from: "debian:13-slim"
  }
  golang: {
    description: "official golang container"
    spec: from: "golang:1.25-trixie"
  }
  node: {
    description: "official node container"
    spec: from: "node:25-trixie"
  }
}