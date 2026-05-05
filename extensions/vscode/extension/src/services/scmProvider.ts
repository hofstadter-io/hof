import * as vscode from 'vscode';
import { makeReq, parseEnvUri, findSession } from './utils';
import { extensionEmitter } from '../comms';

export async function activate(context: vscode.ExtensionContext) {

	vscode.commands.registerCommand('veg.explorer.showDiff', async (uri: vscode.Uri) => {
		console.log("veg.explorer.showDiff.args.CMD.uri", uri)
		scmProvider.showDiff(uri)
	})

	vscode.commands.registerCommand('veg.explorer.mergeDiff', async (arg: any) => {
		console.log("veg.explorer.mergeDiff.args", arg)
		try {
			await scmProvider.mergeDiff(arg, undefined, false)
		} catch (e) {
			console.error("veg.explorer.mergeDiff execution error:", e)
		}
	})

	vscode.commands.registerCommand('veg.explorer.hideDiff', async (arg: any) => {
		console.log("veg.explorer.hideDiff.args", arg)
		scmProvider.hideDiff(arg)
	})

	vscode.commands.registerCommand('veg.explorer.diffAll', async (args: any) => {
		console.log("veg.explorer.diffAll.args", args)
		vscode.window.showInformationMessage("Diff All (not implemented yet)")
	})

	extensionEmitter.event(async (e) => {
		// ...
		switch (e.type) {

			case "session.list.resp":
				scmProvider.setSessions(e.payload)
				break

			case "session.diff":
				console.log("filesys.session.diff", e.payload)
				let currEnv = e.payload.currEnv;
				if (!currEnv && e.payload.sid && e.payload.pos !== undefined) {
					const session = findSession(scmProvider['_sessions'], e.payload.sid) || scmProvider['_sessions'].find((s: any) => s.sid === e.payload.sid);
					currEnv = scmProvider.findCurrEnv(session, e.payload.pos);
				}

				if (currEnv) {
					const uriDiff = vscode.Uri.parse("oci://" + currEnv)
					scmProvider.showDiff(uriDiff)
				}
				break

			case "session.merge":
				console.log("filesys.session.merge", e.payload)
				let mergeEnv = e.payload.currEnv;
				if (!mergeEnv && e.payload.sid && e.payload.pos !== undefined) {
					const session = findSession(scmProvider['_sessions'], e.payload.sid) || scmProvider['_sessions'].find((s: any) => s.sid === e.payload.sid);
					mergeEnv = scmProvider.findCurrEnv(session, e.payload.pos);
				}

				if (mergeEnv) {
					const uriMerge = vscode.Uri.parse("oci://" + mergeEnv)
					const dest = e.payload.dest ? vscode.Uri.parse(e.payload.dest) : undefined
					scmProvider.mergeDiff(uriMerge, dest, e.payload.forceInput)
				}
				break
		}
	})

}
export class VegScmProvider {
	private _sessions: any[] = []
	private _latestEnvs: Map<string, string> = new Map()
	private _scms: Map<string, vscode.SourceControl> = new Map()
	private _resourceGroups: Map<string, Map<string, vscode.SourceControlResourceGroup>> = new Map()


	setSessions(sessions: any[]) {
		this._sessions = sessions
	}

