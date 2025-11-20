import * as vscode from 'vscode';

import { sendMessage } from '../websocket'
import { WebviewProvider } from './provider'

export async function activate(context: vscode.ExtensionContext) {
  console.log("  activating planner webview")
  const provider = new WebviewProvider(context, "planning", onMessage);
  context.subscriptions.push(
    vscode.window.registerWebviewViewProvider(
      `veg-planning-webview`, // This ID must match package.json
      provider
    )
  );
}

// onMessage handles  Webview -> Server
function onMessage(data: any): void {
  switch (data.type) {
    case 'sendConfig':
      // The user typed a message. Send it to the websocket.
      sendMessage({
        type: data.type,
        payload: data.text
      }); // This function is imported from websocket.ts
      break;
  }
}
