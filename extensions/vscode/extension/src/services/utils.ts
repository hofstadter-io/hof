import * as vscode from 'vscode';

export const SERVER_PORT = 2257;
export const SERVER_URL = `http://localhost:${SERVER_PORT}`;
export const SERVER_REGISTRY_HOST = "host.docker.internal:5000";

export type Environ = {
	name?: string
	srcUri?: string
	srcPath?: string
	fromUri?: string
	dstPath?: string
	workdir?: string
}

export type Folder = {
	uri: vscode.Uri
	sid: string
	name?: string | undefined
	base?: string | undefined
	session?: any
	environ?: Environ
}

export function vsUriToVeg(uri: vscode.Uri): vscode.Uri {
	// console.log("vsUriToVeg.input:", uri.toString())
	let authority = uri.authority
	let p = uri.path
	if (p.startsWith("/")) p = p.slice(1)
	const parts = p.split("/")

	let envSegments: string[] = []
	let pathSegments: string[] = []

	if (authority.includes(":")) {
		const full = authority + (uri.path.startsWith("/") ? uri.path : "/" + uri.path)
		const segments = full.split("/")
		let lastColonIndex = -1
		for (let i = 0; i < segments.length; i++) {
			if (segments[i].includes(":")) lastColonIndex = i
		}
		const envEndIndex = Math.max(0, lastColonIndex)
		envSegments = segments.slice(0, envEndIndex + 1)
		pathSegments = segments.slice(envEndIndex + 1)
	} else {
		envSegments = [authority]
		if (parts[0] !== "") {
			envSegments.push(parts[0])
		}
		pathSegments = parts.slice(1)
	}

	const ociAuthority = envSegments[0]
	const ociPath = "/" + envSegments.slice(1).join("/")
	let filePath = pathSegments.join("/")
	if (filePath === "" || filePath === "/") {
		filePath = "./"
	} else if (!filePath.startsWith("./")) {
		filePath = "./" + (filePath.startsWith("/") ? filePath.slice(1) : filePath)
	}

	const q = new URLSearchParams(uri.query)
	q.set("path", filePath)

	const vUri = {
		scheme: "oci",
		authority: ociAuthority,
		path: ociPath,
		query: q.toString(),
	}
	// console.log("vsUriToVeg.output:", JSON.stringify(vUri))
	try {
		const r = vscode.Uri.from(vUri)
		// console.log("convert.return", r)
		return r
	} catch (e: any) {
		console.error("catch!", e)
		throw e;
	}
}

export function parseEnvUri(uri: vscode.Uri) {
	let full = uri.authority + uri.path
	if (full.endsWith("/")) full = full.slice(0, -1)
	
	const firstSlash = full.indexOf("/")
	const lastColon = full.lastIndexOf(":")
	
	let envId = full
	let envVer = "?"
	
	if (lastColon !== -1) {
		// tag is after last colon IF that colon is after the first slash (host/img:tag)
		// OR if there is no slash at all (img:tag)
		if (firstSlash === -1 || lastColon > firstSlash) {
			envId = full.substring(0, lastColon)
			envVer = full.substring(lastColon + 1)
		}
	}
	
	return { envId, envVer, fullPath: full }
}

export function normalizeEnvId(envId: string): string {
	const parts = envId.split("/")
	if (parts.length > 1) {
		// if parts[0] looks like a host (contains dot or colon)
		if (parts[0].includes(":") || parts[0].includes(".")) {
			return parts.slice(1).join("/")
		}
	}
	return envId
}

export function findSession(sessions: any[], envId: string, uri?: vscode.Uri) {
	// only find sessions for images in our internal registry
	if (uri && uri.authority !== SERVER_REGISTRY_HOST) {
		return undefined
	}

	const normId = normalizeEnvId(envId)
	return sessions.find(s => {
		const sEnv = s.state?.currEnv
		if (!sEnv) { return false }
		
		const { envId: sId } = parseEnvUri(vscode.Uri.parse("oci://" + sEnv))
		return normalizeEnvId(sId) === normId
	})
}

export async function makeReq(route: string, uri?: vscode.Uri, diffUri?: vscode.Uri, body?: any, onlyDiff: boolean = true, sid?: string): Promise<Response> {
	// console.log("filesys.makeReq", route, uri, body)
	const url = `${SERVER_URL}${route}`
	var req = body
	if (uri && !body) {
		const ruri = vsUriToVeg(uri)
		req = {
			// idiots at microsoft encode query params, meaning = sign is %'d and
			// ACTUALLY compliant implementations of URLs ignore it... M$ idiots...
			uri: `${ruri.scheme}://${ruri.authority}${ruri.path}?${ruri.query}`,
			diff: onlyDiff,
		}
		if (!!diffUri) {
			const duri = vsUriToVeg(diffUri)
			req.diffUri = `${duri.scheme}://${duri.authority}${duri.path}?${duri.query}`;
		}
		if (sid) {
			req.sid = sid
		}
	}

	// console.log("filesys.makeReq", route, req)
	return fetch(url, {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
			"X-Veg-User": "tony", // TODO: make this configurable or dynamic
		},
		body: JSON.stringify(req)
	})
}
