import * as vscode from 'vscode';

import * as terminals from './terminals'

export function activate(context: vscode.ExtensionContext) {
  console.log(`activating sync features`)

  terminals.activate(context)
}

export function deactivate() {
  console.log(`deactivating sync features`)

  terminals.deactivate()
}