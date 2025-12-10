import fs from 'node:fs/promises'
import * as vscode from 'vscode';
import { extensionEmitter, sendMessage } from '../comms';

const SERVER_PORT = 2257;
const SERVER_URL = `http://localhost:${SERVER_PORT}`;

type Folder = {
	uri: vscode.Uri
	sid: string
	name?: string | undefined
	base?: string | undefined
	session?: any
}
type FolderListing = Array<[string, vscode.FileType]>
// This method is called when your extension is activated
// Your extension is activated the very first time the command is executed
export async function activate(context: vscode.ExtensionContext) {

	const vcp = new VegContentProvider()

	vscode.workspace.registerFileSystemProvider("veg", vcp, {
		isCaseSensitive: true,
		isReadonly: false,
	})

	vscode.commands.registerCommand('veg.explorer.chat', async (uri: any) => {
		console.log("veg.explorer.chat.args", uri)
		const msg = {
			type: "session.create",
			payload: {
				focus: true,
				fromUri: uri.toString(),
			}
		}
		sendMessage(msg)

  })

	vscode.commands.registerCommand('veg.explorer.openEnviron', async (args: any) => {
		console.log("veg.explorer.openEnviron.args", args)
		const value = await vscode.window.showInputBox({
			title: "Veg Open",
			prompt: "Open a directory, git repo, or any image",
			placeHolder: "/path/to/... | github.com/... | image:tag",
		})

		vcp.open(value as string)
  })

	vscode.commands.registerCommand('veg.explorer.forkEnviron', async (args: any) => {
		console.log("veg.explorer.forkEnviron.args", args)
  })

	vscode.commands.registerCommand('veg.explorer.toggleShown', async (args: any) => {
		console.log("veg.explorer.toggleShown.args", args)
    vcp.toggleShown()
		vscode.commands.executeCommand('workbench.files.action.refreshFilesExplorer')
  })

	vscode.commands.registerCommand('veg.explorer.showDiff', async (args: any) => {
		console.log("veg.explorer.showDiff.args", args)
  })

	vscode.commands.registerCommand('veg.explorer.refreshAll', async (args: any) => {
		console.log("veg.explorer.refreshAll.args", args)
		vscode.commands.executeCommand('workbench.files.action.refreshFilesExplorer')
  })


	extensionEmitter.event(async (e) => {
    // ...
		switch (e.type) {
			case "filesys.openEnviron":
				console.log("filesys.openEnviron.payload", e.payload)

				// prefer envUri because it is more specific
				const envUri = e.payload.envUri
				if (envUri) {
					vcp.open(envUri)
					return
				}

				// if session, use curEnv in state
				const session = e.payload.session
				if (session) {
					if (session.state?.currEnv) {
						vcp.open("", session)
					} else {
						console.error("filesys.openEnviron called with invalid params")
					}
					return
				}

				// otherwise session
				console.error("filesys.openEnviron called with invalid params")
				return
		}
	});

}

// This method is called when your extension is deactivated
export async function deactivate() {}

// https://github.com/microsoft/vscode-extension-samples/blob/main/fsprovider-sample/src/fileSystemProvider.ts

// doing this so we can view diffs across agentic writes
// and accept or fork at various points
// todo, we need to persist this somewhere, the database or on the filesystem?
class VegContentProvider implements vscode.FileSystemProvider {

	private _environs: Record<string,any> = {}

	// we track a show (only) diff or everything
  private _showDiff: boolean = true
  toggleShown() {
    this._showDiff = !this._showDiff
		console.log("VEG.dagger.fs.toggleShown", this._showDiff)
  }


	private _emitter = new vscode.EventEmitter<vscode.FileChangeEvent[]>();
	private _bufferedEvents: vscode.FileChangeEvent[] = [];
	private _fireSoonHandle?: NodeJS.Timeout;

	readonly onDidChangeFile: vscode.Event<vscode.FileChangeEvent[]> = this._emitter.event

	//
	// Translator functions, we need more of these. This component is likely the best place to consolidate the path complexity we have throughout the integration with vscode. Our extension/agent server needs to be simple so that it can interface with many systems.
	//

