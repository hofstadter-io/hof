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
}

/**
 * Clean up when the window closes.
 */
export async function deactivate() {
  console.log('Dectivating extension "veg-extension"...');

	sync.deactivate()
	await comms.deactivate()
}