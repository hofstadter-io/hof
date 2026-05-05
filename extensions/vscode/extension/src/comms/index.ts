import * as vscode from 'vscode';

export { extensionEmitter } from './events';
export { sendMessage, isConnected, connectOrSpawnServer } from './websocket';
export { updateStatusBar } from './statusBar';

import { activate as wsActivate, deactivate as wsDeactivate } from './websocket';
import { activate as sbActivate, deactivate as sbDeactivate } from './statusBar';


export async function activate(context: vscode.ExtensionContext) {
  console.log("  activating comms system")

  await sbActivate(context)
  await wsActivate(context)
}

export async function deactivate() {
  await wsDeactivate()
  await sbDeactivate()
}