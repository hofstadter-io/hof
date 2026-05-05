import * as path from 'path';
import * as vscode from 'vscode';

import { extensionEmitter } from '../comms';
import { getHtmlForWebview } from './loader';

export class WebviewProvider implements vscode.WebviewViewProvider {
  private _view?: vscode.WebviewView;

  constructor(
    private readonly _context: vscode.ExtensionContext,
    public readonly _name: string,
    private _onMessage: (data: any) => void
  ) {}

  public async resolveWebviewView(
    webviewView: vscode.WebviewView,
    context: vscode.WebviewViewResolveContext,
    _token: vscode.CancellationToken
  ) {
    this._view = webviewView;

    console.log("resolveWebviewView", this._context)

    webviewView.webview.options = {
      // Allow scripts in the webview
      enableScripts: true,
      // Restrict the webview to only loading content from our extension's directories
      localResourceRoots: [
        vscode.Uri.file(path.join(this._context.extensionPath, 'webviews', this._name)),
      ],
    };

    // Set the HTML content
    webviewView.webview.html = await getHtmlForWebview(webviewView.webview, this._context.extensionPath, this._name);

    // --- Communication ---

    // 1. Listen for messages from the webview (Chat UI -> Extension)
    webviewView.webview.onDidReceiveMessage(this._onMessage);

    // 2. Listen for messages from the websocket (Extension -> Chat UI)
    extensionEmitter.event((e) => {
      // console.log(`${this._name} event:`, e)
      // if (e.type === 'WEBSOCKET_MESSAGE') {
        // We got a message. Pass it to the webview.
        this._view?.webview.postMessage(e);
      // }
    });
  }
}