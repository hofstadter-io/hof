import fs from 'node:fs/promises'
import * as vscode from 'vscode';
import { extensionEmitter } from '../util/events';

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
export function activate(context: vscode.ExtensionContext) {

	const vcp = new VegContentProvider()

	vscode.workspace.registerFileSystemProvider("veg", vcp, {
		isCaseSensitive: true,
		isReadonly: false,
	})

	context.subscriptions.push(
		vscode.workspace.onDidOpenTextDocument(e => {
			if (e.uri.scheme !== 'file') {
				return;
			}
			// console.log("openedDocument", e.fileName)
		})
	);

	vscode.commands.registerCommand('veg.explorer.refreshAll', async () => {
		await vcp.refreshAll()
		vscode.commands.executeCommand('workbench.files.action.refreshFilesExplorer')
  })

	vscode.commands.registerCommand('veg.explorer.toggleShown', () => {
    vcp.toggleShown()
		vscode.commands.executeCommand('workbench.files.action.refreshFilesExplorer')
  })

	vscode.commands.registerCommand('veg.session.showSessionDiff', () => {
		const sid = context.workspaceState.get('sid')
		if (!sid || sid === '') {
			return
		}
		vcp.showDiff(sid as string)
	})

	vscode.commands.registerCommand('veg.session.showFileDiff', (uri: vscode.Uri) => {
		console.log("VEG.showFileDiff", uri)
		const qp = new URLSearchParams(uri.query)
		var sid = qp.get("sid") as string
		if (!sid || sid === '') {
			return
		}
		console.log("DIFF:", uri.path)
		const fileUri = vscode.Uri.parse(`file://${uri.path}`)
		vscode.commands.executeCommand('vscode.diff', fileUri, uri, `veg-diff: ${uri.path}`, {
			preserveFocus: false,
			preview: false,
			viewColumn: 1,
		})

		// const uri = node

		// 		var vegUri = vscode.Uri.parse(`veg://${p}`)
		// 		vegUri = vegUri.with({
		// 			query: `sid=${sid}`
		// 		})
		// 		// console.log("veg.diff", fileUri, vegUri)
		// vcp.showDiff(sid as string)
	});

	extensionEmitter.event(async (e) => {
		if (e.type === "session.fs.open") {
			var { sid, dir, name }: { sid: string, dir: string, name?: string } = e.payload;
			vcp.openFS(sid, dir, name)
			return
		}
		if (e.type === "session.fs.refresh") {
			var { sid, pos: p }: { sid: string, pos?: string } = e.payload;
			vcp.refreshFS(sid, p)
			return
		}

		if (e.type === "session.merge") {
			var { sid }: { sid: string, dir: string, name?: string } = e.payload;
			vcp.mergeFS(sid)
			return
		}

		if (e.type === "session.delete") {
			var { sid }: { sid: string } = e.payload;
			// console.log("FS.delete", sid)
			vcp.closeFS(sid)
			return
		}

		if (e.type === "session.list.resp") {
			vcp.setSessions(e.payload)
		}
		if (e.type === "session.diff.show") {
			var { sid, pos}: {
				sid: string,
				pos: string
			} = e.payload;
			if (pos && pos !== '') {
				sid += "-" + pos
			}
			vcp.showDiff(sid)
			return
		}

		if (e.type === "session.diff.resp") {
			var { sid, pos, show, files}: {
				sid: string,
				pos: string,
				show: boolean,
				files: Record<string,string>,
			} = e.payload;

			if (pos && pos !== '') {
				sid += "-" + pos
			}

			vcp.setDiff(sid, e.payload)

			if (show) {
				vcp.showDiff(sid)
			}
		}
	});

  // const fsdeco = new VegFileDecorationProvider();
  // context.subscriptions.push(
  //   vscode.window.registerFileDecorationProvider(fsdeco)
  // );

	// extensionEmitter.event(async (e) => {
	// 	if (e.type === "session.diff.resp") {
  //     fsdeco.setDiff(e.payload.sid, e.payload)
  //   }
  // })
}

// This method is called when your extension is deactivated
export function deactivate() {}

// https://github.com/microsoft/vscode-extension-samples/blob/main/fsprovider-sample/src/fileSystemProvider.ts

// doing this so we can view diffs across agentic writes
// and accept or fork at various points
// todo, we need to persist this somewhere, the database or on the filesystem?
class VegContentProvider implements vscode.FileSystemProvider {

