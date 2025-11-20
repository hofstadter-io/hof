import * as vscode from 'vscode';

import { extensionEmitter } from '../util/events';
import { sendMessage } from '../websocket'
import { WebviewProvider } from './provider'

export async function activate(context: vscode.ExtensionContext) {
  console.log("  activating chat sidebar")
  const provider = new WebviewProvider(context, "chat", onMessage);
  context.subscriptions.push(
    vscode.window.registerWebviewViewProvider(
      `veg-chat-webview`, // This ID must match package.json
      provider
    )
  );
}

// onMessage handles messages sent by the Webview
function onMessage(data: any): void {
  extensionEmitter.fire(data)
  sendMessage(data); // This function is imported from websocket.ts
}
