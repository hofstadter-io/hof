import * as vscode from 'vscode';
import { extensionEmitter } from '../util/events';
import { sendMessage } from '../websocket';

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
				broadcastEnv()
				break;
		}
	});

	// context.subscriptions.push(disposable);
}

// This method is called when your extension is deactivated
export function deactivate() {}

async function broadcastEnv() {
	const msg = {
		type: "env.info.resp",
		payload: {
			machineId: vscode.env.machineId,
			sessionId: vscode.env.sessionId,
			remoteName: vscode.env.remoteName,
			appRoot: vscode.env.appRoot,
			clipboard: await vscode.env.clipboard.readText(),
		}
	}
	extensionEmitter.fire(msg);
	sendMessage(msg)
}