	// we keep a mock fs per session
	// <session>/<path> = <content>
	private fs: Record<string, Record<string, string>> = {}
	private diff: Record<string,any> = {}
	private _sessions: Record<string,any> = {}


	private _emitter = new vscode.EventEmitter<vscode.FileChangeEvent[]>();
	private _bufferedEvents: vscode.FileChangeEvent[] = [];
	private _fireSoonHandle?: NodeJS.Timeout;

	// constructor(){
	// 	this.fs = {}
	// 	this.diff = {}
	// }

	readonly onDidChangeFile: vscode.Event<vscode.FileChangeEvent[]> = this._emitter.event

  private _shown: boolean = true
  toggleShown() {
    this._shown = !this._shown
		console.log("VEG.vcp.toggleShown", this._shown)
  }

	copy(source: vscode.Uri, destination: vscode.Uri, options: { overwrite: boolean }): void | Thenable<void> {

	}

	createDirectory(uri: vscode.Uri): void | Thenable<void> {

	}

	delete(uri: vscode.Uri, options: { recursive: boolean }): void | Thenable<void> {

	}

	readDirectory(uri: vscode.Uri): FolderListing | Thenable<FolderListing> {
		const f = async () => {
			// console.log("VEG.readDir.uri", uri)
			const s = this.uriToSid(uri)
			this.ensureExists(s)
			const diff = this.diff[s]
			// console.log("VEG.readDir.fs", fs)

			const url = `${SERVER_URL}/fs/list`
			var req = {
				sid: s,
				path: uri.path,
			}
			const resp = await fetch(url, {
				method: "POST",
				headers: {
					"Content-Type": "application/json",
				},
				body: JSON.stringify(req)
			})
			// console.log("VEG.readDir.resp", resp)

			const dag: any = await resp.json()
			// console.log("VEG.readDir.dag", dag)
			const entries: Array<[string,boolean]> = dag?.entries || []

			// our returned listing
			var l: FolderListing = []

			// From Dagger via extension server
			for (const pair of entries) {
				const path = pair[0]
				const isDir = pair[1]
				if (!this._shown) {
					var p = uri.path
					if (!p.endsWith("/") && !path.startsWith("/")) {
						p += "/" + path
					} else {
						p += path
					}
					const match = matchPathInDiff(p, diff)
					// console.log("show&tell", path, isDir, match)
					if (!match) {
						continue
					}
				}
				l.push([path, isDir ? vscode.FileType.Directory : vscode.FileType.File])
			}

			// From our FS, we'd want to update the time from our fs, which we also track (soon(tm))
			//
			// for (const path in fs) {
			// 	if (!path.startsWith(uri.path)) {
			// 		continue
			// 	}
			// 	const rel = path.replace(uri.path + "/", "")
			// 	const parts = rel.split("/")
			// 	const part = parts[0]
			// 	// hack for now, files have a .
			// 	const isDir = !part.includes(".")
			// 	console.log("VEG.readDir.loop", path, rel, parts, part, isDir)
			// 	// const p = path
			// 	l.push([part, isDir ? vscode.FileType.Directory : vscode.FileType.File])
			// }
			return l
		}
		return f()
	}

	readFile(uri: vscode.Uri): Uint8Array | Thenable<Uint8Array> {
		const f = async () => {
			// console.log("VEG.readFile.uri", uri)
			const s = this.uriToSid(uri)

			const url = `${SERVER_URL}/fs/read`
			var req = {
				sid: s,
				path: uri.path,
			}
			const resp = await fetch(url, {
				method: "POST",
				headers: {
					"Content-Type": "application/json",
				},
				body: JSON.stringify(req)
			})
			// console.log("VEG.readFile.resp", resp)

			const dag: any = await resp.json()
			// console.log("VEG.readFile.dag", dag)
			const c: string = dag.contents

			return new TextEncoder().encode(c)

		}
		return f()
	}

	rename(source: vscode.Uri, destination: vscode.Uri, options: { overwrite: boolean }): void | Thenable<void> {
		const s = this.uriToSid(source)
		this.ensureExists(s)

		const c = this.fs[s][source.path]
		this.fs[s][destination.path] = c
	}

