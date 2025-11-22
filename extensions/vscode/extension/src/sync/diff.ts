import * as vscode from 'vscode';
import { extensionEmitter } from '../util/events';
import { Session } from 'inspector';

// This method is called when your extension is activated
// Your extension is activated the very first time the command is executed
export function activate(context: vscode.ExtensionContext) {

	const vcp = new VegContentProvider()

	vscode.workspace.registerTextDocumentContentProvider("veg://", vcp)

	context.subscriptions.push(
		vscode.workspace.onDidOpenTextDocument(e => {
			if (e.uri.scheme !== 'file') {
				return;
			}
			// console.log("openedDocument", e.fileName)
		})
	);
	vscode.commands.registerCommand('veg.debug.diff', () => {
		const dir = "/Users/tony/hof/hof/extensions/vscode/extension"
		const lhs = vscode.Uri.file(dir + "/src/sync/diff.ts")
		const rhs = vscode.Uri.file(dir + "/src/sync/env.ts")
		vscode.commands.executeCommand('vscode.diff', lhs, rhs, "test diff", {
			viewColumn: 1
		})
	})

	extensionEmitter.event((e) => {
		switch (e.type) {
			case "diff.write_file":
				console.log("diff.write_file", e)
				var {
					sid,
					path,
					content,
				}: {
					sid: string,
					path: string,
					content: string,
				} = e.payload;

				if (!path.startsWith("/")) {
					path = "/" + path
				}

				const uri = vscode.Uri.parse(`file://${path}&sid=${sid}`)
				vcp.addFile(uri, content)

      case "diff.show":
				console.log("diff.show", e)
				var {
					lhs,
					rhs,
				}: {
					lhs: string,
					rhs: string,
				} = e.payload;

				vscode.commands.executeCommand('vscode.diff', lhs, rhs, "test diff", {
					viewColumn: 1
				})
				break;

      case "diff.nextStep":
				console.log("diff.nextStep", e)
				const { sid: s } = e.payload
				vcp.nextStep(s)

				break;
		}
	});
}

// This method is called when your extension is deactivated
export function deactivate() {}


// doing this so we can view diffs across agentic writes
// and accept or fork at various points
// todo, we need to persist this somewhere, the database or on the filesystem?
interface SessionFS {
	sid: string
	steps: SessionStep[]
}

interface SessionStep {
	files: Record<string,string>
}

class VegContentProvider implements vscode.TextDocumentContentProvider {
	private fs: Record<string, SessionFS>

	constructor(){
		this.fs = {}
	}

	provideTextDocumentContent(uri: vscode.Uri, token: vscode.CancellationToken): vscode.ProviderResult<string> {
		const qp = new URLSearchParams(uri.query)
		const s = qp.get("session") as string
		const t = Number(qp.get("step"))
		return this.fs[s].steps[t].files[uri.path]
	}

	addFile(uri: vscode.Uri, content: string) {
		const qp = new URLSearchParams(uri.query)
		const s = qp.get("session") as string
		this.ensureExists(s)
		const S = this.fs[s].steps

		S[S.length-1].files[uri.path] = content
	}

	nextStep(session: string) {
		if (this.ensureExists(session)) {
			this.fs[session].steps.push({ files: {} })
		}
	}

	// returns true if exists, creates and returns false if not
	private ensureExists(s: string): boolean {
		if (!(s in this.fs)) {
			this.fs[s] = {
				sid: s,
				steps: [{ files: {} }],
			}
			return false
		}
		return true
	}
}