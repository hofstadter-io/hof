import fs from 'node:fs/promises'
import * as vscode from 'vscode';
import { extensionEmitter, sendMessage } from '../comms';

const SERVER_PORT = 2257;
const SERVER_URL = `http://localhost:${SERVER_PORT}`;

type Folder = {
	readonly uri: vscode.Uri
	readonly sid: string
	readonly name?: string | undefined
	readonly base?: string | undefined
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

	open(value: string) {
		if (!value || value === "") {
			return
		}

		const f = async () => {
			// console.log("filesys.open", value)

			// hacky parsing of user input to something our extension server understands
			if (value?.startsWith("file://")) {
				// no-op
			} else if (value?.startsWith("https://")) {
				// no-op
			} else if (value?.startsWith("/")) {
				value = `file://${value}`
			} else if (value?.startsWith("oci://")) {
				// no-op
			} else if (value?.includes(":")) {
				// assume oci
				// look for no-domain, use docker.io (i.e. implied in docker pull)
				// the backend requires fully qualified oci://domain.com/reg/org/img:tag
				const parts = value.split(":")
				const paths = parts[0].split("/")
				const host = paths[0]
				if (!host.includes(".")) {
					value = `docker.io/${value}`
				}
				value = `oci://${value}`
			} else {
				vscode.window.showErrorMessage(`unsupported environ: ${value}`)
			}

			const resp = await this.makeReq("/fs/open", undefined, {
				fromUri: value,
			})
			// console.log("filesys.resp:", resp)
			if (resp.status !== 200) {
				vscode.window.showErrorMessage(`${resp.status} - ${resp.statusText}`)
				return
			}

			const data: any = await resp.json()
			// console.log("created:", data)

			// convert to something vscode will understand

			const uri = vscode.Uri.parse(`veg://${data.envUri}`)
			const parts = uri.path.split(":")
			var sid = parts[0]
			if (sid.startsWith("/")) {
				sid = sid.substring(1)
			}
			// const tag = parts[1]

			const f: Folder = {
				uri,
				name: value,
				sid: sid,
				// base: dir,
			}

			const wsF = vscode?.workspace?.workspaceFolders as any[]
			const count = vscode.workspace.workspaceFolders?.length || 0

			var ws: vscode.WorkspaceFolder | any = null
			// console.log("openFS:", sid, dir, name)
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