	stat(uri: vscode.Uri): vscode.FileStat | Thenable<vscode.FileStat> {

		const f = async () => {
			// console.log("VEG.stat.uri", uri)
			const s = this.uriToSid(uri)
			const path = uri.path
			const url = `${SERVER_URL}/fs/stat`
			var req = {
				sid: s,
				path,
			}
			// console.log("VEG.readFile.resp", req)
			const resp = await fetch(url, {
				method: "POST",
				headers: {
					"Content-Type": "application/json",
				},
				body: JSON.stringify(req)
			})
			// console.log("VEG.readFile.resp", resp)

			if (resp.status !== 200) {
				// console.error("VEG.stat.fetch.error", resp)
				throw vscode.FileSystemError.FileNotFound(uri)
			}

			const diff: any = this.diff[s]
			const dag: any = await resp.json()
			// console.log("VEG.stat.dag", dag)
			// console.log("VEG.stat", path, diff.addpaths, diff.modpaths)
			var mtime = dag.mtime
			// if path in diff, update mtime
			if (diff?.addpaths?.includes(path) || diff?.modpaths?.includes(path)) {
				// console.log("VEG.modified!", path)
				mtime = new Date().getMilliseconds()
			}

			// todo
			return {
				ctime: dag.ctime,
				mtime: dag.mtime,
				size: dag.size,
				type: dag.dir ? vscode.FileType.Directory : vscode.FileType.File,
			}

		}
		return f()
	}

	watch(uri: vscode.Uri, options: {excludes: readonly string[], recursive: boolean}): vscode.Disposable {
		const handler = () => {

		}
		return new vscode.Disposable(handler)
	}

	writeFile(uri: vscode.Uri, content: Uint8Array, options: {create: boolean, overwrite: boolean}): void | Thenable<void> {
		const s = this.uriToSid(uri)
		this.ensureExists(s)

		const c = new TextDecoder().decode(content)
		this.fs[s][uri.path] = c
	}	


	/*********************************/

	// cleanPath(path: string) {
	// 	var wDir = ""
	// 	const wsF = vscode.workspace.workspaceFolders
	// 	if (wsF && wsF.length > 0) {
	// 		wDir = wsF[0].uri.path	
	// 	}
	// 	return (wDir != "" && path.startsWith(wDir)) ? path.substring(wDir.length + 1) : path

	// }

