import * as vscode from 'vscode';

import { sendMessage } from '../websocket'
import { extensionEmitter } from '../util/events';
import { WebviewProvider } from './provider'

export async function activate(context: vscode.ExtensionContext) {
  console.log("  activating debug sidebar")
  const provider = new WebviewProvider(context, "debug", (_)=>{});
  context.subscriptions.push(
    vscode.window.registerWebviewViewProvider(
      `veg-debug-webview`, // This ID must match package.json
      provider
    )
  );

  const sync = async () => {
    extensionEmitter.fire({
      type: "requestSync",
    });
    sendMessage({
      type: "requestSync",
      payload: null
    })
  }

	vscode.commands.registerCommand('veg.debug.requestSync', () => {
    sync()
	});
  sync()
}

// onMessage handles  Webview -> Server
// function onMessage(data: any): void {
//   switch (data.type) {
//     case 'dummy':
//       extensionEmitter.fire({
//         type: data.type,
//         payload: data.payload
//       });
//       // The user typed a message. Send it to the websocket.
//       sendMessage({
//         type: data.type,
//         payload: data.payload
//       }); // This function is imported from websocket.ts
//       break;
//   }
// }
