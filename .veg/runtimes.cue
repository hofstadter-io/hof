package veg

// https://github.com/qdm12/basedevcontainer
// https://github.com/qdm12/binpot

runenv: [n=string]: { name: n }
runenv: {
  base: {
    description: "a base container with many common tools"
    spec: from: "qmcgaw/basedevcontainer:debian"
  }

  golang: {
    description: "a golang specific container"
    spec: from: "qmcgaw/godevcontainer:debian"
  }
}