	openFS(sid: string, dir: string, name?: string) {
		// console.log("openFS called:", dir)
		var vegUri = vscode.Uri.parse(`veg://${dir}`)
		vegUri = vegUri.with({
			query: `sid=${sid}`
		})
		const sess = this._sessions[sid]
		if ((!name || name === "") && sess) {
			name = sess?.state?.title || sid
		}

		const f: Folder = {
			uri: vegUri,
			sid: sid,
			name: name || sid,
			base: dir,
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
		// console.log("openFS calling", f, count, vegUri)

		const started = vscode.workspace.updateWorkspaceFolders(count, null, f)
		// console.log("openFS started?", started)

	}

	async refreshAll() {
		const wsF = vscode?.workspace?.workspaceFolders as any[]

		var ws: vscode.WorkspaceFolder | any = null
		wsF.forEach((ws, i)=>{
			const qp = new URLSearchParams(ws.uri.query)
			const s = qp.get("sid") as string
			const S = this._sessions[s]
			if (S) {
				const wsN = { 
					uri: ws.uri,
					name:  S?.state?.title || "refresh",
				}
				const result = vscode.workspace.updateWorkspaceFolders(i, 1, wsN)
			}
		})
	}

	async refreshFS(sid: string, pos?: string) {
		const wsF = vscode?.workspace?.workspaceFolders as any[]
		const count = vscode.workspace.workspaceFolders?.length || 0

		wsF.forEach((ws, i)=>{
			const qp = new URLSearchParams(wsF[i].uri.query)
			var s = qp.get("sid") as string
			if (s === sid) {
				const S = this._sessions[sid]
				if (S) {
					ws.name = S.state?.title || sid
					const wsN = { 
						uri: ws.uri,
						name:  S?.state?.title || "refresh",
					}
					const result = vscode.workspace.updateWorkspaceFolders(i, 1, wsN)
				}
				return
			}
		})
	}

	async mergeFS(sid: string) {
		console.log("Merging", sid)
		// TODO, get diff so we are for sure updated
		// TODO, accept pos when we do so we can get snapshots
		const diff: any = this.diff[sid]
		for (var key in diff.files) {
			const val = diff.files[key]
			await fs.writeFile(key, val)
		}
	}

	closeFS(sid: string) {
		const wsF = vscode?.workspace?.workspaceFolders as any[]
		var p = -1
		for (var i in wsF){
			// wtf is i not a number...? fucking typescript
			p++
			
			const qp = new URLSearchParams(wsF[i].uri.query)
			var s = qp.get("sid") as string
			// console.log("closeFS.loop:", i, wsF[i].sid === sid, wsF[i], p, qp, s, s === sid)
			if (s === sid) {
				// console.log("CLOSING!", sid)
				vscode.workspace.updateWorkspaceFolders(p, 1)
				return
			}

		}
	}

	provideTextDocumentContent(uri: vscode.Uri, token: vscode.CancellationToken): vscode.ProviderResult<string> {
		const s = this.uriToSid(uri)
		this.ensureExists(s)
		return this.fs[s][uri.path]
	}

	addFile(uri: vscode.Uri, content: string) {
		const s = this.uriToSid(uri)
		this.ensureExists(s)
		const wsF = vscode.workspace.workspaceFolders
		var wDir: string | undefined
		if (wsF && wsF.length > 0) {
			wDir = wsF[0].uri.path	
		}

		this.fs[s][`${wDir}/${uri.path}`] = content
		this.fs[s][uri.path] = content
		delete this.fs[s][uri.path]
	}

	addFileFS(sid: string, path: string, content: string) {
		this.ensureExists(sid)
		const wsF = vscode.workspace.workspaceFolders
		var wDir: string | undefined
		if (wsF && wsF.length > 0) {
			wDir = wsF[0].uri.path	
		}

		this.fs[sid][`${wDir}/${path}`] = content
	}

	setDiff(sid: string, diff: any) {
		this.diff[sid] = diff
	}

	setSessions(sessions: any[]) {
		sessions.forEach((s) => {
			this._sessions[s.sid] = s
		})
	}

	setFS(sid: string, files: Record<string,string>) {
		// console.log("veg.setFS", sid)
		this.fs[sid] = {}
		for (const p in files) {
			// console.log("veg.setFS", p)
			this.addFileFS(sid, p, files[p])
		}
	}

	async showDiff(sid: string, path?: string) {
		for (const p in this.diff[sid].files) {
			const ap = this.diff[sid].addpaths as string[]
			const mp = this.diff[sid].modpaths as string[]
			// see if it was added or modified
			if (ap.includes(p)) {
				// console.log("ADD:", p)
				var vegUri = vscode.Uri.parse(`veg://${p}`)
				vegUri = vegUri.with({
					query: `sid=${sid}`
				})
				vscode.window.showTextDocument(vegUri, {
					preserveFocus: false,
					preview: false,
					viewColumn: 1,
				})
				continue
			} 

			if (mp.includes(p)) {
				// console.log("DIFF:", p)
				const fileUri = vscode.Uri.parse(`file://${p}`)
				var vegUri = vscode.Uri.parse(`veg://${p}`)
				vegUri = vegUri.with({
					query: `sid=${sid}`
				})
				// console.log("veg.diff", fileUri, vegUri)
				vscode.commands.executeCommand('vscode.diff', fileUri, vegUri, `veg-diff: ${p}`, {
					preserveFocus: false,
					preview: false,
					viewColumn: 1,
				})
			}
		}
	}

	// ensures a session at least exists in the "fs"
	// returns true if exists, creates and returns false if not
	private ensureExists(s: string): boolean {
		if (s in this.fs) {
			return true
		}
		this.fs[s] = {}
		return false
	}

	// returns our fs key from the uri (<session>[-<pos>])
	private uriToSid(uri: vscode.Uri): string {
		const qp = new URLSearchParams(uri.query)
		var s = qp.get("sid") as string
		if (!s) {
			console.error("sid unset in vscode.Uri", uri, qp, s)
			return "error"
		}
		// const p = qp.get("pos") as string
		// if (p && p !== '') {
		// 	s += "-" + p
		// }
		return s
	}

}

export class VegFileDecorationProvider implements vscode.FileDecorationProvider {
  
  // 1. Event Emitter: Essential for updating the UI when state changes
  private _onDidChangeFileDecorations: vscode.EventEmitter<vscode.Uri | vscode.Uri[] | undefined> = new vscode.EventEmitter<vscode.Uri | vscode.Uri[] | undefined>();
  readonly onDidChangeFileDecorations: vscode.Event<vscode.Uri | vscode.Uri[] | undefined> = this._onDidChangeFileDecorations.event;

