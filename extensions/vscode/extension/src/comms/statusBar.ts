import * as vscode from 'vscode';

let myStatusBarItem: vscode.StatusBarItem;

export async function activate(context: vscode.ExtensionContext) {
  // Create the status bar item
  myStatusBarItem = vscode.window.createStatusBarItem(
    vscode.StatusBarAlignment.Left,
    100
  );
  myStatusBarItem.command = 'veg.connect';
  context.subscriptions.push(myStatusBarItem);

  // Set initial "connecting" state
  updateStatusBar('Veg Server', 'Connecting...', 'sync~spin');
}

export async function deactivate() {
  if (myStatusBarItem) {
    myStatusBarItem.dispose();
  }
}

/**
 * A shared function to update the status bar text, tooltip, and icon.
 */
export function updateStatusBar(
  text: string,
  tooltip: string,
  iconName?: string // e.g., 'sync~spin', 'check', 'error'
) {
  if (!myStatusBarItem) {
    return;
  }
  myStatusBarItem.text = iconName ? `$(${iconName}) ${text}` : text;
  myStatusBarItem.tooltip = tooltip;
  myStatusBarItem.show();
}