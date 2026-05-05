import * as vscode from 'vscode';

import * as env from './env'
import * as terminals from './terminals'
import * as window from './window'
import * as workspace from './workspace'

const cs = [
  env,
  window,
  terminals,
  workspace,
]

export function activate(context: vscode.ExtensionContext) {
  console.log(`activating sync features`)
  cs.forEach( c => c.activate(context) )
}

export function deactivate() {
  console.log(`deactivating sync features`)
  cs.reverse().forEach( c => c.deactivate() )
}