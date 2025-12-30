import config from 'dotenv'
import { IdResolver, HandleResolver } from '@atproto/identity'
import { AtpAgent } from '@atproto/api'

console.log("Login Example")

// Load some .env values
config.config({
  path: './examples/.env'
})
const plcDomain = process.env.PLC
const pdsDomain = process.env.PDS
const handle = process.env.HANDLE
const password = process.env.PASSWORD

console.log("PLC:    ", plcDomain)
console.log("PDS:    ", pdsDomain)
console.log("Handle: ", handle)

// Identity resolver for testnet
const idr = new IdResolver({
  plcUrl: `https://${plcDomain}`
})

// Handle -> DID
const did = await idr.handle.resolve(handle)
console.log("DID:    ", did)

// DID -> Doc (PDS)
const doc = await idr.did.resolve(did)
// console.log("DOC:", JSON.stringify(doc, null, "  "))
const pds = doc["service"]?.filter(s => s.id === "#atproto_pds")[0].serviceEndpoint
console.log("PDS:    ", pds)

// create an agent
const agent = new AtpAgent({ service: pds as string })

const repo = await agent.com.atproto.repo.describeRepo({
  repo: did as string
})

console.log("Repo:", JSON.stringify(repo.data, null, "  "))

// login user to PDS
const session = await agent.login({
  identifier: handle as string,
  password: password as string,
})

console.log("Session:", JSON.stringify(session.data, null, "  "))
