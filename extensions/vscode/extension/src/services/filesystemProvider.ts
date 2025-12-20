import fs from 'node:fs/promises'
import * as vscode from 'vscode';
import { extensionEmitter, sendMessage } from '../comms';

const SERVER_PORT = 2257;
const SERVER_URL = `http://localhost:${SERVER_PORT}`;

type Environ = {
	name?: string
	srcUri?: string
	srcPath?: string
	fromUri?: string
	dstPath?: string
	workdir?: string
}

type Folder = {
	uri: vscode.Uri
	sid: string
	name?: string | undefined
	base?: string | undefined
	session?: any
	environ?: Environ
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
	vscode.commands.registerCommand('veg.filesys.hack', async (arg1: any) => {
		console.log("veg.filesys.hack", arg1)
	})

	vscode.commands.registerCommand('veg.explorer.chat', async (uri: vscode.Uri) => {
		console.log("veg.explorer.chat.args", uri)

		// we'll have to decipher uri to figure out the new details
		//   ... when we get to forking sessions?
		//   for now, how do we decipher? file:// vs veg:// ?
		// how do we handle subdirectories that have been clicked?

		const msg: any = {
			type: "session.create",
			payload: {
				focus: true,
			}
		}
		switch (uri.scheme) {

			case "file":
				msg.payload.environ = {
					srcUri: uri.toString(),
				}
				break

			case "veg":
				// todo, see if there is a session query param (sid)
				msg.payload.environ = {
					fromUri: uri.toString(),
				}
				break

			default:
				vscode.window.showErrorMessage(`unsupported chat uri: ${uri}`)
				return
		}

		console.log("veg.explorer.chat.msg", msg)
		sendMessage(msg)

	})

	vscode.commands.registerCommand('veg.explorer.openEnviron', async (args: any) => {
		console.log("veg.explorer.openEnviron.args", args)
		const value = await vscode.window.showInputBox({
			title: "Veg Open",
			prompt: "Open a directory, git repo, or any image",
			placeHolder: "/path/to/... | https://github.com/... | image:tag",
		})

		vcp.open(value as string)
	})

	vscode.commands.registerCommand('veg.explorer.openSession', async (session: any) => {
		console.log("veg.explorer.openEnviron.session", session)

		// vcp.open(value as string)
	})

	vscode.commands.registerCommand('veg.explorer.forkEnviron', async (args: any) => {
		console.log("veg.explorer.forkEnviron.args", args)
	})

	vscode.commands.registerCommand('veg.explorer.toggleShown', async (args: any) => {
		console.log("veg.explorer.toggleShown.args", args)
		vcp.toggleShown()
		vscode.commands.executeCommand('workbench.files.action.refreshFilesExplorer')
	})

	vscode.commands.registerCommand('veg.explorer.showDiff', async (uri: vscode.Uri) => {
		console.log("veg.explorer.showDiff.args.CMD.uri", uri)
		vcp.showDiff(uri)
	})

	vscode.commands.registerCommand('veg.explorer.mergeDiff', async (arg: any) => {
		console.log("veg.explorer.mergeDiff.args", arg)
		vcp.mergeDiff(arg, undefined, true)
	})

	vscode.commands.registerCommand('veg.explorer.hideDiff', async (arg: any) => {
		console.log("veg.explorer.hideDiff.args", arg)
		vcp.hideDiff(arg)
	})

	vscode.commands.registerCommand('veg.explorer.copyPath', async (uri: vscode.Uri) => {
		console.log("veg.explorer.copyPath.uri", uri.toString())
		await vscode.env.clipboard.writeText(uri.toString())
	})

	vscode.commands.registerCommand('veg.explorer.diffAll', async (args: any) => {
		console.log("veg.explorer.diffAll.args", args)
		vscode.window.showInformationMessage("Diff All (not implemented yet)")
	})

	vscode.commands.registerCommand('veg.explorer.refreshAll', async (args: any) => {
		console.log("veg.explorer.refreshAll.args", args)
		vcp.refreshAll()
		vscode.commands.executeCommand('workbench.files.action.refreshFilesExplorer')
	})


	extensionEmitter.event(async (e) => {
		// ...
		switch (e.type) {

			case "session.list.resp":
				vcp.setSessions(e.payload)
				break

			case "session.diff":
				console.log("filesys.session.diff", e.payload)
				const uriDiff = vscode.Uri.parse("oci://" + e.payload.currEnv)
				vcp.showDiff(uriDiff)
				break

			case "session.merge":
				console.log("filesys.session.merge", e.payload)
				const uriMerge = vscode.Uri.parse("oci://" + e.payload.currEnv)
				const dest = e.payload.dest ? vscode.Uri.parse(e.payload.dest) : undefined
				vcp.mergeDiff(uriMerge, dest, e.payload.forceInput)
				break

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
export async function deactivate() { }

// https://github.com/microsoft/vscode-extension-samples/blob/main/fsprovider-sample/src/fileSystemProvider.ts

// doing this so we can view diffs across agentic writes
// and accept or fork at various points
// todo, we need to persist this somewhere, the database or on the filesystem?
class VegContentProvider implements vscode.FileSystemProvider {

	private _environs: Record<string, any> = {}
	private _sessions: any[] = []
	private _latestEnvs: Map<string, string> = new Map()

	setSessions(sessions: any[]) {
		this._sessions = sessions
	}

	refreshAll() {
		const wsF = (vscode.workspace.workspaceFolders || []) as any[]
		for (let i = 0; i < wsF.length; i++) {
			const wf = wsF[i]
			if (wf.uri.scheme !== 'veg') {
				continue
			}

			// extract envId
			let p = wf.uri.authority + wf.uri.path
			if (p.startsWith("/")) p = p.slice(1)
			const lastColon = p.lastIndexOf(":")
			const envId = lastColon !== -1 ? p.substring(0, lastColon) : p

			const session = this._sessions.find(s => {
				const sEnv = s.state?.currEnv
				// console.log("  - ", sEnv, s)
				if (!sEnv) { return false }
				const lastColon = sEnv.lastIndexOf(":")
				const sId = lastColon !== -1 ? sEnv.substring(0, lastColon) : sEnv
				return sId === envId
			})

			if (session) {
				const f: Folder = {
					uri: vscode.Uri.parse("veg://" + session.state?.currEnv),
					name: session.state?.title || session.sid,
					sid: session.sid,
					session: session,
					environ: {
						fromUri: "oci://" + session.state?.currEnv,
						name: session.state?.title || session.sid,
					}
				}
				vscode.workspace.updateWorkspaceFolders(i, 1, f)
			}
		}
	}

	// we track a show (only) diff or everything
	private _onlyDiff: boolean = true
	toggleShown() {
		this._onlyDiff = !this._onlyDiff
		console.log("VEG.dagger.fs.toggleShown", this._onlyDiff)
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
		} catch (e: any) {
			console.error("catch!", e)
			throw e;
		}
	}

	private async makeReq(route: string, uri?: vscode.Uri, diffUri?: vscode.Uri, body?: any): Promise<Response> {
		// console.log("filesys.makeReq", route, uri, body)
		const url = `${SERVER_URL}${route}`
		var req = body
		if (uri && !body) {
			const ruri = this.vsUriToVeg(uri)
			req = {
				// idiots at microsoft encode query params, meaning = sign is %'d and
				// ACTUALLY compliant implementations of URLs ignore it... M$ idiots...
				uri: `${ruri.scheme}://${ruri.authority}${ruri.path}?${ruri.query}`,
				diff: this._onlyDiff,
			}
			if (!!diffUri) {
				const duri = this.vsUriToVeg(diffUri)
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

	// maybe path is not needed here, it is part of the uri, or we pass it around separately? there are places where we only have a single string to work with, which is why we started stuffing things into a URI, which is pretty flexible tbh
	// open(anyUri: string, path?: string, session?: any) {
	open(inputUri: string, session?: any) {
		console.log("filesys.open", inputUri, session)
		if ((!inputUri || inputUri === "") && !session) {
			return
		}

		// TODO, what if there is a session?
		// this should be a major condition below, because we are starting from one and can get most info there?

		// we need to do some logic here to clean things up
		// Uri parsing may be insufficient, it barfs on certain inputs
		// const anyUri = vscode.Uri.parse(inputUri)
		// So _i think_ this was switched to a more manual processing, though we could probably do some url parsing down under the hood?
		const f = async () => {
			var environ: Environ = {}

			var sid = ""

			// hacky parsing of user input to something our extension server understands
			if (inputUri?.startsWith("veg://")) {
				// hmmm, what do we do here
				// do we use a veg://session.env/<sid>/... or does the sid come earlier?
				// we want a way to always display the latest as it gets updated
				// and thus we also need a way to indicate to refresh when as session gets updated
			} else if (inputUri?.startsWith("file://")) {
				// no-op
				environ.srcUri = inputUri
				environ.name = inputUri.substring(7)
			} else if (inputUri?.startsWith("/")) {
				environ.srcUri = `file://${inputUri}`
				environ.name = inputUri

			} else if (inputUri?.startsWith("https://")) {
				// no-op
				environ.srcUri = inputUri
				environ.name = inputUri.substring(8)

			} else if (inputUri?.startsWith("oci://")) {
				// no-op
				environ.fromUri = inputUri
				environ.name = inputUri.substring(6)
			} else if (inputUri?.includes(":")) {
				// assume oci
				// look for no-domain, use docker.io (i.e. implied in docker pull)
				// the backend requires fully qualified oci://domain.com/reg/org/img:tag
				const parts = inputUri.split(":")
				const paths = parts[0].split("/")
				const host = paths[0]
				var ociUri = inputUri
				if (!host.includes(".")) {
					ociUri = `docker.io/${inputUri}`
				}
				environ.fromUri = `oci://${ociUri}`
				environ.name = inputUri
			} else {
				if (session?.state?.currEnv) {
					// these will always be oci
					sid = session.sid
					environ.name = session.state.title || sid
					// TODO, this should be a veg://<session>... something? we need to sort this out eventually, translators, more fields so we can differentiate, more alignment with server in open(...args) too?
					// seems we can add extra without borking things up?
					// this should probably be handled on the server during state managemtn
					environ.fromUri = "oci://" + session.state.currEnv
				} else {
					vscode.window.showErrorMessage(`unsupported environ: ${inputUri}`)
					return
				}
			}

			// what about subpaths... 
			var envUri = environ?.srcUri
			if (!envUri || envUri === "") {
				envUri = environ?.fromUri
			}
			if (!envUri || envUri === "") {
				vscode.window.showErrorMessage(`Bad error, see console for details: ${inputUri} -> ${envUri}`)
				console.error("shouldn't get here", inputUri, envUri, environ)
				return
			}
			const tmpUri = vscode.Uri.parse(envUri)
			const uri = vscode.Uri.from({
				...tmpUri,
				scheme: "veg",
			})

			const f: Folder = {
				uri,
				name: environ.name,
				sid,
				session,
				environ,
				// base: dir,
			}

			console.log("filesys.open.midway", environ, uri, f, sid, envUri)

			// if not a session, then we need to open it on the backend and get a uri
			if (inputUri && !session) {
				const resp = await this.makeReq("/fs/open", undefined, undefined, environ)
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
			}

			const wsF = vscode?.workspace?.workspaceFolders as any[]
			const count = vscode.workspace.workspaceFolders?.length || 0
			var pos = count

			var ws: vscode.WorkspaceFolder | any = null
			console.log("filesys.open.wsFolders:", f, wsF)

			// skip if already open, replace if matchind sid
			for (var i: number = 0; i < wsF.length; i++) {
				if (f.uri === wsF[i].uri) {
					vscode.window.showErrorMessage(`uri already open: ${f.uri}`)
					return
				}
				if (f.sid === wsF[i].sid) {
					pos = i
					break
				}
			}
			// console.log("filesys.open calling", f, count, uri)

			const started = vscode.workspace.updateWorkspaceFolders(pos, pos < count ? 1 : null, f)
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

	writeFile(uri: vscode.Uri, content: Uint8Array, options: { create: boolean, overwrite: boolean }): void | Thenable<void> {
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

	watch(uri: vscode.Uri, options: { excludes: readonly string[], recursive: boolean }): vscode.Disposable {
		const handler = () => {
			// DO NOT IMPLEMENT YET
		}
		return new vscode.Disposable(handler)
	}

	private _scms: Map<string, vscode.SourceControl> = new Map()
	private _resourceGroups: Map<string, Map<string, vscode.SourceControlResourceGroup>> = new Map()

	private getScmInfo(source: any): { uri?: vscode.Uri, session?: any, scmId?: string, groupId?: string } {
		// @ts-ignore
		if (!source.id || source.scheme) {
			return { uri: source as vscode.Uri };
		}

		// @ts-ignore
		const id = source.id;
		let uri: vscode.Uri | undefined;
		let session: any;
		let scmId: string | undefined;
		let groupId: string | undefined;

		// try to find which SCM/Group this is
		for (const [sid, scm] of this._scms.entries()) {
			if (scm === source) {
				scmId = sid;
				const latest = this._latestEnvs.get(sid);
				if (latest) uri = vscode.Uri.parse("veg://" + latest);
				session = this._sessions.find(s => s.sid === sid);
				break;
			}
		}

		if (!scmId) {
			for (const [sid, groups] of this._resourceGroups.entries()) {
				for (const [gid, group] of groups.entries()) {
					if (group === source) {
						scmId = sid;
						groupId = gid;
						uri = vscode.Uri.parse("veg://" + gid);
						session = this._sessions.find(s => s.sid === sid);
						break;
					}
				}
				if (scmId) break;
			}
		}

		if (!scmId) {
			// Fallback to old logic
			session = this._sessions.find(s => s.state?.currEnv === id || s.sid === id);
			if (session && id === session.sid) {
				scmId = session.sid;
				const sEnv = session.state?.currEnv;
				uri = vscode.Uri.parse("veg://" + (sEnv || id));
			} else {
				scmId = id;
				uri = vscode.Uri.parse("veg://" + id);
			}
		}

		return { uri, session, scmId, groupId };
	}

	// todo, we probably need a diffUri here
	showDiff(source: vscode.Uri | vscode.SourceControlResourceGroup, destination?: vscode.Uri): void | Thenable<void> {
		const f = async () => {
			console.log("filesys.showDiff.ARGS", source, destination)

			const info = this.getScmInfo(source);
			let uri = info.uri;
			let session = info.session;
			let envId = "";
			let envVer = "";

			if (!uri) {
				return
			}

			// extract envId and envVer from uri
			if (uri.scheme === 'veg' || uri.scheme === 'oci') {
				let p = uri.authority + uri.path
				if (p.startsWith("/")) p = p.slice(1)
				const lastColon = p.lastIndexOf(":")
				if (lastColon !== -1) {
					envId = p.substring(0, lastColon)
					envVer = p.substring(lastColon + 1)
				} else {
					envId = p
					envVer = "?"
				}
			}

			if (!session) {
				session = this._sessions.find(s => {
					const sEnv = s.state?.currEnv
					if (!sEnv) { return false }
					const lastColon = sEnv.lastIndexOf(":")
					const sId = lastColon !== -1 ? sEnv.substring(0, lastColon) : sEnv
					return sId === envId
				})
			}

			const scmId = info.scmId || session?.sid || envId
			const scmTitle = session?.state?.title || scmId
			let groupId = info.groupId || (uri.authority + uri.path)
			if (groupId.startsWith("/")) groupId = groupId.slice(1)
			const groupTitle = `${envVer} : ${envId}`

			// Track latest
			const currentLatest = this._latestEnvs.get(scmId);
			const currentVer = parseInt(envVer);
			if (!currentLatest || (!isNaN(currentVer) && currentVer > parseInt(currentLatest.split(":")[1]))) {
				this._latestEnvs.set(scmId, `${envId}:${envVer}`);
			}

			const resp = await this.makeReq("/fs/diff", uri)
			if (resp.status !== 200) {
				// console.error("filesys.showDiff.makeReq error:", resp)
				throw vscode.FileSystemError.FileNotFound(uri)
			}
			// console.log("filesys.showDiff.resp", uri, resp)

			const diff: any = await resp.json()
			console.log("filesys.showDiff.diff", diff)

			// get vscode uris for the environ basepath
			const prevUri = vscode.Uri.from({ ...vscode.Uri.parse(diff.prev), scheme: "veg" })
			const nextUri = vscode.Uri.from({ ...vscode.Uri.parse(diff.next), scheme: "veg" })

			let scm = this._scms.get(scmId)
			if (!scm) {
				scm = vscode.scm.createSourceControl("veg", scmTitle)
				scm.inputBox.visible = false
				this._scms.set(scmId, scm)
			}

			console.log("SCM id", scmId, "Group id", groupId)
			
			let scmGroups = this._resourceGroups.get(scmId)
			if (!scmGroups) {
				scmGroups = new Map()
				this._resourceGroups.set(scmId, scmGroups)
			}

			let group = scmGroups.get(groupId)
			if (!group) {
				group = scm.createResourceGroup(groupId, groupTitle)
				scmGroups.set(groupId, group)
			} else {
				group.label = groupTitle
			}
			
			const multiDiffResources: { originalUri: vscode.Uri | undefined; modifiedUri: vscode.Uri | undefined }[] = [];
			const resources: vscode.SourceControlResourceState[] = []

			// 1. Modified paths
			for (var p of diff.modPaths) {

				const pfUri = vscode.Uri.from({ ...prevUri, path: prevUri.path + p })
				const nfUri = vscode.Uri.from({ ...nextUri, path: nextUri.path + p })
				// console.log("diff", pfUri, nfUri)

				// 1. Add to SCM view (Single file diff)
				resources.push({
					resourceUri: nfUri,
					contextValue: 'modified',
					decorations: { 
						tooltip: `Modified: ${p}`,
						iconPath: new vscode.ThemeIcon('diff-modified', new vscode.ThemeColor('gitDecoration.modifiedResourceForeground')),
					},
					command: {
						command: 'vscode.diff',
						title: 'Show Diff',
						arguments: [pfUri, nfUri, `veg-diff: ${p}`]
					}
				});

				// 2. Collect for Multi Diff View
				multiDiffResources.push({
					originalUri: pfUri,
					modifiedUri: nfUri,
				});
			}

			// 2. Added paths
			for (var p of diff.addPaths) {
				if ((p.startsWith("/") && p.endsWith("/")) || p === "/stdout.txt" || p === "/stderr.txt") {
					continue
				}
				const nfUri = vscode.Uri.from({ ...nextUri, path: nextUri.path + p })

				resources.push({
					resourceUri: nfUri,
					contextValue: 'added',
					decorations: { 
						tooltip: `Added: ${p}`,
						iconPath: new vscode.ThemeIcon('diff-added', new vscode.ThemeColor('gitDecoration.addedResourceForeground')),
					},
					command: {
						command: 'vscode.open',
						title: 'Open File',
						arguments: [nfUri]
					}
				});

				multiDiffResources.push({
					// @ts-ignore
					originalUri: undefined,
					modifiedUri: nfUri,
				});
			}

			// 3. Deleted paths
			for (var p of diff.delPaths) {
				if ((p.startsWith("/") && p.endsWith("/")) || p === "/stdout.txt" || p === "/stderr.txt") {
					continue
				}
				const pfUri = vscode.Uri.from({ ...prevUri, path: prevUri.path + p })

				resources.push({
					resourceUri: pfUri,
					contextValue: 'deleted',
					decorations: { 
						tooltip: `Deleted: ${p}`, 
						strikeThrough: true,
						faded: true,
						iconPath: new vscode.ThemeIcon('diff-removed', new vscode.ThemeColor('gitDecoration.deletedResourceForeground')),
					},
					command: {
						command: 'vscode.open',
						title: 'Open File',
						arguments: [pfUri]
					}
				});

				multiDiffResources.push({
					originalUri: pfUri,
					// @ts-ignore
					modifiedUri: undefined,
				});
			}
			group.resourceStates = resources

			

			console.log("multiDiffResources:", multiDiffResources)

			// Define the options object exactly as the internal interface expects
			const options = {
				title: groupTitle,
				// This provides a unique ID for the tab so VS Code can manage it
				multiDiffSourceUri: vscode.Uri.from({ scheme: 'veg', path: `/diff-session/${groupId}` }),
				resources: multiDiffResources,
			};

			vscode.commands.executeCommand('workbench.view.scm')

			// Execute the internal workbench command
			await vscode.commands.executeCommand('_workbench.openMultiDiffEditor', options);



			return
		}
		return f()
	}

	// todo, we probably need a diffUri here
	mergeDiff(source: vscode.Uri | vscode.SourceControlResourceGroup, destination?: vscode.Uri, forceInput?: boolean): void | Thenable<void> {
		const f = async () => {
			const info = this.getScmInfo(source);
			let dest = destination
			let session = info.session;
			let uri = info.uri;
			let envId = "";
			let scmId = info.scmId || "";

			if (!uri) {
				return
			}

			// extract envId from uri
			if (uri.scheme === 'veg' || uri.scheme === 'oci') {
				let p = uri.authority + uri.path
				if (p.startsWith("/")) p = p.slice(1)
				const lastColon = p.lastIndexOf(":")
				if (lastColon !== -1) {
					envId = p.substring(0, lastColon)
				} else {
					envId = p
				}
			}

			if (!session) {
				session = this._sessions.find(s => {
					const sEnv = s.state?.currEnv
					if (!sEnv) { return false }
					const lastColon = sEnv.lastIndexOf(":")
					const sId = lastColon !== -1 ? sEnv.substring(0, lastColon) : sEnv
					return sId === envId
				})
			}

			if (!scmId) scmId = session?.sid || envId

			// Track latest
			let envVer = ""
			if (uri.scheme === 'veg') {
				let p = uri.path
				if (p.startsWith("/")) p = p.slice(1)
				envVer = p.split("/")[0].split(":")[1] || "?"
			} else if (uri.scheme === 'oci') {
				envVer = uri.path.split("/")[1].split(":")[1] || "?"
			}
			const currentLatest = this._latestEnvs.get(scmId);
			const currentVer = parseInt(envVer);
			if (!currentLatest || (!isNaN(currentVer) && currentVer > parseInt(currentLatest.split(":")[1]))) {
				this._latestEnvs.set(scmId, `${envId}:${envVer}`);
			}

			if (!dest && !forceInput) {
				if (session?.state?.initEnv?.srcUri?.startsWith("file://")) {
					dest = vscode.Uri.parse(session.state.initEnv.srcUri)
				}
			}

			if (!dest || forceInput) {
				const value = await vscode.window.showInputBox({
					title: "Merge Diff",
					prompt: "Pick a directory or Veg environ to merge into",
					placeHolder: "/path/on/disk/... | veg://...",
					value: session?.state?.initEnv?.srcUri?.startsWith("file://") ? session.state.initEnv.srcUri : "",
				})
				if (!value) {
					return
				}
				dest = vscode.Uri.parse(value as string)
			}

			if (dest.scheme !== 'file') {
				vscode.window.showErrorMessage(`unsupported target, only file://: ${dest}`)
				return
			}

			console.log("filesys.mergeDiff.args", source, dest)
			const resp = await this.makeReq("/fs/diff", uri)
			if (resp.status !== 200) {
				// console.error("filesys.mergeDiff.makeReq error:", resp)
				throw vscode.FileSystemError.FileNotFound(uri)
			}
			// console.log("filesys.mergeDiff.resp", uri, resp)

			const diff: any = await resp.json()
			console.log("filesys.mergeDiff.diff", diff)

			// write to disk
			if (dest.scheme === 'file') {
				for (var path of diff.addPaths) {
					// skip ugh...
					if ((path.startsWith("/") && path.endsWith("/")) || path === "/stdout.txt" || path === "/stderr.txt") {
						continue
					}
					const val = diff.files[path]
					const key = dest.path + path
					await fs.writeFile(key, val)
				}

				for (var path of diff.modPaths) {
					// skip ugh...
					if ((path.startsWith("/") && path.endsWith("/")) || path === "/stdout.txt" || path === "/stderr.txt") {
						continue
					}
					const val = diff.files[path]
					const key = dest.path + path
					await fs.writeFile(key, val)
				}

				for (var path of diff.delPaths) {
					// skip ugh...
					if ((path.startsWith("/") && path.endsWith("/")) || path === "/stdout.txt" || path === "/stderr.txt") {
						continue
					}
					const key = dest.path + path
					await fs.rm(key)
				}
			}

			return
		}
		return f()
	}

	hideDiff(source: vscode.Uri | vscode.SourceControlResourceGroup) {
		const info = this.getScmInfo(source);
		let scmId = info.scmId || "";
		let groupId = info.groupId || "";
		let isScm = !groupId;

		if (!scmId) {
			// Fallback if getScmInfo didn't find it in maps
			// @ts-ignore
			if (source.id && !source.scheme) {
				// @ts-ignore
				const id = source.id
				// check if it is an SCM provider or a resource group
				// @ts-ignore
				if (source.resourceStates === undefined) {
					isScm = true
					scmId = id
				} else {
					groupId = id
				}

				const lastColon = id.lastIndexOf(":")
				const envIdFromId = lastColon !== -1 ? id.substring(0, lastColon) : id

				const session = this._sessions.find(s => {
					if (s.sid === id || s.state?.currEnv === id) return true
					const sEnv = s.state?.currEnv
					if (!sEnv) return false
					const sLastColon = sEnv.lastIndexOf(":")
					const sId = sLastColon !== -1 ? sEnv.substring(0, sLastColon) : sEnv
					return sId === envIdFromId
				})

				if (session) {
					scmId = session.sid
					if (!isScm) groupId = id
				} else {
					scmId = envIdFromId
					if (!isScm) groupId = id
				}
			} else {
				const uri = source as vscode.Uri
				// extract envId
				let envId = "?"
				if (uri.scheme === 'veg' || uri.scheme === 'oci') {
					let p = uri.authority + uri.path
					if (p.startsWith("/")) p = p.slice(1)
					const lastColon = p.lastIndexOf(":")
					if (lastColon !== -1) {
						envId = p.substring(0, lastColon)
					} else {
						envId = p
					}
				}

				const session = this._sessions.find(s => {
					const sEnv = s.state?.currEnv
					if (!sEnv) { return false }
					const lastColon = sEnv.lastIndexOf(":")
					const sId = lastColon !== -1 ? sEnv.substring(0, lastColon) : sEnv
					return sId === envId
				})

				scmId = session?.sid || envId
				groupId = uri.authority + uri.path
				if (groupId.startsWith("/")) groupId = groupId.slice(1)
			}
		}
		console.log("hideDiff SCM id", scmId, "Group id", groupId, "isScm", isScm)

		const scmGroups = this._resourceGroups.get(scmId)
		if (scmGroups) {
			if (isScm) {
				// Hide all groups in this SCM
				for (const group of scmGroups.values()) {
					group.dispose()
				}
				scmGroups.clear()
			} else {
				const group = scmGroups.get(groupId)
				if (group) {
					group.dispose()
					scmGroups.delete(groupId)
				}
			}
			
			if (scmGroups.size === 0) {
				this._resourceGroups.delete(scmId)
				const scm = this._scms.get(scmId)
				if (scm) {
					scm.dispose()
					this._scms.delete(scmId)
				}
			}
		} else if (isScm) {
			// Even if no groups, dispose the SCM if it exists
			const scm = this._scms.get(scmId)
			if (scm) {
				scm.dispose()
				this._scms.delete(scmId)
			}
		}
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
