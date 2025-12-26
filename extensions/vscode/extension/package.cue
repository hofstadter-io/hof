package extension

name:        "veg"
displayName: "Veg"
description: "Veg VS Code Extension"
publisher:   "verdverm"
version:     "0.0.1"
engines: {
	vscode: "^1.106.0"
}
categories: [
	"Other",
]
enabledApiProposals: [
	"scmMultiDiffEditor",
]
activationEvents: [
	"onStartupFinished",
]
main: "./out/extension.js"
capabilities: {
	virtualWorkspaces: true
}

contributes: {

	viewsContainers: {
		activitybar: [
			{id: "veg-chat-sidebar", title: "Veggie", icon: "media/leafy-green.svg"},
			{id: "veg-manage-sidebar", title: "Veggie", icon: "media/carrot.svg"},
		]
		panel: [
			{id: "veg-debug-panel", title: "Veggie", icon: "media/sprout.svg"},
		]
	}

	views: {
		"veg-chat-sidebar": [
			{id: "veg-chat", name: "Chat", type: "webview", contextualTitle: "Chat"},
		]
		"veg-manage-sidebar": [
			{id: "veg-manage", name: "Manage", type: "webview", contextualTitle: "Manage"},
		]
		"veg-debug-panel": [
			{id: "veg-debug", name: "Debug", type: "webview", contextualTitle: "Debug"},
			{id: "veg-sessions", name: "Sessions", contextualTitle: "Sessions"},
			// { id: "veg-planning", name: "Planning", contextualTitle: "Planning" },
			// { id: "veg-agents", name: "Agents", contextualTitle: "Agents" },
		]
	}

	keybindings: [
		{command: "veg-chat.focus", key: "ctrl+g", mac: "cmd+g"},
		{command: "veg-chat.focus", key: "ctrl+shift+g", mac: "cmd+shift+g"},
		{command: "veg-chat.focus", key: "alt+g g", mac: "alt+g g"},
		{command: "veg-sessions.focus", key: "alt+g s", mac: "alt+g s"},
		// { command: "veg-agents.focus", key: "alt+g a", mac: "alt+g a" },
		// { command: "veg-planning.focus", key: "alt+g p", mac: "alt+g p" },
		{command: "veg-debug.focus", key: "alt+g d", mac: "alt+g d"},
	]

	commands: [
		{command: "veg.connect", category: "Veg Connect", title: "Veg Connect"},

		{command: "veg.filesys.hack", category: "Veg", title: "Hack Filesys (Veg)", icon: "$(symbol-class)"},

		{command: "veg.explorer.terminal", category: "Veg", title: "Terminal (Veg)", icon: "$(terminal)"},
		{command: "veg.explorer.copyPath", category: "Veg", title: "Copy Path (Veg)", icon: "$(copy)"},
		{command: "veg.explorer.chat", category: "Veg", title: "Chat (Veg)", icon: "$(comment-discussion)"},
		{command: "veg.explorer.openEnviron", category: "Veg", title: "Open (Veg)", icon: "$(folder-opened)"},
		{command: "veg.explorer.forkEnviron", category: "Veg", title: "Fork (Veg)", icon: "$(gist-fork)"},
		{command: "veg.explorer.toggleShown", category: "Veg", title: "Toggle Diff View (Veg)", icon: "$(filter)"},
		{command: "veg.explorer.showDiff", category: "Veg", title: "Show Diff (Veg)", icon: "$(diff-multiple)"},
		{command: "veg.explorer.showFileDiff", category: "Veg", title: "Show File Diff (Veg)", icon: "$(diff)"},
		{command: "veg.explorer.mergeDiff", category: "Veg", title: "Merge Diff (Veg)", icon: "$(git-pull-request-create)"},
		{command: "veg.explorer.hideDiff", category: "Veg", title: "HideDiff Diff (Veg)", icon: "$(eye-closed)"},
		{command: "veg.explorer.diffAll", category: "Veg", title: "Diff All (Veg)", icon: "$(group-by-ref-type)"},
		{command: "veg.explorer.refreshAll", category: "Veg", title: "Refresh (Veg)", icon: "$(clear-all)"},

		{command: "veg.session.openEnviron", category: "Veg", title: "Show Session Files (Veg)", icon: "$(list-tree)"},
		{command: "veg.session.showSessionDiff", category: "Veg", title: "Show Session Diff (Veg)", icon: "$(diff-multiple)"},
		{command: "veg.session.showFileDiff", category: "Veg", title: "Show File Diff (Veg)", icon: "$(diff-multiple)"},
		{command: "veg.session.chat", category: "Veg", title: "veg.sessions.chat", icon: "$(comment-discussion)"},
		{command: "veg.session.fork", category: "Veg", title: "veg.sessions.fork", icon: "$(git-branch)"},
		{command: "veg.session.edit", category: "Veg", title: "veg.sessions.edit", icon: "$(pencil)"},
		{command: "veg.session.delete", category: "Veg", title: "veg.sessions.delete", icon: "$(trash)"},

		{command: "veg.sessions.refresh", category: "Veg", title: "Refresh Session List (Veg)", icon: "$(refresh)"},
		{command: "veg.sessions.create", category: "Veggie", title: "Fresh Veggie", icon: "$(add)"},

		{command: "veg.debug.requestSync", category: "Veg Debug", title: "veg.debug.requestSync", icon: "$(refresh)"},
		// { command: "veg.debug.terminal", category: "Veg Debug", title: "veg.debug.requestSync", icon: "$(terminal)" },
	]

	submenus: [{
		id:    "veg.explorerSubmenu"
		label: "More (Veg)"
	}]

	menus: {
		"explorer/context": [
			{command: "veg.explorer.chat", group: "_veg@0"},
			{submenu: "veg.explorerSubmenu", group: "_veg@1"},
		]
		"veg.explorerSubmenu": [
			{command: "veg.explorer.openEnviron", group: "_veg@2"},
			{command: "veg.explorer.forkEnviron", group: "_veg@3"},
			{command: "veg.explorer.toggleShown", group: "_veg@4"},
			{command: "veg.explorer.showDiff", group: "_veg@5"},
			{command: "veg.explorer.mergeDiff", group: "_veg@6"},
			{command: "veg.explorer.copyPath", group: "_veg@7"},
			{command: "veg.explorer.refreshAll", group: "_veg@8"},
		]
		"editor/title": [
			{command: "veg.explorer.refreshAll", group: "navigation"},
			{command: "veg.explorer.toggleShown", group: "navigation"},
			{command: "veg.session.showSessionDiff", group: "navigation"},
		]
		"view/title": [
			{command: "veg.sessions.create", group: "navigation", when: "view == veg-chat"},
			{command: "veg.debug.requestSync", group: "navigation", when: "view == veg-manage"},
			{command: "veg.debug.requestSync", group: "navigation", when: "view == veg-debug"},
			{command: "veg.debug.requestSync", group: "navigation", when: "view == veg-sessions"},
			{command: "veg.sessions.create", group: "navigation", when: "view == veg-sessions"},
		]
		"view/item/context": [
      for i, cmd in [
        "chat",
        "openEnviron",
        "terminal",
        "showDiff",
        "mergeDiff",
        "clone",
        "delete"
      ] {
        when: "view == veg-sessions && viewItem == session"
        group: "inline@\(i)"
        command: "veg.session.\(cmd)"
      }
		]
		"scm/title": [
			{command: "veg.explorer.chat", group: "navigation@0", when: "scmProvider == git || scmProvider == veg"},

			{command: "veg.explorer.showDiff", group: "navigation@11", when: "scmProvider == veg"},
			{command: "veg.explorer.mergeDiff", group: "navigation@12", when: "scmProvider == veg"},
			{command: "veg.explorer.hideDiff", group: "navigation@13", when: "scmProvider == veg"},
			{command: "veg.explorer.diffAll", group: "navigation@14", when: "scmProvider == veg"},
		]

		"scm/resourceGroup/context": [
			{command: "veg.explorer.showDiff", group: "inline@1", when: "scmProvider == veg"},
			{command: "veg.explorer.mergeDiff", group: "inline@2", when: "scmProvider == veg"},
			{command: "veg.explorer.hideDiff", group: "inline@3", when: "scmProvider == veg"},
			{command: "veg.explorer.terminal", group: "inline@4", when: "scmProvider == veg"},
		]
	}

}
"scripts": {
	"vscode:prepublish": "pnpm run compile"
	"compile":           "tsc -p ./"
	"watch":             "tsc -watch -p ./"
	"pretest":           "pnpm run compile && pnpm run lint"
	"lint":              "eslint src"
	"test":              "vscode-test"
}
"devDependencies": {
	"@types/mocha":          "^10.0.10"
	"@types/node":           "22.x"
	"@types/vscode":         "^1.106.0"
	"@types/ws":             "^8.18.1"
	"@vscode/test-cli":      "^0.0.12"
	"@vscode/test-electron": "^2.5.2"
	"eslint":                "^9.39.1"
	"typescript":            "^5.9.3"
	"typescript-eslint":     "^8.46.3"
}
"dependencies": {
	"ansi-colors":  "^4.1.3"
	"jsonc-parser": "^3.3.1"
	"ws":           "^8.18.3"
}
