import * as vscode from 'vscode';
import * as path from 'path';
import * as os from 'os';
import { spawn } from 'child_process';
import WebSocket from 'ws';

// --- Import our new shared modules ---
import { updateStatusBar } from './statusBar';
import { extensionEmitter } from './events';

// --- Config & Types (same as before) ---
const SERVER_PORT = 2257;
const SERVER_URL = `ws://localhost:${SERVER_PORT}`;
interface Message<T> { type: string; payload: T; }
interface HelloPayload { version: string; }
interface EchoPayload { text: string; }
// ------------------------------------

let ws: WebSocket | undefined;

export function isConnected(): boolean {
  return ws !== undefined && ws.readyState === WebSocket.OPEN;
}

// Note: We rename `activate` to `activateWebsocket` to avoid
// naming conflicts when we import it.
export async function activate(context: vscode.ExtensionContext) {
  // We no longer manage the status bar item here.
  // We just call the `updateStatusBar` function when our state changes.
  await connectOrSpawnServer(context);

  context.subscriptions.push(
    vscode.commands.registerCommand('veg.connect', async () => {
      deactivate()
      await connectOrSpawnServer(context);
    })
  );

}

export async function deactivate() {
  if (ws) {
    ws.close();
    ws = undefined;
  }
}

export async function connectOrSpawnServer(context: vscode.ExtensionContext) {
  updateStatusBar('Veg Server', 'Connecting...', 'sync~spin'); // <-- We call the imported function
  try {
    ws = await connectToWebSocket();
    console.log('Successfully connected to existing server.');
    updateStatusBar('Veg Server', 'Connected', 'check');
    setupWebSocketHandlers(ws);
  } catch (error) {
    console.log('Server not found. Attempting to start...');
    updateStatusBar('Veg Server', 'Starting...', 'sync~spin');
    if (!startServer(context)) {
      updateStatusBar('Veg Server', 'Failed to start', 'error');
      return;
    }
    try {
      ws = await retryConnection();
      console.log('Successfully started and connected to server.');
      updateStatusBar('Veg Server', 'Connected', 'check');
      setupWebSocketHandlers(ws);
    } catch (retryErr) {
      const msg = retryErr instanceof Error ? retryErr.message : 'Unknown';
      updateStatusBar('Veg Server', `Error: ${msg}`, 'error');
    }
  }
}

function setupWebSocketHandlers(socket: WebSocket) {
  socket.on('message', (data) => {
    const messageStr = data.toString();
    try {
      const msg: Message<unknown> = JSON.parse(messageStr);
      // console.log(`[SERVER]:`, msg);
      extensionEmitter.fire(msg);
    } catch (e) {
      console.error('Error parsing server message', e);
    }
  });

  socket.on('close', () => {
    console.log('WebSocket connection closed.');
    updateStatusBar('Veg Server', 'Connection lost.', 'error');
    ws = undefined;
  });

  socket.on('error', async (err) => {
    console.error(`WebSocket error: ${err.message}`);
    updateStatusBar('Veg Server', `Error: ${err.message}`, 'error');
    ws = undefined;
    try {
      ws = await retryConnection(60);
      console.log('Successfully started and connected to server.');
      updateStatusBar('Veg Server', 'Connected', 'check');
      setupWebSocketHandlers(ws);
    } catch (retryErr) {
      const msg = retryErr instanceof Error ? retryErr.message : 'Unknown';
      updateStatusBar('Veg Server', `Error: ${msg}`, 'error');
    }
  });

  // sendHello(socket);
}

export function sendMessage<T>(msg: Message<T>) {
  if (!ws) {
    vscode.window.showErrorMessage('Server not connected. (sendMessage)');
    return;
  }
  // console.log(`[VSCODE]:`, msg);
  ws.send(JSON.stringify(msg));
}


export function sendEcho(textToSend: string) {
  if (!ws) {
    vscode.window.showErrorMessage('Server not connected. (sendEcho)');
    return;
  }
  const payload: EchoPayload = { text: textToSend };
  const message: Message<EchoPayload> = { type: 'echo', payload };
  ws.send(JSON.stringify(message));
}

// --- Helper Functions (no changes, just not exported) ---

function showQuickPickOptions() {
  // (This logic is unchanged from before)
  vscode.window.showQuickPick([{ label: 'Send Echo Message', command: 'alpha.echo' }])
    .then(selection => {
      if (selection) {
        vscode.commands.executeCommand(selection.command);
      }
    });
}

function sendHello(ws: WebSocket) {
  const payload: HelloPayload = { version: '1.0.0' };
  const message: Message<HelloPayload> = { type: 'hello', payload };
  ws.send(JSON.stringify(message));
}

/**
 * Creates a single connection attempt.
 */
function connectToWebSocket(): Promise<WebSocket> {
  return new Promise((resolve, reject) => {
    const socket = new WebSocket(SERVER_URL, {
      headers: {
        'X-Veg-User': 'tony', // TODO: make this configurable or dynamic
      }
    });
    socket.on('open', () => {
      socket.off('error', reject);
      resolve(socket);
    });
    socket.on('error', (err) => {
      reject(err);
    });
  });
}

/**
 * Tries to connect, retrying every second for 5 seconds.
 */
async function retryConnection(retries: number = 5, timeout: number = 1000): Promise<WebSocket> {
  for (let i = 0; i < retries; i++, timeout *= 1.4) {
    try {
      if (i > 0) {
        await new Promise((resolve) => setTimeout(resolve, timeout));
      }
      return await connectToWebSocket();
    } catch (error) { /* continue */ }
  }
  throw new Error(`Failed to connect to the server after ${retries} attempts.`);
}

/**
 * Tries to spawn the server process in the background.
 */
function startServer(context: vscode.ExtensionContext): boolean {
  // const serverPath = getServerPath(context);
  // if (!serverPath) {
  //   vscode.window.showErrorMessage('Veg server binary not found.');
  //   return false;
  // }
  try {
    const proc = spawn("hof", [`agent`, `--port`, `${SERVER_PORT}`], {
      detached: true,
      stdio: 'ignore',
      windowsHide: true,
    });
    proc.unref();
    return true;
  } catch (err) {
    vscode.window.showErrorMessage(`Error spawning server: ${err}`);
    return false;
  }
}

/**
 * Finds the correct binary path.
 */
function getServerPath(context: vscode.ExtensionContext): string | undefined {
  const platform = os.platform();
  const arch = os.arch();
  let binaryName: string;

  if (platform === 'win32') binaryName = `server-win-${arch}.exe`;
  else if (platform === 'darwin') binaryName = `server-darwin-${arch}`;
  else if (platform === 'linux') binaryName = `server-linux-${arch}`;
  else return undefined;

  return path.join(context.extensionPath, 'server', binaryName);
}