import * as vscode from 'vscode';

// This emitter will be used to pass messages between
// the websocket (which receives) and the sidebar (which displays).
export const extensionEmitter = new vscode.EventEmitter<any>();