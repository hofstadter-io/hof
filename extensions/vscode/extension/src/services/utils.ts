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
	var p = uri.path
	if (p.startsWith("/")) {
		p = p.slice(1)
	}
	const parts = p.split("/")
	const e = parts[0]
	const path = parts.slice(1).join("/")
	const q = new URLSearchParams(uri.query)
	if (path !== "") {
		q.set("path", path)
	}
	// console.log("convert.parts", e, q, parts)
	const vUri = {
		scheme: "oci",
		authority: uri.authority,
		path: "/" + e,
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
