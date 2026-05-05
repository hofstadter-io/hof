import * as vscode from 'vscode';

import { extensionEmitter, sendMessage } from '../comms';

var tasks: any = [];

export function activate(context: vscode.ExtensionContext) {
	const rootPath = (vscode.workspace.workspaceFolders && (vscode.workspace.workspaceFolders.length > 0))
		? vscode.workspace.workspaceFolders[0].uri.fsPath : undefined;

	// Samples of `window.registerTreeDataProvider`
	const planningProvider = new PlanningProvider(context, rootPath);
	vscode.window.registerTreeDataProvider('veg-planning', planningProvider);
	vscode.commands.registerCommand('veg.planning.refresh', () => planningProvider.refresh());
	vscode.commands.registerCommand('veg.planning.create', () => {
		vscode.window.showInformationMessage(`Successfully called add entry.`)
		sendMessage({
			type: "planning.create",
			payload: {
				focus: true
			}
		})
	});
	vscode.commands.registerCommand('veg.planning.createSub', (node: Task) => {
		vscode.window.showInformationMessage(`Successfully called createSub entry on ${node.label}.`)
	});
	vscode.commands.registerCommand('veg.planning.status', (node: Task) => {
		vscode.window.showInformationMessage(`Successfully called status entry on ${node.label}.`)
	});
	vscode.commands.registerCommand('veg.planning.edit', (node: Task) => vscode.window.showInformationMessage(`Successfully called edit entry on ${node.label}.`));
	vscode.commands.registerCommand('veg.planning.delete', (node: Task) => {
		sendMessage({
			type: "planning.delete",
			payload: {
				id: node.id
			}
		})
		vscode.window.showInformationMessage(`Successfully called delete entry on ${node.label}.`)
	});

	extensionEmitter.event((e) => {
		switch (e.type) {
			case "planning.list":
				console.log("tasks", e.payload)
				tasks = e.payload
				planningProvider.refresh()
				break;
		}
	});

}


export class PlanningProvider implements vscode.TreeDataProvider<Task> {

	private _onDidChangeTreeData: vscode.EventEmitter<Task | undefined | void> = new vscode.EventEmitter<Task | undefined | void>();
	readonly onDidChangeTreeData: vscode.Event<Task | undefined | void> = this._onDidChangeTreeData.event;

	constructor(
		private readonly context: vscode.ExtensionContext,
		private readonly workspaceRoot: string | undefined
	) { }

	refresh(): void {
		this._onDidChangeTreeData.fire();
	}

	getTreeItem(element: Task): vscode.TreeItem {
		return element;
	}

	getChildren(element?: Task): Thenable<Task[]> {

		if (element) {
			console.log("elemental element", element)
			var nodes: Task[] = []
			for (const p of element.data.tasks) {
				const n = new Task(p, this.context.extensionUri, vscode.TreeItemCollapsibleState.Collapsed)
				nodes.push(n)
			}
			return Promise.resolve(nodes);
			return Promise.resolve([]);
		} else {
			console.log("elementless child", tasks)
			if (tasks.length < 1) {
				tasks = defaultTasks
			}
			var nodes: Task[] = []
			for (const p of tasks) {
				const n = new Task(p, this.context.extensionUri, vscode.TreeItemCollapsibleState.Collapsed)
				nodes.push(n)
			}
			return Promise.resolve(nodes);
		}
	}

}

export class Task extends vscode.TreeItem {
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

	contextValue = 'task';
}

const defaultTasks: any = [{
	id: "task1",
	status: "inprogress",
	tasks: [{
		id: "task1.1",
		status: "done",
		tasks: [{
			id: "task1.1.1",
			status: "done",
		},{
			id: "task1.1.2",
			status: "done",
		}]
	},{
		id: "task1.2",
		status: "inprogress",
		tasks: [{
			id: "task1.2.1",
			status: "done",
		},{
			id: "task1.2.2",
			status: "inprogress",
		},{
			id: "task1.2.3",
			status: "open",
		}]
	}]
},{
	id: "task2",
	status: "inprogress",
	tasks: [{
		id: "task2.1",
		status: "inprogress",
		tasks: [{
			id: "task2.1.1",
			status: "open",
		},{
			id: "task2.1.2",
			status: "open",
		}]
	},{
		id: "task2.2",
		status: "open",
		tasks: [{
			id: "task2.2.1",
			status: "open",
		},{
			id: "task2.2.2",
			status: "open",
		}]
	}]
}]