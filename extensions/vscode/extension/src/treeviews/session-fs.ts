import * as vscode from 'vscode';

import { extensionEmitter, sendMessage } from '../comms';

var agents: any = [];

export function activate(context: vscode.ExtensionContext) {
	const rootPath = (vscode.workspace.workspaceFolders && (vscode.workspace.workspaceFolders.length > 0))
		? vscode.workspace.workspaceFolders[0].uri.fsPath : undefined;

	// Samples of `window.registerTreeDataProvider`
	const agentProvider = new AgentProvider(context, rootPath);
	vscode.window.registerTreeDataProvider('veg-session-fs', agentProvider);
	vscode.commands.registerCommand('veg.agents.refresh', () => agentProvider.refresh());
	vscode.commands.registerCommand('veg.agents.create', () => {
		vscode.window.showInformationMessage(`Successfully called add entry.`)
		sendMessage({
			type: "agents.create",
			payload: {
				focus: true
			}
		})
	});
	vscode.commands.registerCommand('veg.agents.chat', (node: Agent) => {
		vscode.window.showInformationMessage(`Successfully called chat entry on ${node.label}.`)
	});
	vscode.commands.registerCommand('veg.agents.edit', (node: Agent) => vscode.window.showInformationMessage(`Successfully called edit entry on ${node.label}.`));
	vscode.commands.registerCommand('veg.agents.delete', (node: Agent) => {
		sendMessage({
			type: "agents.delete",
			payload: {
				id: node.id
			}
		})
		vscode.window.showInformationMessage(`Successfully called delete entry on ${node.label}.`)
	});

	extensionEmitter.event((e) => {
		switch (e.type) {
			case "agents.list.resp":
				agents = e.payload
				agentProvider.refresh()
				break;
		}
	});

}


export class AgentProvider implements vscode.TreeDataProvider<Agent> {

	private _onDidChangeTreeData: vscode.EventEmitter<Agent | undefined | void> = new vscode.EventEmitter<Agent | undefined | void>();
	readonly onDidChangeTreeData: vscode.Event<Agent | undefined | void> = this._onDidChangeTreeData.event;

	constructor(
		private readonly context: vscode.ExtensionContext,
		private readonly workspaceRoot: string | undefined
	) { }

	refresh(): void {
		this._onDidChangeTreeData.fire();
	}

	getTreeItem(element: Agent): vscode.TreeItem {
		return element;
	}

	getChildren(element?: Agent): Thenable<Agent[]> {

		console.log("agents.getChildren element", element)

		// if no element, root node, so return our agent list
		if (!element) {
			var nodes: Agent[] = []
			for (const p in agents) {
				const n = new Agent(p, this.context.extensionUri, vscode.TreeItemCollapsibleState.Collapsed)
				nodes.push(n)
			}
			return Promise.resolve(nodes);
		}

		// otherwise, figure out what is in this node and make approprate node types
		// var node = element.data
		// var nodes: Agent[] = []
		// for (const p in node.agents) {
		// 	const n = new Agent(p, this.context.extensionUri, vscode.TreeItemCollapsibleState.Collapsed)
		// 	nodes.push(n)
		// }
		// return Promise.resolve(nodes);

		// return no children
		return Promise.resolve([]);
	}

}

export class Agent extends vscode.TreeItem {
	// public label?: string = ""
	constructor(
		public readonly label: any,
		// public readonly data: any,
		public readonly extensionRoot: vscode.Uri,
		public readonly collapsibleState: vscode.TreeItemCollapsibleState,
	) {

		super(label, collapsibleState);

		// this.iconPath = vscode.ThemeIcon.name

		// this.tooltip = `${this.id} ${this.data.title}`;
		// this.description = this.version;
	}

	contextValue = 'agent';
}