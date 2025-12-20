import * as vscode from 'vscode';

export const SERVER_PORT = 2257;
export const SERVER_URL = `http://localhost:${SERVER_PORT}`;

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
	// console.log("convert.uri", uri)
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
	const filePath = pathSegments.join("/")

	const q = new URLSearchParams(uri.query)
	if (filePath !== "") {
		q.set("path", filePath)
	}

	const vUri = {
		scheme: "oci",
		authority: ociAuthority,
		path: ociPath,
		query: q.toString(),
	}
	// console.log("convert.vUri", vUri)
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
	let authority = uri.authority
	let p = uri.path
	if (p.startsWith("/")) p = p.slice(1)
	const parts = p.split("/")

	let envStr = ""
	if (authority.includes(":")) {
		const full = authority + (uri.path.startsWith("/") ? uri.path : "/" + uri.path)
		const segments = full.split("/")
		let lastColonIndex = -1
		for (let i = 0; i < segments.length; i++) {
			if (segments[i].includes(":")) lastColonIndex = i
		}
		const envEndIndex = Math.max(0, lastColonIndex)
		envStr = segments.slice(0, envEndIndex + 1).join("/")
	} else {
		envStr = authority
		if (parts[0] !== "" && parts[0] !== undefined) {
			envStr += (envStr !== "" ? "/" : "") + parts[0]
		}
	}

	const lastColon = envStr.lastIndexOf(":")
	
	let envId = envStr
	let envVer = "?"

	if (lastColon !== -1) {
		envId = envStr.substring(0, lastColon)
		envVer = envStr.substring(lastColon + 1)
	}

	return { envId, envVer, fullPath: uri.authority + uri.path }
}

export function findSession(sessions: any[], envId: string) {
	return sessions.find(s => {
		const sEnv = s.state?.currEnv
		if (!sEnv) { return false }
		const lastColon = sEnv.lastIndexOf(":")
		const sId = lastColon !== -1 ? sEnv.substring(0, lastColon) : sEnv
		return sId === envId
	})
}

export async function makeReq(route: string, uri?: vscode.Uri, diffUri?: vscode.Uri, body?: any, onlyDiff: boolean = true): Promise<Response> {
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
	}

	// console.log("filesys.makeReq", route, req)
	return fetch(url, {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
		},
		body: JSON.stringify(req)
	})
}
