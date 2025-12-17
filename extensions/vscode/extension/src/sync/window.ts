import * as vscode from 'vscode';
import { extensionEmitter, sendMessage } from '../comms';

// This method is called when your extension is activated
// Your extension is activated the very first time the command is executed
export function activate(context: vscode.ExtensionContext) {

	context.subscriptions.push(
		vscode.workspace.onDidOpenTextDocument(e => {
			if (e.uri.scheme !== 'file') {
				return;
			}
			// console.log("openedDocument", e.fileName)
		})
	);

	extensionEmitter.event((e) => {
		switch (e.type) {
      case "sync.request":
      case "sync.request.window":
			case "requestSync":
				broadcastWindow()
				break;
		}
	});

	// context.subscriptions.push(disposable);
}

// This method is called when your extension is deactivated
export function deactivate() {}

function broadcastWindow() {

	const tgs = vscode.window.tabGroups.all.map((tg: vscode.TabGroup) => {

		return {
			active: tg.isActive,
			column: tg.viewColumn,
			tabs: tg.tabs.map((t: vscode.Tab) => {
				return {
					label: t.label,
					isActive: t.isActive,
					isDirty: t.isDirty,
					isPinned: t.isPinned,
					isPreview: t.isPreview,
				}
			})
		}
	})
	const msg = {
		type: "window.info.resp",
		payload: {
			state: vscode.window.state,
			// activeTextEditor: vscode.window.activeTextEditor,
			// activeTerminal: vscode.window.activeTerminal,
			tabgroups: tgs,
			visibleTextEditors: vscode.window.visibleTextEditors,
		}
	}
	extensionEmitter.fire(msg);
	sendMessage(msg)
}