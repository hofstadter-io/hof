import * as vscode from 'vscode';

import { extensionEmitter, sendMessage } from '../comms';
import { WebviewProvider } from './provider'

export async function activate(context: vscode.ExtensionContext) {
  console.log("  activating debug sidebar")
  const provider = new WebviewProvider(context, "debug", (_)=>{});
  context.subscriptions.push(
    vscode.window.registerWebviewViewProvider(
      `veg-debug`, // This ID must match package.json
      provider,
      {
        webviewOptions: { retainContextWhenHidden: true }
      }
    )
  );

  const sync = async (s: string) => {
    extensionEmitter.fire({
      type: "requestSync",
      payload: {
        sid: s,
      }
    });
    sendMessage({
      type: "requestSync",
      payload: {
        sid: s,
      }
    })
  }


  // incoming messages
	extensionEmitter.event((e) => {
		// console.log(`debug.panel event:`, e)
		switch (e.type) {
			case "chat.loadSession":
        // console.log("WS SAVE SID:", e.payload?.sid)
        context.workspaceState.update("sid", e.payload?.sid)
				break;
      case "session.delete":
        const curr = context.workspaceState.get("sid")
        if (curr && e.payload?.sid === curr ) {
          context.workspaceState.update("sid", null)
        }

        break;
		}
	});


	vscode.commands.registerCommand('veg.debug.requestSync', () => {
    const sid = context.workspaceState.get("sid") as string
    // console.log("WS LOAD SID:", sid)
    if (sid && sid !== "") {
      sync(sid)
    } else {
      sync("")
    }
	});
}