  private _diff: Record<string,any> = {}
	setDiff(sid: string, diff: any) {
		this._diff[sid] = diff

    // todo, broadcast decorations
		const uris: vscode.Uri[] = []
		if (diff.addpaths) {
			for (const p of diff.addpaths) {
				var uri = vscode.Uri.parse(`veg://${p}`)
				uri = uri.with({
					query: `sid=${sid}`
				})
				uris.push(uri)
			}
		}
		if (diff.modpaths) {
			for (const p of diff.modpaths) {
				var uri = vscode.Uri.parse(`veg://${p}`)
				uri = uri.with({
					query: `sid=${sid}`
				})
				uris.push(uri)
			}
		}
		if (diff.delpaths) {
			for (const p of diff.delpaths) {
				var uri = vscode.Uri.parse(`veg://${p}`)
				uri = uri.with({
					query: `sid=${sid}`
				})
				uris.push(uri)
			}
		}
		var uri = vscode.Uri.parse(`veg:///Users/tony/adk/go`)
		uri = uri.with({
			query: `sid=${sid}`
		})
		uris.push(uri)

    console.log("BROADCASTING!!!", diff, uris)
		this._onDidChangeFileDecorations.fire(uris)
	}

  // 2. The Core Logic
  provideFileDecoration(uri: vscode.Uri): vscode.ProviderResult<vscode.FileDecoration> {
    // A. Guard: Only run for your specific scheme
    if (uri.scheme !== 'veg') {
        return undefined;
    }
		const qp = new URLSearchParams(uri.query)
		var s = qp.get("sid") as string
    if (!s) {
      return undefined
    }

    const diff = this._diff[s]
    const path = uri.path
    const hasMods = diff?.addpaths?.length > 0 || diff?.modpaths?.length > 0 || diff?.delpaths?.length > 0
    console.log("VEG.deco", uri, diff, hasMods)

    // if ((!path || path === "" || path === "/") && diff) {
    //   return {
    //     badge: 'M', // 1-2 characters max
    //     color: new vscode.ThemeColor('gitDecoration.modifiedResourceForeground'), // Use theme colors
    //     tooltip: 'New item'
    //   };
    // }

		return matchPathInDiff(path, diff)

    // Return undefined if no decoration is needed (default look)
    return undefined;
  }

  // 3. Method to trigger updates manually
  updateDecorations(uri: vscode.Uri) {
    this._onDidChangeFileDecorations.fire(uri);
  }
}

function matchPathInDiff(path: string, diff: any): any {
	if (diff?.modpaths?.includes(path)) {
		// console.log("found-mp-i", path)
	  return {
	    badge: '🍋', // 1-2 characters max
	    color: new vscode.ThemeColor('gitDecoration.modifiedResourceForeground'), // Use theme colors
	    tooltip: 'Modified'
	  };
	}
	if (diff?.addpaths?.includes(path)) {
		// console.log("found-ap-i", path)
	  return {
	    badge: '🌵', // 1-2 characters max
	    color: new vscode.ThemeColor('gitDecoration.addedResourceForeground'), // Use theme colors
	    tooltip: 'Created'
	  };
	}
	if (diff?.delpaths?.includes(path)) {
		// console.log("found-dp-i", path)
	  return {
	    badge: '🍄', // 1-2 characters max
	    color: new vscode.ThemeColor('gitDecoration.modifiedResourceForeground'), // Use theme colors
	    tooltip: 'Deleted'
	  };
	}

	// look for prefixes, i.e. when a dir
	// modify first, for directories
	for (const p of diff?.modpaths) {
			// console.log("found-mp-s", path)
		if (p.startsWith(path)) {
			return {
				badge: '🍋', // 1-2 characters max
				color: new vscode.ThemeColor('gitDecoration.modifiedResourceForeground'), // Use theme colors
				tooltip: 'Modifed'
			};
		}
	}
	for (const p of diff?.addpaths) {
		if (p.startsWith(path)) {
			// console.log("found-ap-s", path)
			return {
				badge: '🌵', // 1-2 characters max
				color: new vscode.ThemeColor('gitDecoration.addedResourceForeground'), // Use theme colors
				tooltip: 'Created'
			};
		}
	}
	for (const p of diff?.delpaths) {
		if (p.startsWith(path)) {
			// console.log("found-dp-s", path)
			return {
				badge: '🍄', // 1-2 characters max
				color: new vscode.ThemeColor('gitDecoration.modifiedResourceForeground'), // Use theme colors
				tooltip: 'Deleted'
			};
		}
	}

	// console.log("not found", path)
	return undefined;
}