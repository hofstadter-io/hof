import * as vscode from 'vscode';
import fs from 'node:fs/promises';
import { makeReq } from './utils';
import { extensionEmitter } from '../comms';

export class VegScmProvider {
	private _sessions: any[] = []
	private _latestEnvs: Map<string, string> = new Map()
	private _scms: Map<string, vscode.SourceControl> = new Map()
	private _resourceGroups: Map<string, Map<string, vscode.SourceControlResourceGroup>> = new Map()

	setSessions(sessions: any[]) {
		this._sessions = sessions
	}

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
			const resp = await makeReq("/fs/diff", uri)
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
}

export const scmProvider = new VegScmProvider()

export async function activate(context: vscode.ExtensionContext) {

	vscode.commands.registerCommand('veg.explorer.showDiff', async (uri: vscode.Uri) => {
		console.log("veg.explorer.showDiff.args.CMD.uri", uri)
		scmProvider.showDiff(uri)
	})

	vscode.commands.registerCommand('veg.explorer.mergeDiff', async (arg: any) => {
		console.log("veg.explorer.mergeDiff.args", arg)
		scmProvider.mergeDiff(arg, undefined, true)
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
				const uriDiff = vscode.Uri.parse("oci://" + e.payload.currEnv)
				scmProvider.showDiff(uriDiff)
				break

			case "session.merge":
				console.log("filesys.session.merge", e.payload)
				const uriMerge = vscode.Uri.parse("oci://" + e.payload.currEnv)
				const dest = e.payload.dest ? vscode.Uri.parse(e.payload.dest) : undefined
				scmProvider.mergeDiff(uriMerge, dest, e.payload.forceInput)
				break
		}
	})

}
