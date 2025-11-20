import * as vscode from 'vscode';

import * as statusBar from './statusBar'
import * as sync from './sync'
import * as treeviews from './treeviews';
import * as websocket from './websocket'
import * as webviews from './webviews';

import { extensionEmitter } from './util/events';

export async function activate(context: vscode.ExtensionContext) {
  console.log('Activating extension "veg-extension"...');

	// important subsystems first
	statusBar.activate(context)
	websocket.activate(context)

	// background monitoring
	sync.activate(context)

	// ui components
	treeviews.activate(context)
	webviews.activate(context)

	// two-way refresh with server
	setTimeout(function() {
		console.log('Sending startup sync');
		extensionEmitter.fire({
			type: "requestSync",
		})
		websocket.sendMessage({
			type: "requestSync",
			payload: {}
		})
	}, 2000);
}

/**
 * Clean up when the window closes.
 */
export function deactivate() {
  console.log('Dectivating extension "veg-extension"...');

	sync.deactivate()
	websocket.deactivate()
	statusBar.deactivate()
}