	private vsUriToVeg(uri: vscode.Uri): vscode.Uri {
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
		} catch(e: any) {
			console.error("catch!", e)
			throw e;
		}
	}

	private async makeReq(route: string, uri?: vscode.Uri, body?: any): Promise<Response> {
		// console.log("filesys.makeReq", route, uri, body)
		const url = `${SERVER_URL}${route}`
		var req = body
		if (uri && !body) {
			const ruri = this.vsUriToVeg(uri)
			req = {
				// fucking idiots at microsoft encode query params, meaning = sign is %'d and
				// ACTUALLY compliant implementations of URLs ignore it... FUCKING M$ IDIOTS!
				uri: `${ruri.scheme}://${ruri.authority}${ruri.path}?${ruri.query}`,
				diff: this._showDiff,
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

	// maybe path is not needed here, it is part of the uri, or we pass it around separately? there are places where we only have a single string to work with, which is why we started stuffing things into a URI, which is pretty flexible tbh
	// open(anyUri: string, path?: string, session?: any) {
	open(inputUri: string, session?: any) {
		console.log("filesys.open", inputUri, session)
		if ((!inputUri || inputUri === "") && !session) {
			return
		}

		// we need to do some logic here to clean things up
		// Uri parsing may be insufficient, it barfs on certain inputs
		// const anyUri = vscode.Uri.parse(inputUri)
		// So _i think_ this was switched to a more manual processing, though we could probably do some url parsing down under the hood?
		var anyUri = inputUri
		var path = ""

		const f = async () => {
			var sid = ""
			var envUri = ""
			var name = anyUri

			// hacky parsing of user input to something our extension server understands
			if (anyUri?.startsWith("veg://")) {
				// hmmm, what do we do here
				// do we use a veg://session.env/<sid>/... or does the sid come earlier?
				// we want a way to always display the latest as it gets updated
				// and thus we also need a way to indicate to refresh when as session gets updated
			} else if (anyUri?.startsWith("file://")) {
				// no-op
				name = anyUri.substring(7)
			} else if (anyUri?.startsWith("https://")) {
				// no-op
				name = anyUri.substring(8)
			} else if (anyUri?.startsWith("/")) {
				anyUri = `file://${anyUri}`
			} else if (anyUri?.startsWith("oci://")) {
				// no-op
				name = anyUri.substring(6)
			} else if (anyUri?.includes(":")) {
				// assume oci
				// look for no-domain, use docker.io (i.e. implied in docker pull)
				// the backend requires fully qualified oci://domain.com/reg/org/img:tag
				const parts = anyUri.split(":")
				const paths = parts[0].split("/")
				const host = paths[0]
				if (!host.includes(".")) {
					anyUri = `docker.io/${anyUri}`
				}
				anyUri = `oci://${anyUri}`
			} else {
				if (session?.state?.currEnv) {
					// these will always be oci
					sid = session.sid
					name = session.state.title || sid
					// TODO, this should be a veg://<session>... something? we need to sort this out eventually, translators, more fields so we can differentiate, more alignment with server in open(...args) too?
					// seems we can add extra without borking things up?
					// this should probably be handled on the server during state managemtn
					envUri = "oci://" + session.state.currEnv
				} else {
					vscode.window.showErrorMessage(`unsupported environ: ${anyUri}`)
				}
			}

			if (envUri === "") {
				envUri = anyUri
			}

			const tmpUri = vscode.Uri.parse(envUri)
			const uri = vscode.Uri.from({
				...tmpUri,
				scheme: "veg",
			})

			const f: Folder = {
				uri,
				name,
				sid,
				session,
				// base: dir,
			}

			console.log("filesys.open.midway", anyUri, uri, f, path, sid, envUri)

			const resp = await this.makeReq("/fs/open", undefined, {
				fromUri: envUri,
				path,
			})
			// console.log("filesys.resp:", resp)
			if (resp.status !== 200) {
				vscode.window.showErrorMessage(`${resp.status} - ${resp.statusText}`)
				return
			}

			const data: any = await resp.json()
			console.log("filesys.open.api.resp:", data)

			const vegUri = vscode.Uri.parse(`veg://${data.envUri}`)
			f.uri = vegUri
			// const parts = uri.path.split(":")
			// var sid = parts[0]
			// if (sid.startsWith("/")) {
			// 	sid = sid.substring(1)
			// }
			// const tag = parts[1]

			// convert to something vscode will understand

			const wsF = vscode?.workspace?.workspaceFolders as any[]
			const count = vscode.workspace.workspaceFolders?.length || 0

			var ws: vscode.WorkspaceFolder | any = null
			console.log("filesys.open.wsFolders:", f, wsF)
			for (var i in wsF){
				// console.log("openFS.loop:", i, wsF[i].sid === sid, wsF[i])
				const qp = new URLSearchParams(wsF[i].uri.query)
				var s = qp.get("sid") as string
				if (s === sid) {
					// HMMM, I doubt this is right
					return
				}

			}
			// console.log("filesys.open calling", f, count, uri)

			const started = vscode.workspace.updateWorkspaceFolders(count, null, f)
			// console.log("filesys.open started?", started)
		}
		return f()
	}

	stat(uri: vscode.Uri): vscode.FileStat | Thenable<vscode.FileStat> {
		// console.log("fs.stat.uri", uri)
		const f = async () => {
			const resp = await this.makeReq("/fs/stat", uri)
			// console.log("fs.stat.resp", resp)
			if (resp.status !== 200) {
				throw vscode.FileSystemError.FileNotFound(uri)
			}

			const data: any = await resp.json()
			// console.log("fs.stat.data", uri, data)
			return {
				ctime: data.ctime,
				mtime: data.mtime,
				size: data.size,
				type: data.dir ? vscode.FileType.Directory : vscode.FileType.File,
			}

		}
		return f()
	}

	readFile(uri: vscode.Uri): Uint8Array | Thenable<Uint8Array> {
		const f = async () => {
			// console.log("filesys.readFile.uri", uri)
			const resp = await this.makeReq("/fs/read", uri)
			// console.log("filesys.readFile.resp", uri, resp)
			if (resp.status !== 200) {
				// console.error("readFile.makeReq error:", uri, resp)
				throw vscode.FileSystemError.FileNotFound(uri)
			}

			const data: any = await resp.json()
			// console.log("filesys.readFile.data", uri, data)
			const c: string = data.content
			return new TextEncoder().encode(c)

		}
		return f()
	}

	readDirectory(uri: vscode.Uri): FolderListing | Thenable<FolderListing> {
		const f = async () => {
			// console.log("filesys.readDir.uri", uri)
			const resp = await this.makeReq("/fs/list", uri)
			if (resp.status !== 200) {
				console.error("readDirectory.makeReq error:", resp)
				throw vscode.FileSystemError.FileNotFound(uri)
			}
			// console.log("filesys.readDir.resp", uri, resp)

			const data: any = await resp.json()
			// console.log("filesys.readDir.data", uri, data, data?.entries)

			// our returned listing
			var l: FolderListing = []
			for (const entry of data?.entries) {
				// console.log("filesys.readDir.entry", entry)
				const path = entry.name
				const isDir = entry.dir
				l.push([path, isDir ? vscode.FileType.Directory : vscode.FileType.File])
			}
			// console.log("filesys.readDir.list", uri, l)

			return l
		}
		return f()
	}

	writeFile(uri: vscode.Uri, content: Uint8Array, options: {create: boolean, overwrite: boolean}): void | Thenable<void> {
	}	


	createDirectory(uri: vscode.Uri): void | Thenable<void> {

	}

	delete(uri: vscode.Uri, options: { recursive: boolean }): void | Thenable<void> {

	}

	rename(source: vscode.Uri, destination: vscode.Uri, options: { overwrite: boolean }): void | Thenable<void> {
		// move
	}

	copy(source: vscode.Uri, destination: vscode.Uri, options: { overwrite: boolean }): void | Thenable<void> {

	}

	watch(uri: vscode.Uri, options: {excludes: readonly string[], recursive: boolean}): vscode.Disposable {
		const handler = () => {
      // DO NOT IMPLEMENT YET
		}
		return new vscode.Disposable(handler)
	}

	// returns our fs key from the uri (<session>[-<pos>])
	private uriToEid(uri: vscode.Uri): string {
		const qp = new URLSearchParams(uri.query)
		var e = qp.get("eid") as string
		if (!e) {
			e = qp.get("sid") as string
		}
		if (!e) {
			console.error("eid unset in vscode.Uri", uri, qp)
			return "error"
		}

		// const p = qp.get("pos") as string
		// if (p && p !== '') {
		// 	s += "-" + p
		// }
		return e
	}

}