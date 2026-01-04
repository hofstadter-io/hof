package veg

// import (
//   "github.com/hofstadter-io/hof/lib/env/devex"
// )

_reg: "ghcr.io/hofstadter-io"
_ver: "v0.7.0-alpha.1"

environs: {
  for short in ["dev", "hof"] {
    (short): {
      name: "veg-\(short)"
      description: "our veg-\(short) image for self & agents"
      spec: {
        from: "\(_reg)/veg-\(short):\(_ver)"
      }
    }

  }
}

// environs: [n=string]: { name: string | *n }
// environs: {
//   for k, env in devex.veg {
//     (k): {
//       name: env.name
//       description: env.description | "\(env.name) image"
//       spec: {
//         from: "\(_reg)/\(env.name):local"
//       }
//     }
//   }
// }