import * as vscode from 'vscode';
import { extensionEmitter, sendMessage } from '../comms';
import { makeReq, normalizeEnvId, Environ, Folder, parseEnvUri, findSession, vsUriToVeg } from './utils';
import * as scm from './scmProvider';

// This method is called when your extension is activated
// Your extension is activated the very first time the command is executed
export async function activate(context: vscode.ExtensionContext) {

  var userInput: any = null

	const vcp = new VegContentProvider()

	vscode.workspace.registerFileSystemProvider("veg", vcp, {
		isCaseSensitive: true,
		isReadonly: false,
	})
	
	// Activate SCM provider
	await scm.activate(context);

	vscode.commands.registerCommand('veg.filesys.hack', async (arg1: any) => {
		console.log("veg.filesys.hack", arg1)
	})

	vscode.commands.registerCommand('veg.explorer.chat', async (uri: vscode.Uri | vscode.SourceControl) => {
		console.log("veg.explorer.chat.args", uri)

    // parse inputs
    if (uri instanceof vscode.Uri) {
      uri = uri
    } else if (uri?.rootUri instanceof vscode.Uri) {
      uri = uri.rootUri
    } else {
      console.log("veg.explorer.chat.unknown-type:", typeof uri)
      console.log("veg.explorer.chat.args", uri)
      return
    }
		console.log("veg.explorer.chat.args", uri)

    // construct message
		const msg: any = {
			type: "session.create",
			payload: {
				focus: true,
        agent: userInput?.agent,
        model: userInput?.model,
        envName: userInput?.environ,
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

    // send message
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

	vscode.commands.registerCommand('veg.explorer.copyPath', async (uri: vscode.Uri) => {
		console.log("veg.explorer.copyPath.uri", uri.toString())
		await vscode.env.clipboard.writeText(uri.toString())
	})

	vscode.commands.registerCommand('veg.explorer.saveToVeg', async (uri: vscode.Uri) => {
		console.log("veg.explorer.saveToVeg.args", uri)

		const value = await vscode.window.showInputBox({
			title: "Save to Veg",
			prompt: "Enter destination veg:// URI (folder)",
			placeHolder: "veg://...",
		})
		if (!value) return

		const destRoot = vscode.Uri.parse(value)
		if (destRoot.scheme !== 'veg') {
			vscode.window.showErrorMessage("Destination must be a veg:// URI")
			return
		}
		
		const srcName = uri.path.split('/').pop() || "unknown"
		const finalDest = vscode.Uri.joinPath(destRoot, srcName)

		const copyRec = async (src: vscode.Uri, dst: vscode.Uri) => {
			const stat = await vscode.workspace.fs.stat(src)
			if (stat.type & vscode.FileType.Directory) {
				await vscode.workspace.fs.createDirectory(dst)
				const entries = await vscode.workspace.fs.readDirectory(src)
				for (const [name, type] of entries) {
					await copyRec(vscode.Uri.joinPath(src, name), vscode.Uri.joinPath(dst, name))
				}
			} else {
				const content = await vscode.workspace.fs.readFile(src)
				await vscode.workspace.fs.writeFile(dst, content)
			}
		}

		await vscode.window.withProgress({
			location: vscode.ProgressLocation.Notification,
			title: `Uploading ${srcName}...`,
			cancellable: false
		}, async () => {
			try {
				await copyRec(uri, finalDest)
				vscode.window.showInformationMessage(`Uploaded ${srcName}`)
			} catch(e: any) {
				vscode.window.showErrorMessage(`Upload failed: ${e.message}`)
			}
		})
	})

	vscode.commands.registerCommand('veg.explorer.refreshAll', async (args: any) => {
		console.log("veg.explorer.refreshAll.args", args)
		vcp.refreshAll()
		vscode.commands.executeCommand('workbench.files.action.refreshFilesExplorer')
	})

	vscode.commands.registerCommand('veg.explorer.terminal', async (arg: any) => {
		console.log("veg.explorer.terminal.args", arg)
		vscode.window.showInformationMessage(`Veg Terminal triggered: ${arg?.id || 'no id'}`);
		scm.scmProvider.openTerminal(arg)
	})


	extensionEmitter.event(async (e) => {
		// ...
		switch (e.type) {

			case "chat.userInput.resp":
				userInput = e.payload
				break

			case "session.list.resp":
				vcp.setSessions(e.payload)
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
						vcp.open("", session, e.payload.pos)
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

	private _sessions: any[] = []

	setSessions(sessions: any[]) {
		this._sessions = sessions
		this.refreshAll()
	}

	refreshAll() {
		const wsF = (vscode.workspace.workspaceFolders || []) as any[]
		for (let i = 0; i < wsF.length; i++) {
			const wf = wsF[i]
			if (wf.uri.scheme !== 'veg') {
				continue
			}

			// extract envId using helper
			const { envId, envVer } = parseEnvUri(wf.uri)
			const q = new URLSearchParams(wf.uri.query)
			const uriPos = q.get("pos")

			// use helper
			const session = findSession(this._sessions, envId)

			if (session) {
				const sEnv = session.state?.currEnv
				const { envVer: sVer } = sEnv ? parseEnvUri(vscode.Uri.parse("oci://" + sEnv)) : { envVer: "" }

				let nextUri = wf.uri
				if (sVer !== "" && sVer !== envVer) {
					const tmpUri = vscode.Uri.parse("oci://" + sEnv)
					nextUri = vscode.Uri.from({
						...tmpUri,
						scheme: "veg",
						query: wf.uri.query,
					})
				}

				let name = session.state?.title || session.sid

				const ver = sVer !== "" ? sVer : envVer
				if (ver && ver !== "" && ver !== "?" && !isNaN(parseInt(ver))) {
					name = `(${ver}) ${name}`
				}

				let posVal = uriPos
				if (!posVal) {
					const posMatch = wf.name.match(/\[(\d+)\]/);
					if (posMatch) posVal = posMatch[1]
				}

				if (posVal) {
					name = `[${posVal}] ${name}`;
				}

				const f: Folder = {
					uri: nextUri,
					name: name,
					sid: session.sid,
					session: session,
					environ: {
						fromUri: "oci://" + (nextUri.authority + nextUri.path).replace(/^\//, ''),
						name: name,
					}
				}
				// only update if something changed to avoid unnecessary refreshes
				if (wf.uri.toString() !== f.uri.toString() || wf.name !== f.name) {
					vscode.workspace.updateWorkspaceFolders(i, 1, f)
				}
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


	// maybe path is not needed here, it is part of the uri, or we pass it around separately? there are places where we only have a single string to work with, which is why we started stuffing things into a URI, which is pretty flexible tbh
	// open(anyUri: string, path?: string, session?: any) {
	open(inputUri: string, session?: any, pos?: number) {
		console.log("filesys.open", inputUri, session, pos)
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
				// HMMM, perhaps this needs to be the first checked condition?
				if (session?.state?.currEnv) {
					// these will always be oci
					sid = session.sid
					environ.name = session.state.title || sid

					const img = session.state.currEnv
					const lastColon = img.lastIndexOf(":")
					const tag = lastColon !== -1 ? img.substring(lastColon + 1) : ""
					if (tag !== "" && !isNaN(parseInt(tag))) {
						environ.name = `(${tag}) ${environ.name}`
					}
					if (pos !== undefined && pos >= 0) {
						environ.name = `[${pos}] ${environ.name}`
					}
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
				query: pos !== undefined ? `pos=${pos}` : tmpUri.query,
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
				const resp = await makeReq("/fs/open", undefined, undefined, environ, this._onlyDiff)
				// console.log("filesys.resp:", resp)
				if (resp.status !== 200) {
					vscode.window.showErrorMessage(`${resp.status} - ${resp.statusText}`)
					return
				}

				const data: any = await resp.json()
				console.log("filesys.open.api.resp:", data)

				const vegUri = vscode.Uri.parse(data.envUri).with({ scheme: 'veg' })
				f.uri = vegUri
				// const parts = uri.path.split(":")
				// var sid = parts[0]
				// if (sid.startsWith("/")) {
				// 	sid = sid.substring(1)
				// }
				// const tag = parts[1]

				// convert to something vscode will understand
			}

			const wsF = (vscode?.workspace?.workspaceFolders || []) as any[]
			const count = wsF.length
			var wsIndex = count

			// skip if already open, replace if matching sid
			const { envId: fEnvId } = parseEnvUri(f.uri)
			const normFEnvId = normalizeEnvId(fEnvId)

			for (var i: number = 0; i < wsF.length; i++) {
				const wf = wsF[i]
				if (wf.uri.scheme !== 'veg') continue

				if (f.uri.toString() === wf.uri.toString()) {
					vscode.window.showErrorMessage(`uri already open: ${f.uri}`)
					return
				}

				const { envId: wfEnvId } = parseEnvUri(wf.uri)
				if (normFEnvId === normalizeEnvId(wfEnvId)) {
					wsIndex = i
					// try to preserve pos from name if not provided
					if (pos === undefined) {
						const posMatch = wf.name.match(/\[(\d+)\]/);
						if (posMatch) {
							f.name = `[${posMatch[1]}] ${f.name}`
						}
					}
					break
				}
			}
			// console.log("filesys.open calling", f, count, uri)

			const started = vscode.workspace.updateWorkspaceFolders(wsIndex, wsIndex < count ? 1 : null, f)
			// console.log("filesys.open started?", started)
		}
		return f()
	}

	stat(uri: vscode.Uri): vscode.FileStat | Thenable<vscode.FileStat> {
		// console.log("fs.stat.uri", uri)
		const f = async () => {
			const ociUri = vsUriToVeg(uri)
			const { envId } = parseEnvUri(ociUri)
			const session = findSession(this._sessions, envId, ociUri)
			const sid = session?.sid

			const resp = await makeReq("/fs/stat", uri, undefined, undefined, this._onlyDiff, sid)
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
			const ociUri = vsUriToVeg(uri)
			const { envId } = parseEnvUri(ociUri)
			const session = findSession(this._sessions, envId, ociUri)
			const sid = session?.sid

			const resp = await makeReq("/fs/read", uri, undefined, undefined, this._onlyDiff, sid)
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
			console.log("readDirectory.uri", uri.toString())
			const ociUri = vsUriToVeg(uri)
			const { envId } = parseEnvUri(ociUri)
			const session = findSession(this._sessions, envId, ociUri)
			const sid = session?.sid

			const resp = await makeReq("/fs/list", uri, undefined, undefined, this._onlyDiff, sid)
			if (resp.status !== 200) {
				console.error("readDirectory.makeReq error:", resp)
				throw vscode.FileSystemError.FileNotFound(uri)
			}
			// console.log("filesys.readDir.resp", uri, resp)

			const data: any = await resp.json()
			console.log("readDirectory.data", uri.toString(), data?.entries?.length, JSON.stringify(data?.entries))

			// our returned listing
			var l: FolderListing = []
			for (const entry of data?.entries) {
				// console.log("filesys.readDir.entry", entry)
				let name = entry.name
				// if absolute path, we want just the name
				if (name.includes("/")) {
					name = name.split("/").pop() || name
				}
				if (name === "") continue;

				const isDir = entry.dir
				l.push([name, isDir ? vscode.FileType.Directory : vscode.FileType.File])
			}
			console.log("readDirectory.list", uri.toString(), JSON.stringify(l))

			return l
		}
		return f()
	}

	async writeFile(uri: vscode.Uri, content: Uint8Array, options: { create: boolean, overwrite: boolean }): Promise<void> {
		console.log("VegContentProvider.writeFile", uri.toString(), options)
		const ociUri = vsUriToVeg(uri)
		const params = new URLSearchParams(ociUri.query)
		const path = params.get("path") || ""

		const { envId } = parseEnvUri(ociUri)
		const session = findSession(this._sessions, envId, ociUri)
		const sid = session?.sid

		// uri should not include path, need to reconstruct

		const resp = await makeReq("/fs/write", uri, undefined, {
			uri: ociUri.toString(),
			path: path,
			content: new TextDecoder().decode(content),
			sid: sid,
		}, this._onlyDiff)

		if (resp.status !== 200) {
			throw vscode.FileSystemError.Unavailable(uri)
		}

		const data: any = await resp.json()
		const nextUri = vscode.Uri.parse(data.envUri).with({ scheme: 'veg' })
		this.updateUri(uri, nextUri)
	}


	async createDirectory(uri: vscode.Uri): Promise<void> {
		const ociUri = vsUriToVeg(uri)
		const params = new URLSearchParams(ociUri.query)
		const path = params.get("path") || ""

		const { envId } = parseEnvUri(ociUri)
		const session = findSession(this._sessions, envId, ociUri)
		const sid = session?.sid

		const resp = await makeReq("/fs/mkdir", uri, undefined, {
			uri: ociUri.toString(),
			path: path,
			sid: sid,
		}, this._onlyDiff)

		if (resp.status !== 200) {
			throw vscode.FileSystemError.Unavailable(uri)
		}

		const data: any = await resp.json()
		const nextUri = vscode.Uri.parse(data.envUri).with({ scheme: 'veg' })
		this.updateUri(uri, nextUri)
	}

	async delete(uri: vscode.Uri, options: { recursive: boolean }): Promise<void> {
		const ociUri = vsUriToVeg(uri)
		const params = new URLSearchParams(ociUri.query)
		const path = params.get("path") || ""

		const { envId } = parseEnvUri(ociUri)
		const session = findSession(this._sessions, envId, ociUri)
		const sid = session?.sid

		const resp = await makeReq("/fs/delete", uri, undefined, {
			uri: ociUri.toString(),
			path: path,
			sid: sid,
		}, this._onlyDiff)

		if (resp.status !== 200) {
			throw vscode.FileSystemError.Unavailable(uri)
		}

		const data: any = await resp.json()
		const nextUri = vscode.Uri.parse(data.envUri).with({ scheme: 'veg' })
		this.updateUri(uri, nextUri)
	}

	async rename(source: vscode.Uri, destination: vscode.Uri, options: { overwrite: boolean }): Promise<void> {
		const ociSrc = vsUriToVeg(source)
		const srcParams = new URLSearchParams(ociSrc.query)
		
		const ociDst = vsUriToVeg(destination)
		const dstParams = new URLSearchParams(ociDst.query)

		const { envId } = parseEnvUri(ociSrc)
		const session = findSession(this._sessions, envId, ociSrc)
		const sid = session?.sid

		const resp = await makeReq("/fs/rename", source, undefined, {
			uri: ociSrc.toString(),
			src: srcParams.get("path") || "",
			dst: dstParams.get("path") || "",
			sid: sid,
		}, this._onlyDiff)

		if (resp.status !== 200) {
			throw vscode.FileSystemError.Unavailable(source)
		}

		const data: any = await resp.json()
		const nextUri = vscode.Uri.parse(data.envUri).with({ scheme: 'veg' })
		this.updateUri(source, nextUri)
	}

	async copy(source: vscode.Uri, destination: vscode.Uri, options: { overwrite: boolean }): Promise<void> {
		const ociSrc = vsUriToVeg(source)
		const srcParams = new URLSearchParams(ociSrc.query)
		
		const ociDst = vsUriToVeg(destination)
		const dstParams = new URLSearchParams(ociDst.query)

		const { envId } = parseEnvUri(ociSrc)
		const session = findSession(this._sessions, envId, ociSrc)
		const sid = session?.sid

		const resp = await makeReq("/fs/copy", source, undefined, {
			uri: ociSrc.toString(),
			src: srcParams.get("path") || "",
			dst: dstParams.get("path") || "",
			sid: sid,
		}, this._onlyDiff)

		if (resp.status !== 200) {
			throw vscode.FileSystemError.Unavailable(source)
		}

		const data: any = await resp.json()
		const nextUri = vscode.Uri.parse(data.envUri).with({ scheme: 'veg' })
		this.updateUri(source, nextUri)
	}

	private updateUri(oldUri: vscode.Uri, nextUri: vscode.Uri) {
		const wsF = (vscode.workspace.workspaceFolders || []) as any[]
		for (let i = 0; i < wsF.length; i++) {
			const wf = wsF[i]
			if (wf.uri.scheme !== 'veg') continue

			const { envId: oldEnvId } = parseEnvUri(oldUri)
			const { envId: wfEnvId } = parseEnvUri(wf.uri)

			if (oldEnvId === wfEnvId) {
				const { envVer } = parseEnvUri(nextUri)
				const q = new URLSearchParams(wf.uri.query)
				const uriPos = q.get("pos")

				const session = findSession(this._sessions, wfEnvId, nextUri)
				let name = session?.state?.title || session?.sid || wfEnvId

				if (envVer && envVer !== "" && envVer !== "?" && !isNaN(parseInt(envVer))) {
					name = `(${envVer}) ${name}`
				}

				let posVal = uriPos
				if (!posVal) {
					const posMatch = wf.name.match(/\[(\d+)\]/);
					if (posMatch) posVal = posMatch[1]
				}

				if (posVal) {
					name = `[${posVal}] ${name}`;
				}

				const f: Folder = {
					uri: nextUri,
					name: name,
					sid: session?.sid || wfEnvId,
					session: session,
					environ: {
						fromUri: "oci://" + (nextUri.authority + nextUri.path).replace(/^\//, ''),
						name: name,
					}
				}
				vscode.workspace.updateWorkspaceFolders(i, 1, f)
				break
			}
		}
	}

	watch(uri: vscode.Uri, options: { excludes: readonly string[], recursive: boolean }): vscode.Disposable {
		const listener = extensionEmitter.event(e => {
			// console.log("veg.fs.watch.event", e)
			if (e.type === "filesys.change") {
				// TODO: check if e.payload matches watched uri
				this._emitter.fire([{ type: vscode.FileChangeType.Changed, uri }])
			}
		})
		return listener
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

type FolderListing = Array<[string, vscode.FileType]>
