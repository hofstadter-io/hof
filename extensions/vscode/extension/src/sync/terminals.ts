import * as vscode from 'vscode';
import c from 'ansi-colors';
import { extensionEmitter, sendMessage } from '../comms';

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
		type: "terminal.info",
		payload: {
			terminals: getTerminals(),
		}
	}
	extensionEmitter.fire(msg);
	sendMessage(msg)
}

function findTerm(vsterm: vscode.Terminal): Terminal | null { 
	// console.log("find:", vsterm, trackedTerminals)
	for (const term of trackedTerminals.values()) {
		if (term.terminal == vsterm) {
			return term
		}
	}
	return null
}

function deleteTerm(vsterm: vscode.Terminal): Terminal | null { 
	// console.log("find:", vsterm, trackedTerminals)
	for (const term of trackedTerminals.values()) {
		if (term.terminal == vsterm) {
			trackedTerminals.delete(term)
			return term
		}
	}
	return null
}

function finalizeExec(term: Terminal, end: vscode.TerminalShellExecutionEndEvent) {
	if (end.execution.commandLine.value === "") {
		console.warn("terminal.finalizeExec: ignoring empty command")
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
	var sessions: any[] = []
	let disposable = vscode.commands.registerCommand('veg.debug.terminal', () => {
			// 1. Create the terminal
			const terminal = vscode.window.createTerminal({
					name: "Veg Terminal",
					// shellPath: "..." // Optional: specify a shell
			});

			// 2. Show the terminal to the user (focus it)
			terminal.show();

			// 3. Send the command text
			// The second argument 'true' implies adding a newline to execute immediately
			terminal.sendText("echo 'Hello from VS Code Extension!'", true);
	});

	context.subscriptions.push(disposable);

	context.subscriptions.push(

		vscode.window.onDidChangeTerminalShellIntegration(e => {
			console.log("changeIntegration", e, trackedTerminals)
			if (!findTerm(e.terminal)){
				// console.log("changeIntegration.newTerminal", e)
				const t = new Terminal()
				t.terminal = e.terminal
				t.termIndex = termIndex
				termIndex++
				trackedTerminals.add(t);
				broadcastTerminals()
			}
		}),

		vscode.window.onDidCloseTerminal(async e => {
			console.log("close Terminal", e)
			var t = findTerm(e)
			if (t) {
				deleteTerm(e)
				console.log("DELETED:", t)
			} else {
				console.log("did not find!")
			}

		}),

		vscode.window.onDidStartTerminalShellExecution(async e => {
			// console.log("execStart", e, trackedTerminals)
			var t = findTerm(e.terminal)
			if (!t) {
				// console.log("execStart.newTerminal", e)
				t = new Terminal()
				t.terminal = e.terminal
				t.termIndex = termIndex
				termIndex++
				trackedTerminals.add(t);
				broadcastTerminals()
			}

			if (e.execution.commandLine.value === "") {
				// console.warn("ignoring empty command")
				return
			}

			// new history entry
			const h = new Exec();
			t.history.push(h)
			h.start = e
			broadcastTerminals()

			// console.log("hist", h)
			// collect output stream
			const stream = e.execution.read();
			for await (const data of stream) {
				h.output += c.unstyle(data)
			}	
			broadcastTerminals()
		}),

		vscode.window.onDidEndTerminalShellExecution(e => {
			if (e.execution.commandLine.value === "") {
				// console.warn("ignoring empty command")
				return
			}
			// console.log("execEnd", e)
			const t = findTerm(e.terminal)
			if (!t) {
				console.error("failed to find terminal for:", e)
				return
			}
			finalizeExec(t, e)
			// console.log("done:", t)
			broadcastTerminals()
		}),

		vscode.workspace.onDidOpenTextDocument(e => {
			if (e.uri.scheme !== 'file') {
				return;
			}
			// console.log("openedDocument", e.fileName)
		})
	);

	extensionEmitter.event((e) => {
		// console.log(`sync.terminals event:`, e)
		switch (e.type) {
			case "requestSync":
				broadcastTerminals()
				break;
			case "session.list.resp":
				sessions = e.payload
				break;

			case "session.term.open":
				// process inputs
				const { sid, pos, image, termId } = e.payload


				const S = sessions.filter((s) => s.sid === sid)[0]
				if (!S) {
					console.error("unknown session", sid)
				}

				let img = image || S.state?.currEnv || "debian:13-slim"
				if (!image && pos !== undefined && S.events) {
					for (let i = pos; i >= 0; i--) {
						const event = S.events[i];
						if (event?.Actions?.StateDelta?.currEnv) {
							img = event.Actions.StateDelta.currEnv;
							break;
						}
					}
				}

				const lastColon = img.lastIndexOf(":")
				const tag = lastColon !== -1 ? img.substring(lastColon + 1) : ""
				let name = S.state?.title || sid
				if (tag !== "" && !isNaN(parseInt(tag))) {
					name = `(${tag}) ${name}`
				}
				if (pos !== undefined && pos >= 0) {
					name = `[${pos}] ${name}`
				}

				// todo, go to position in events to get dagger ref
				const workdir = S.state?.basedir || S.state?.env?.workdir

				// run dagger via hof for arg handling
				let env = `_EXPERIMENTAL_DAGGER_RUNNER_HOST=container://veg-dagger-engine`
        const busted = new Date().toLocaleString()

				let cmd = `${env} dagger -i core container with-env-variable --name BUSTED_CACHE --value '${busted}' from --address "${img}"`
				if (workdir) {
					cmd += ` with-workdir --path "${workdir}"`
				}
				cmd += ` terminal --cmd zsh`

				// Create and show the terminal
				let terminal: vscode.Terminal;
				const T = termId !== undefined ? findTermById(termId) : null;
				if (T && T.terminal) {
					terminal = T.terminal;
				} else {
					terminal = vscode.window.createTerminal({ 
						name,
						iconPath: new vscode.ThemeIcon("hubot")
					});
				}
				terminal.show();

				terminal.sendText(cmd, true);

				break;
		}
	});

	// context.subscriptions.push(disposable);
}

function findTermById(id: number): Terminal | null { 
	for (const term of trackedTerminals.values()) {
		if (term.termIndex == id) {
			return term
		}
	}
	return null
}

// This method is called when your extension is deactivated
export function deactivate() {}

