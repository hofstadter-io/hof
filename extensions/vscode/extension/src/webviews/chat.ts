import * as path from 'path';
import * as vscode from 'vscode';

import { extensionEmitter } from '../util/events';
import { sendMessage } from '../websocket'
import { WebviewProvider } from './provider'
import { getHtmlForWebview } from '../webview';

export async function activate(context: vscode.ExtensionContext) {
  console.log("  activating chat sidebar")
  const provider = new WebviewProvider(context, "chat", onMessage);
  context.subscriptions.push(
    vscode.window.registerWebviewViewProvider(
      `veg-chat-webview`, // This ID must match package.json
      provider,
      {
        webviewOptions: { retainContextWhenHidden: true }
      }
    )
  );

  // incoming messages
	extensionEmitter.event((e) => {
		switch (e.type) {
			case "chat.loadSession":
        context.workspaceState.update("sid", e.payload?.sid)
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