	private getScmInfo(source: any): { uri?: vscode.Uri, session?: any, scmId?: string, groupId?: string } {
		console.log("getScmInfo: start", { 
			hasSource: !!source, 
			id: source?.id, 
			scheme: source?.scheme,
			isUri: source instanceof vscode.Uri
		});

		if (!source) return {};

		if (source instanceof vscode.Uri || source.scheme) {
			return { uri: source as vscode.Uri };
		}

		const id = source.id;
		if (!id) {
			console.log("getScmInfo: no id found on source");
			return {};
		}

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
				console.log("getScmInfo: found in _scms", { scmId, uri: uri?.toString() });
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
						console.log("getScmInfo: found in _resourceGroups", { scmId, groupId, uri: uri?.toString() });
						break;
					}
				}
				if (scmId) break;
			}
		}

		if (!scmId) {
			// Fallback to old logic
			session = this._sessions.find(s => s.state?.currEnv === id || s.sid === id);
			if (session) {
				scmId = session.sid;
				const sEnv = session.state?.currEnv;
				uri = vscode.Uri.parse("veg://" + (sEnv || id));
				console.log("getScmInfo: found via session fallback", { scmId, uri: uri.toString() });
			} else {
				scmId = id;
				// avoid double scheme if id already has one
				const uriStr = id.includes("://") ? id : "veg://" + id;
				uri = vscode.Uri.parse(uriStr);
				console.log("getScmInfo: final fallback", { scmId, uri: uri.toString() });
			}
		}

		return { uri, session, scmId, groupId };
	}

	private isIgnoredPath(p: string): boolean {
		return (p.startsWith("/") && p.endsWith("/")) ||
			p.includes("/.git/") || 
			p === "/stdout.txt" || p === "/stderr.txt";
	}

	private processDiffEntry(
		type: 'modified' | 'added' | 'deleted',
		path: string,
		prevBaseUri: vscode.Uri,
		nextBaseUri: vscode.Uri,
		resources: vscode.SourceControlResourceState[],
		multiDiffResources: { originalUri: vscode.Uri | undefined; modifiedUri: vscode.Uri | undefined }[]
	) {
		if (type !== 'modified' && this.isIgnoredPath(path)) {
			return;
		}

		const pfUri = vscode.Uri.from({ ...prevBaseUri, path: prevBaseUri.path + path })
		const nfUri = vscode.Uri.from({ ...nextBaseUri, path: nextBaseUri.path + path })

		let resourceUri: vscode.Uri;
		let command: vscode.Command;
		let icon: vscode.ThemeIcon;
		let tooltip: string;
		let strikeThrough: boolean | undefined;
		let faded: boolean | undefined;

		switch (type) {
			case 'modified':
				resourceUri = nfUri;
				command = {
					command: 'vscode.diff',
					title: 'Show Diff',
					arguments: [pfUri, nfUri, `veg-diff: ${path}`]
				};
				icon = new vscode.ThemeIcon('diff-modified', new vscode.ThemeColor('gitDecoration.modifiedResourceForeground'));
				tooltip = `Modified: ${path}`;
				multiDiffResources.push({ originalUri: pfUri, modifiedUri: nfUri });
				break;

			case 'added':
				resourceUri = nfUri;
				command = {
					command: 'vscode.open',
					title: 'Open File',
					arguments: [nfUri]
				};
				icon = new vscode.ThemeIcon('diff-added', new vscode.ThemeColor('gitDecoration.addedResourceForeground'));
				tooltip = `Added: ${path}`;
				multiDiffResources.push({ originalUri: undefined, modifiedUri: nfUri });
				break;

			case 'deleted':
				resourceUri = pfUri;
				command = {
					command: 'vscode.open',
					title: 'Open File',
					arguments: [pfUri]
				};
				icon = new vscode.ThemeIcon('diff-removed', new vscode.ThemeColor('gitDecoration.deletedResourceForeground'));
				tooltip = `Deleted: ${path}`;
				strikeThrough = true;
				faded = true;
				multiDiffResources.push({ originalUri: pfUri, modifiedUri: undefined });
				break;
		}

		const decorations: vscode.SourceControlResourceDecorations = {
			iconPath: icon!,
			tooltip,
			strikeThrough,
			faded
		};

		resources.push({
			resourceUri,
			contextValue: type,
			decorations,
			command
		});
	}

	public findCurrEnv(session: any, pos?: number): string | undefined {
		let currEnv = session?.state?.currEnv;
		if (pos !== undefined && session?.events) {
			for (let i = pos; i >= 0; i--) {
				const event = session.events[i];
				if (event?.Actions?.StateDelta?.currEnv) {
					currEnv = event.Actions.StateDelta.currEnv;
					break;
				}
			}
		}
		return currEnv;
	}

	// todo, we probably need a diffUri here
	showDiff(source: vscode.Uri | vscode.SourceControlResourceGroup, destination?: vscode.Uri): void | Thenable<void> {
		const f = async () => {
			console.log("filesys.showDiff.ARGS", source, destination)

			const info = this.getScmInfo(source);
			let uri = info.uri;
			let session = info.session;

			if (!uri) {
				return
			}

			let { envId, envVer } = parseEnvUri(uri);

			if (!session) {
				session = findSession(this._sessions, envId, uri)
			}

			const scmId = info.scmId || session?.sid || envId
			const scmTitle = session?.state?.title || scmId
			let groupId = info.groupId || (envId + (envVer !== "?" ? ":" + envVer : ""))
			const groupTitle = `${envVer} : ${envId}`

			// Track latest
			const currentLatest = this._latestEnvs.get(scmId);
			const currentVer = parseInt(envVer);
			if (!currentLatest || (!isNaN(currentVer) && currentVer > parseInt(currentLatest.split(":")[1]))) {
				this._latestEnvs.set(scmId, `${envId}:${envVer}`);
			}

			const resp = await makeReq("/fs/diff", uri)
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

			// Ensure the current group is at the top (newest first)
			// // @ts-ignore
			// const otherGroups = scm.resourceGroups.filter(g => g !== group)
			// // @ts-ignore
			// scm.resourceGroups = [group, ...otherGroups]

			const multiDiffResources: { originalUri: vscode.Uri | undefined; modifiedUri: vscode.Uri | undefined }[] = [];
			const resources: vscode.SourceControlResourceState[] = []

			for (const p of diff.modPaths) {
				this.processDiffEntry('modified', p, prevUri, nextUri, resources, multiDiffResources);
			}
			for (const p of diff.addPaths) {
				this.processDiffEntry('added', p, prevUri, nextUri, resources, multiDiffResources);
			}
			for (const p of diff.delPaths) {
				this.processDiffEntry('deleted', p, prevUri, nextUri, resources, multiDiffResources);
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
			let scmId = info.scmId || "";

			if (!uri) {
				return
			}

			let { envId, envVer } = parseEnvUri(uri);

			if (!session) {
				session = findSession(this._sessions, envId, uri)
			}

			console.log("mergeDiff: session info", { 
				sid: session?.sid, 
				hasInitEnv: !!session?.state?.initEnv,
				srcUri: session?.state?.initEnv?.srcUri 
			})

			if (!scmId) scmId = session?.sid || envId

			// Track latest
			const currentLatest = this._latestEnvs.get(scmId);
			const currentVer = parseInt(envVer);
			if (!currentLatest || (!isNaN(currentVer) && currentVer > parseInt(currentLatest.split(":")[1]))) {
				this._latestEnvs.set(scmId, `${envId}:${envVer}`);
			}

			if (!dest && !forceInput) {
				if (session?.state?.initEnv?.srcUri?.startsWith("file://")) {
					dest = vscode.Uri.parse(session.state.initEnv.srcUri)
					console.log("mergeDiff: using session initEnv srcUri", dest.toString())
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
			} else {
				vscode.window.showInformationMessage(`Merging into ${dest.fsPath}`)
			}

			// if (dest.scheme !== 'file') {
			// 	vscode.window.showErrorMessage(`unsupported target, only file://: ${dest}`)
			// 	return
			// }

			console.log("filesys.mergeDiff.args", source, dest)
			const resp = await makeReq("/fs/diff", uri)
			if (resp.status !== 200) {
				// console.error("filesys.mergeDiff.makeReq error:", resp)
				throw vscode.FileSystemError.FileNotFound(uri)
			}
			// console.log("filesys.mergeDiff.resp", uri, resp)

			const diff: any = await resp.json()
			console.log("filesys.mergeDiff.diff", diff)

			// write to disk or veg
			const writes = [...diff.addPaths, ...diff.modPaths];
			for (var w of writes) {
				if (this.isIgnoredPath(w)) continue;
				// if not ignored, lets trim the prefix / that the server has to add for shenanigans elsewhere in vscode
				// this also makes sure we can never write to the root unless the user gives that path, only below that directory
				const npath = w.substring(1)
				const val = diff.files[npath]
				// const key = path.join(dest.path, w)
				const destUri = vscode.Uri.joinPath(dest, npath)
				// await fs.writeFile(key, val)
				await vscode.workspace.fs.writeFile(destUri, new TextEncoder().encode(val))
			}

			for (const path of diff.delPaths) {
				if (this.isIgnoredPath(path)) continue;
				// const key = path.join(dest.path, w)
				const npath = path.substring(1)
				const destUri = vscode.Uri.joinPath(dest, npath)
				// await fs.rm(key, { force: true, recursive: true }).catch(() => { })
				await vscode.workspace.fs.delete(destUri, { recursive: true })
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

				const session = findSession(this._sessions, envIdFromId)

				if (session) {
					scmId = session.sid
					if (!isScm) groupId = id
				} else {
					scmId = envIdFromId
					if (!isScm) groupId = id
				}
			} else {
				const uri = source as vscode.Uri
				const { envId, envVer } = parseEnvUri(uri)
				const session = findSession(this._sessions, envId, uri)

				scmId = session?.sid || envId
				groupId = envId + (envVer !== "?" ? ":" + envVer : "")
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

	public openTerminal(source: vscode.Uri | vscode.SourceControlResourceGroup): void {
		const info = this.getScmInfo(source);
		let uri = info.uri;
		let session = info.session;

		if (!uri) return;

		let { envId } = parseEnvUri(uri);

		if (!session) {
			session = findSession(this._sessions, envId, uri)
		}

		if (session) {
			extensionEmitter.fire({
				type: "session.term.open",
				payload: {
					sid: session.sid,
					image: session.state?.currEnv
				}
			})
		}
	}
}

export const scmProvider = new VegScmProvider()