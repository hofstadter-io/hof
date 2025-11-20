import * as vscode from 'vscode';

import { activate as planning } from './planning'
import { activate as sessions } from './sessions'
import { activate as agents } from './agents'

export function activate(context: vscode.ExtensionContext) {
  console.log(`activating sidebars`)

  planning(context)
	sessions(context)
	agents(context)
}


