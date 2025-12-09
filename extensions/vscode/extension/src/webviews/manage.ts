import * as vscode from 'vscode';

import { extensionEmitter, sendMessage } from '../comms'
import { WebviewProvider } from './provider'

export async function activate(context: vscode.ExtensionContext) {
  console.log("  activating manage sidebar")
  const provider = new WebviewProvider(context, "manage", onMessage);
  context.subscriptions.push(
    vscode.window.registerWebviewViewProvider(
      `veg-manage`, // This ID must match package.json
      provider
    )
  );

  // incoming messages
	extensionEmitter.event((e) => {
		// console.log(`manage.panel event:`, e)
		switch (e.type) {
		}
	});
}

// onMessage handles  Webview -> Server
function onMessage(data: any): void {
  switch (data.type) {
  }
}