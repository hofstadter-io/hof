import * as vscode from 'vscode';

import * as comms from './comms'
import * as sync from './sync'
import * as treeviews from './treeviews';
import * as webviews from './webviews';
import * as filesys from './services/filesystemProvider'

export async function activate(context: vscode.ExtensionContext) {
  console.log('Activating extension "veg-extension"...');

	const state = vscode.window.state
  console.log("  state", context.globalState.keys())

	// important subsystems first
	await comms.activate(context)

	await filesys.activate(context)

	// background monitoring
	sync.activate(context)

	// ui components
	await webviews.activate(context)
	await treeviews.activate(context)

	// two-way refresh with server
	console.log('Sending startup sync');
	comms.extensionEmitter.fire({
		type: "requestSync",
	})
	comms.sendMessage({
		type: "requestSync",
		payload: {}
	})

	// Setup a background timer for periodic synchronization
	const syncTimer = setInterval(async () => {
		// Try to reconnect if not connected
		if (!comms.isConnected()) {
			await comms.connectOrSpawnServer(context);
		}

		const sid = context.workspaceState.get("sid") as string;
		
		// Always request general sync
		const msg = { type: "requestSync", payload: { sid } };
		comms.extensionEmitter.fire(msg);
		comms.sendMessage(msg);

		// If we have an active session, refresh its state
		if (sid) {
			const sessionMsg = { type: "session.get", payload: { sid } };
			comms.extensionEmitter.fire(sessionMsg);
			comms.sendMessage(sessionMsg);
		}
	}, 5 * 1000);

	context.subscriptions.push({ dispose: () => clearInterval(syncTimer) });
}

/**
 * Clean up when the window closes.
 */
export async function deactivate() {
  console.log('Dectivating extension "veg-extension"...');

	sync.deactivate()
	await comms.deactivate()
}