import * as vscode from 'vscode';

import { activate as config } from './config'
import { activate as debug } from './debug'
import { activate as chat } from './chat'

export function activate(context: vscode.ExtensionContext) {
  console.log(`activating webviews`)

	// chat panel
  chat(context)

  // manage panel
  config(context)
  debug(context)

}

