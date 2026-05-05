import * as vscode from 'vscode';

import { activate as manage } from './manage'
import { activate as debug } from './debug'
import { activate as chat } from './chat'

export function activate(context: vscode.ExtensionContext) {
  console.log(`activating webviews`)

  debug(context)
  manage(context)
  chat(context)

}

