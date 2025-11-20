import * as vscode from 'vscode';

import { extensionEmitter } from '../util/events';
import { sendMessage } from '../websocket';

// todo, this is probably bad (being global)
var sessions: any = [];

export function activate(context: vscode.ExtensionContext) {
	const rootPath = (vscode.workspace.workspaceFolders && (vscode.workspace.workspaceFolders.length > 0))
		? vscode.workspace.workspaceFolders[0].uri.fsPath : undefined;

	// Samples of `window.registerTreeDataProvider`
	const sessionsProvider = new SessionsProvider(context, rootPath);
	vscode.window.registerTreeDataProvider('veg-sessions', sessionsProvider);
	vscode.commands.registerCommand('veg.sessions.refresh', () => sessionsProvider.refresh());

	vscode.commands.registerCommand('veg.sessions.create', () => {
		sendMessage({ type: "session.create", payload: { focus: true } })
	});
	vscode.commands.registerCommand('veg.sessions.chat', (node: Session) => {
		vscode.commands.executeCommand('veg-chat-webview.focus')
		const payload = { id: node.id }
		setTimeout(() => {
			extensionEmitter.fire({ type: "chat.loadSession", payload })
			sendMessage({ type: "session.get", payload })
		}, 1000)
	});
	vscode.commands.registerCommand('veg.sessions.edit', (node: Session) => vscode.window.showInformationMessage(`Successfully called edit entry on ${node.label}.`));
	vscode.commands.registerCommand('veg.sessions.delete', (node: Session) => {
		sendMessage({
			type: "session.delete",
			payload: {
				id: node.id
			}
		})
		vscode.window.showInformationMessage(`Successfully called delete entry on ${node.label}.`)
	});

	// incoming messages
	extensionEmitter.event((e) => {
		// console.log(`sessions event:`, e)
		switch (e.type) {
			case "session.list":
				console.log("sessions", e.payload)
				sessions = e.payload
				sessionsProvider.refresh()
				break;
		}
	});

	extensionEmitter.fire({
		type: "requestSync"	
	})
}


export class SessionsProvider implements vscode.TreeDataProvider<Session> {

	private _onDidChangeTreeData: vscode.EventEmitter<Session | undefined | void> = new vscode.EventEmitter<Session | undefined | void>();
	readonly onDidChangeTreeData: vscode.Event<Session | undefined | void> = this._onDidChangeTreeData.event;

	constructor(
		private readonly context: vscode.ExtensionContext,
		private readonly workspaceRoot: string | undefined
	) { }

	refresh(): void {
		this._onDidChangeTreeData.fire();
	}

	getTreeItem(element: Session): vscode.TreeItem {
		return element;
	}

	getChildren(element?: Session): Thenable<Session[]> {

		if (element) {
			// console.log("elemental element", element)
			return Promise.resolve([]);
		} else {
			// console.log("elementless child", sessions)
			var nodes: Session[] = []
			for (const s of sessions) {
				const n = new Session(s.id, s.id, s.lastUpdate, vscode.TreeItemCollapsibleState.Collapsed)
				nodes.push(n)
			}
			return Promise.resolve(nodes);
		}
	}

}

export class Session extends vscode.TreeItem {
	constructor(
		public readonly id: string,
		public readonly label: string,
		private readonly version: string,
		public readonly collapsibleState: vscode.TreeItemCollapsibleState,
		public readonly command?: vscode.Command
	) {

		super(id, collapsibleState);

		this.tooltip = `${this.label}\n${this.version}`;
		this.description = this.version;

		// this.iconPath = {
		// 	light: vscode.Uri.joinPath(extensionRoot, 'resources', 'light', 'list-tree.svg'),
		// 	dark: vscode.Uri.joinPath(extensionRoot, 'resources', 'dark', 'list-tree.svg')
		// };
	}

	contextValue = 'session';
}