import * as path from 'path';
import * as vscode from 'vscode';

import { extensionEmitter, sendMessage } from '../comms';
import { WebviewProvider } from './provider'
import { makeReq } from '../services/utils';

export async function activate(context: vscode.ExtensionContext) {
  console.log("  activating chat sidebar")
  const provider = new WebviewProvider(context, "chat", onMessage);
  context.subscriptions.push(
    vscode.window.registerWebviewViewProvider(
      `veg-chat`, // This ID must match package.json
      provider,
      {
        webviewOptions: { retainContextWhenHidden: true }
      }
    )
  );

  // incoming messages
	extensionEmitter.event(async (e) => {
		switch (e.type) {
			case "chat.loadSession":
        context.workspaceState.update("sid", e.payload?.sid)
				break;
			case "session.prompt":
				await handleSessionPrompt(e.payload);
				break;
		}
	});

  // not working yet
	// vscode.commands.registerCommand('veg.chat.newTab', async () => {
  //   console.log("NEW CHAT")
  //   const panel = vscode.window.createWebviewPanel(
  //     'veg-chat-webview',
  //     'Veg Chat',
  //     vscode.ViewColumn.One,
  //     {
  //       // retainContextWhenHidden: true
  //       enableScripts: true,
  //       // // // Restrict the webview to only loading content from our extension's directories
  //       localResourceRoots: [
  //         vscode.Uri.file(path.join(context.extensionPath, 'webviews', "chat")),
  //       ],

  //     }
  //   )
  //   console.log("WEBVIEW:", panel)

  //   //     // Set the HTML content
  //   // panel.webview.html = await getHtmlForWebview(panel.webview, context.extensionPath, "chat");

  //   // // --- Communication ---

  //   // // 1. Listen for messages from the webview (Chat UI -> Extension)
  //   // panel.webview.onDidReceiveMessage(onMessage);

  //   // // 2. Listen for messages from the websocket (Extension -> Chat UI)
  //   // extensionEmitter.event((e) => {
  //   //   // console.log(`${this._name} event:`, e)
  //   //   // if (e.type === 'WEBSOCKET_MESSAGE') {
  //   //     // We got a message. Pass it to the webview.
  //   //     panel.webview.postMessage(e);
  //   //   // }
  //   // });

  //   // context.subscriptions.push(panel)
	// });
}

// onMessage handles messages sent by the Webview
function onMessage(data: any): void {
  extensionEmitter.fire(data)
  sendMessage(data); // This function is imported from websocket.ts
}

async function handleSessionPrompt(payload: any) {
	const { from, pos, agent, model, environ } = payload;
	try {
		const resp = await makeReq("/prompt/render", undefined, undefined, {
			sid: from,
			pos: pos,
			agent: agent,
			model: model,
			environ: environ,
		});
		if (resp.ok) {
			const data = await resp.json();
			const prompt = data.prompt;

			// Create a filename
			const filename = `sid-${from}${agent ? '-' + agent : ''}.md`;
			const uri = vscode.Uri.parse(`untitled:${filename}`);

			const doc = await vscode.workspace.openTextDocument(uri);
			const edit = new vscode.WorkspaceEdit();
			edit.insert(uri, new vscode.Position(0, 0), prompt);
			await vscode.workspace.applyEdit(edit);
			await vscode.window.showTextDocument(doc, { preview: true });
		} else {
			vscode.window.showErrorMessage(`Failed to render prompt: ${resp.statusText}`);
		}
	} catch (err) {
		vscode.window.showErrorMessage(`Error rendering prompt: ${err}`);
	}
}
