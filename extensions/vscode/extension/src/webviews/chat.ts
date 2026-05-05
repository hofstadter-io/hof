import * as path from 'path';
import * as vscode from 'vscode';

import { extensionEmitter, sendMessage } from '../comms';
import { WebviewProvider } from './provider'
import { makeReq } from '../services/utils';

const prompts = new Map<string, string>();
const promptProvider = new class implements vscode.TextDocumentContentProvider {
	onDidChangeEmitter = new vscode.EventEmitter<vscode.Uri>();
	onDidChange = this.onDidChangeEmitter.event;
	provideTextDocumentContent(uri: vscode.Uri): string {
		return prompts.get(uri.toString()) || "";
	}
};

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

	context.subscriptions.push(vscode.workspace.registerTextDocumentContentProvider('veg-prompt', promptProvider));

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
			const data = await resp.json() as { prompt: string };
			const prompt = data.prompt;

			// Create a filename
			const filename = `sid-${from}${agent ? '-' + agent : ''}.md`;
			const uri = vscode.Uri.parse(`veg-prompt:${filename}`);

			prompts.set(uri.toString(), prompt);
			promptProvider.onDidChangeEmitter.fire(uri);

			const doc = await vscode.workspace.openTextDocument(uri);
			await vscode.window.showTextDocument(doc, { preview: true });
			await vscode.languages.setTextDocumentLanguage(doc, 'markdown');
		} else {
			vscode.window.showErrorMessage(`Failed to render prompt: ${resp.statusText}`);
		}
	} catch (err) {
		vscode.window.showErrorMessage(`Error rendering prompt: ${err}`);
	}
}
