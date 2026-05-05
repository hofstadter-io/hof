import * as vscode from 'vscode';

import { extensionEmitter, sendMessage } from '../comms';

// todo, this is probably bad (being global)
var sessions: Map<string, any> = new Map();
var sort: string = 'lastUpdate'

export function activate(context: vscode.ExtensionContext) {

	// Samples of `window.registerTreeDataProvider`
	const sessionsProvider = new SessionsProvider(context);
	vscode.window.registerTreeDataProvider('veg-sessions', sessionsProvider);

	vscode.commands.registerCommand('veg.sessions.refresh', () => sessionsProvider.refresh());

	vscode.commands.registerCommand('veg.sessions.create', (node?: Session) => {
		console.log("create session", node)
		sendMessage({ type: "session.create", payload: { focus: true } })
	});

	vscode.commands.registerCommand('veg.session.openEnviron', (node?: Session) => {
		console.log("open session", node)
		if (!node || !node.session) {
			console.error("veg.session.openEnviron was called without input or session info")
			return
		}
		// sendMessage({ type: "session.create", payload: { focus: true, dir: rootPath } })
		const msg = {
			type: "filesys.openEnviron",
			payload: {
				session: node.session
			}
		}
		extensionEmitter.fire(msg)
	});

	vscode.commands.registerCommand('veg.session.chat', (node: Session) => {
		vscode.commands.executeCommand('veg-chat.focus')
		const payload = { sid: node.session.sid }
		extensionEmitter.fire({ type: "chat.loadSession", payload })
		sendMessage({ type: "session.diff", payload })
	});

	vscode.commands.registerCommand('veg.session.terminal', (node: Session) => {
		extensionEmitter.fire({
			type: "session.term.open",
			payload: {
				sid: node.session.sid,
				image: node.session.state?.currEnv,
			}
		})
	});

	vscode.commands.registerCommand('veg.session.showDiff', (node: Session) => {
		extensionEmitter.fire({
			type: "session.diff",
			payload: {
				sid: node.session.sid,
				show: true,
				currEnv: node.session.state?.currEnv,
			}
		})
	});

	vscode.commands.registerCommand('veg.session.mergeDiff', (node: Session) => {
		extensionEmitter.fire({
			type: "session.merge",
			payload: {
				sid: node.session.sid,
				currEnv: node.session.state?.currEnv,
			}
		})
	});

	vscode.commands.registerCommand('veg.session.clone', (node: Session) => {
		extensionEmitter.fire({
			type: "session.clone",
			payload: {
				sid: node.session.sid,
				pos: node.session.events?.length || 0,
				focus: true,
			}
		})
	});

	vscode.commands.registerCommand('veg.session.fork', (node: Session) => {
		console.log("fork session", node)

	})

	vscode.commands.registerCommand('veg.session.edit', (node: Session) => {
		vscode.window.showInformationMessage(`Successfully called edit entry on ${node.label}.`)
	});
	vscode.commands.registerCommand('veg.session.delete', (node: Session) => {
		const msg = {
			type: "session.delete",
			payload: {
				sid: node.session.sid
			}
		}
		sendMessage(msg)
		extensionEmitter.fire(msg)
		vscode.window.showInformationMessage(`Successfully called delete entry on ${node.label}.`)
	});

	// incoming messages
	extensionEmitter.event(async (e) => {
		// console.log(`sessions event:`, e)
		switch (e.type) {
			case "session.list.resp":
				// console.log("sessions", e.payload)
				sessions.clear()
				for (const s of e.payload) {
					sessions.set(s.sid, s)
				}
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
	) { }

	refresh(): void {
		this._onDidChangeTreeData.fire();
	}

	getTreeItem(element: Session): vscode.TreeItem {
		return element;
	}

	getChildren(element?: Session): Thenable<Session[] | undefined> {

		if (element) {
			// console.log("elemental element", element)
			// return Promise.resolve([]);
			return Promise.resolve(undefined)
		} else {
			// root, so we work with the sessions we know about
			// console.log("elementless child", sessions)
			var nodes: Session[] = []
			for (const s of sessions.values()) {
				const l = s.state?.title || s.sid
				const n = new Session(s, l, vscode.TreeItemCollapsibleState.Collapsed)
				nodes.push(n)
			}
			nodes.sort((a: Session,b: Session) => {
				if (sort === "lastUpdate") {
					// newest at the tope
					if(a.session[sort] > b.session[sort]) {
						return -1
					}
					if(a.session[sort] < b.session[sort]) {
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
			if (nodes.length === 0) {
				return Promise.resolve(undefined)
			}
			return Promise.resolve(nodes);
		}
	}

}

export class Session extends vscode.TreeItem {
	constructor(
		public readonly session: any,
		public readonly label: string,
		public readonly collapsibleState: vscode.TreeItemCollapsibleState,
		public readonly command?: vscode.Command,
	) {
		super(session.sid, collapsibleState);

		this.tooltip = `${this.label}\n${this.session.sid}\n${this.session.lastUpdate}`;
		this.description = this.session.lastUpdate;

		// this.iconPath = new vscode.ThemeIcon("list-tree")
	}

	contextValue = 'session';
}
