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
      case "sync.request.workspace":
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
	const msg = {
		type: "workspace.info.resp",
		payload: {
			name: vscode.workspace.name,
			isTrusted: vscode.workspace.isTrusted,
			folders: vscode.workspace.workspaceFolders,
			docs: vscode.workspace.textDocuments,
		}
	}
	extensionEmitter.fire(msg);
	sendMessage(msg)
}