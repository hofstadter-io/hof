import * as vscode from 'vscode';

import { extensionEmitter, sendMessage } from '../comms';

// todo, this is probably bad (being global)
var sessions: any = [];
var sort: string = 'lastUpdate'

export function activate(context: vscode.ExtensionContext) {
	const rootPath = (vscode.workspace.workspaceFolders && (vscode.workspace.workspaceFolders.length > 0))
		? vscode.workspace.workspaceFolders[0].uri.fsPath : undefined;

	// Samples of `window.registerTreeDataProvider`
	const sessionsProvider = new EnvironsProvider(context, rootPath);
	vscode.window.registerTreeDataProvider('veg-environs', sessionsProvider);


	// vscode.commands.registerCommand('veg.sessions.refresh', () => sessionsProvider.refresh());

	// vscode.commands.registerCommand('veg.sessions.create', () => {
	// 	sendMessage({ type: "session.create", payload: { focus: true, dir: rootPath } })
	// });
	// vscode.commands.registerCommand('veg.sessions.chat', (node: Environ) => {
	// 	vscode.commands.executeCommand('veg-chat-webview.focus')
	// 	const payload = { sid: node.sid }
	// 	extensionEmitter.fire({ type: "chat.loadSession", payload })
	// 	sendMessage({ type: "session.diff", payload })
	// });
	// vscode.commands.registerCommand('veg.sessions.edit', (node: Environ) => vscode.window.showInformationMessage(`Successfully called edit entry on ${node.label}.`));
	// vscode.commands.registerCommand('veg.sessions.delete', (node: Environ) => {
	// 	const msg = {
	// 		type: "session.delete",
	// 		payload: {
	// 			sid: node.sid
	// 		}
	// 	}
	// 	sendMessage(msg)
	// 	extensionEmitter.fire(msg)
	// 	vscode.window.showInformationMessage(`Successfully called delete entry on ${node.label}.`)
	// });

	// // incoming messages
	// extensionEmitter.event(async (e) => {
	// 	// console.log(`sessions event:`, e)
	// 	switch (e.type) {
	// 		case "session.list":
	// 			// console.log("sessions", e.payload)
	// 			sessions = e.payload
	// 			sessionsProvider.refresh()
	// 			break;
	// 	}
	// });

	// extensionEmitter.fire({
	// 	type: "requestSync"	
	// })
}


export class EnvironsProvider implements vscode.TreeDataProvider<Environ> {

	private _onDidChangeTreeData: vscode.EventEmitter<Environ | undefined | void> = new vscode.EventEmitter<Environ | undefined | void>();
	readonly onDidChangeTreeData: vscode.Event<Environ | undefined | void> = this._onDidChangeTreeData.event;

	constructor(
		private readonly context: vscode.ExtensionContext,
		private readonly workspaceRoot: string | undefined
	) { }

	refresh(): void {
		this._onDidChangeTreeData.fire();
	}

	getTreeItem(element: Environ): vscode.TreeItem {
		return element;
	}

	getChildren(element?: Environ): Thenable<Environ[]> {

		if (element) {
			console.log("environ element", element)
			return Promise.resolve([]);
		} else {
			console.log("environ child", element)
			// console.log("elementless child", sessions)
			var nodes: Environ[] = []
			for (const s of sessions) {
				const l = s.state?.title || s.sid
				const n = new Environ(s.sid, l, s.lastUpdate, vscode.TreeItemCollapsibleState.Collapsed)
				nodes.push(n)
			}
			nodes.sort((a: Environ,b: Environ) => {
				if (sort === "lastUpdate") {
					// newest at the tope
					if(a[sort] > b[sort]) {
						return -1
					}
					if(a[sort] < b[sort]) {
						return 1
					}
					return 0
				}

				if(a.label < b.label) {
					return -1
				}
				if(a.label > b.label) {
					return 1
				}
				return 0

			})
			return Promise.resolve(nodes);
		}
	}

}

export class Environ extends vscode.TreeItem {
	constructor(
		public readonly sid: string,
		public readonly label: string,
		private readonly lastUpdate: string,
		public readonly collapsibleState: vscode.TreeItemCollapsibleState,
		public readonly command?: vscode.Command
	) {

		super(sid, collapsibleState);

		this.tooltip = `${this.label}\n${this.sid}\n${this.lastUpdate}`;
		this.description = this.lastUpdate;

		// this.iconPath = {
		// 	light: vscode.Uri.joinPath(extensionRoot, 'resources', 'light', 'list-tree.svg'),
		// 	dark: vscode.Uri.joinPath(extensionRoot, 'resources', 'dark', 'list-tree.svg')
		// };
	}

	contextValue = 'session';
}