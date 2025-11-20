import * as vscode from 'vscode';

import { extensionEmitter } from '../util/events';
import { sendMessage } from '../websocket';

var agents: any = [];

export function activate(context: vscode.ExtensionContext) {
	const rootPath = (vscode.workspace.workspaceFolders && (vscode.workspace.workspaceFolders.length > 0))
		? vscode.workspace.workspaceFolders[0].uri.fsPath : undefined;

	// Samples of `window.registerTreeDataProvider`
	const planningProvider = new PlanningProvider(context, rootPath);
	vscode.window.registerTreeDataProvider('veg-models', planningProvider);
	vscode.commands.registerCommand('veg.models.refresh', () => planningProvider.refresh());
	vscode.commands.registerCommand('veg.models.chat', (node: Agent) => {
		vscode.window.showInformationMessage(`Successfully called chat entry on ${node.label}.`)
	});
	extensionEmitter.event((e) => {
		switch (e.type) {
			case "models.list":
				console.log("models", e.payload)
				agents = e.payload
				planningProvider.refresh()
				break;
		}
	});

}


export class PlanningProvider implements vscode.TreeDataProvider<Agent> {

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

		if (element) {
			console.log("elemental element", element)
			var nodes: Agent[] = []
			for (const p of element.data.agents) {
				const n = new Agent(p, this.context.extensionUri, vscode.TreeItemCollapsibleState.Collapsed)
				nodes.push(n)
			}
			return Promise.resolve(nodes);
		} else {
			console.log("elementless child", agents)
			if (agents.length < 1) {
				agents = defaultTasks
			}
			var nodes: Agent[] = []
			for (const p of agents) {
				const n = new Agent(p, this.context.extensionUri, vscode.TreeItemCollapsibleState.Collapsed)
				nodes.push(n)
			}
			return Promise.resolve(nodes);
		}
	}

}

export class Agent extends vscode.TreeItem {
	public label?: string = ""
	constructor(
		public readonly data: any,
		public readonly extensionRoot: vscode.Uri,
		public readonly collapsibleState: vscode.TreeItemCollapsibleState,
	) {

		const label = `${data.id}`

		super(data.id, collapsibleState);
		this.label = label

		// this.iconPath = vscode.ThemeIcon.name

		this.tooltip = `${this.id} ${this.data.title}`;
		// this.description = this.version;
	}

	contextValue = 'agent';
}

const defaultTasks: any = [{
	id: "coding",
},{
	id: "planning",
},{
	id: "research",
},{
	id: "general",
}]