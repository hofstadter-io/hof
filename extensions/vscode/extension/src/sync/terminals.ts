import * as vscode from 'vscode';
import c from 'ansi-colors';
import { extensionEmitter } from '../util/events';
import { sendMessage } from '../websocket';

const trackedTerminals = new Set<Terminal>();
var termIndex = 0

class Exec {
	start: vscode.TerminalShellExecutionStartEvent | null = null
	output: string = ""
	end: vscode.TerminalShellExecutionEndEvent | null = null
}

class Terminal {
	termIndex: number = 0;
	terminal: vscode.Terminal | null = null;
	history: Exec[] = [];	
}

class TerminalPayload {
	id: number = -1;
	name?: string;
	history?: HistoryPayload[];
}

class HistoryPayload {
	cmd?: any;
	cwd?: string;
	out?: string;
	exit?: number;
}

export function getTerminals(): TerminalPayload[] {
	const terms: TerminalPayload[] = []
	for (const term of trackedTerminals.values()) {
		terms.push({
			id: term.termIndex,
			name: term.terminal?.name,
			history: term.history?.map((h => {
				const H: HistoryPayload = {
					cmd: h.end?.execution.commandLine,
					cwd: h.start?.terminal.shellIntegration?.cwd?.path,
					out: h.output,
					exit: h.end?.exitCode,
				}

				return H
			}))
		})
	}
	return terms
}

function broadcastTerminals() {
	const msg = {
		type: "terminalInfo",
		payload: {
			terminals: getTerminals(),
		}
	}
	extensionEmitter.fire(msg);
	sendMessage(msg)
}

function findTerm(vsterm: vscode.Terminal): Terminal | null { 
	console.log("find:", vsterm, trackedTerminals)
	for (const term of trackedTerminals.values()) {
		if (term.terminal == vsterm) {
			return term
		}
	}
	return null
}

function finalizeExec(term: Terminal, end: vscode.TerminalShellExecutionEndEvent) {
	if (end.execution.commandLine.value === "") {
		console.warn("ignoring empty command")
		return
	}
	// search backwards because we push to history, most probable at the end
	for (var t = term.history.length - 1; t >= 0; t--) {
		const h = term.history[t]
		if (h.start?.execution === end.execution) {
			h.end = end
			break
		}
	}
}

// This method is called when your extension is activated
// Your extension is activated the very first time the command is executed
export function activate(context: vscode.ExtensionContext) {

	context.subscriptions.push(

		vscode.window.onDidChangeTerminalShellIntegration(e => {
			console.log("changeIntegration", e, trackedTerminals)
			if (!findTerm(e.terminal)){
				console.log("changeIntegration.newTerminal", e)
				const t = new Terminal()
				t.terminal = e.terminal
				t.termIndex = termIndex
				termIndex++
				trackedTerminals.add(t);
				broadcastTerminals()
			}
		}),

		vscode.window.onDidStartTerminalShellExecution(async e => {
			console.log("execStart", e, trackedTerminals)
			var t = findTerm(e.terminal)
			if (!t) {
				console.log("execStart.newTerminal", e)
				t = new Terminal()
				t.terminal = e.terminal
				t.termIndex = termIndex
				termIndex++
				trackedTerminals.add(t);
				broadcastTerminals()
			}

			if (e.execution.commandLine.value === "") {
				console.warn("ignoring empty command")
				return
			}

			// new history entry
			const h = new Exec();
			t.history.push(h)
			h.start = e
			broadcastTerminals()

			console.log("hist", h)
			// collect output stream
			const stream = e.execution.read();
			for await (const data of stream) {
				h.output += c.unstyle(data)
			}	
			broadcastTerminals()
		}),

		vscode.window.onDidEndTerminalShellExecution(e => {
			if (e.execution.commandLine.value === "") {
				console.warn("ignoring empty command")
				return
			}
			console.log("execEnd", e)
			const t = findTerm(e.terminal)
			if (!t) {
				console.error("failed to find terminal for:", e)
				return
			}
			finalizeExec(t, e)
			console.log("done:", t)
			broadcastTerminals()
		}),

		vscode.workspace.onDidOpenTextDocument(e => {
			if (e.uri.scheme !== 'file') {
				return;
			}
			// console.log("openedDocument", e.fileName)
		})
	);
	// Use the console to output diagnostic information (console.log) and errors (console.error)
	// This line of code will only be executed once when your extension is activated
	// console.log('Congratulations, your extension "alpha" is now active!');

	// The command has been defined in the package.json file
	// Now provide the implementation of the command with registerCommand
	// The commandId parameter must match the command field in package.json
	const disposable = vscode.commands.registerCommand('alpha.helloWorld', () => {
		// The code you place here will be executed every time your command is executed
		// Display a message box to the user
		console.log(JSON.stringify(vscode.window.terminals, null, "  "));
		vscode.window.showInformationMessage('Hallo World from alpha!');
	});

	vscode.commands.registerCommand('alpha.infoTerminals', () => {
		// The code you place here will be executed every time your command is executed
		// Display a message box to the user
		const terms = getTerminals();
		console.log(JSON.stringify(terms, null, "  "));
		broadcastTerminals()
	});

	extensionEmitter.event((e) => {
		console.log(`sync.terminals event:`, e)
		switch (e.type) {
			case "requestSync":
				broadcastTerminals()
				break;
		}
	});

	context.subscriptions.push(disposable);
}

// This method is called when your extension is deactivated
export function deactivate() {}
