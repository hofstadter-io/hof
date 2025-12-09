import * as vscode from 'vscode';

import { activate as sessions } from './sessions'

export async function activate(context: vscode.ExtensionContext) {
  console.log(`activating treeviews`)

	await sessions(context)
}


