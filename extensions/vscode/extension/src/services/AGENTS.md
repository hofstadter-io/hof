# Services Subsystem (`src/services`)

This directory is currently dedicated to the **Veg Virtual Filesystem Provider**, which enables interacting with remote environments (like Dagger containers or session states) as if they were local workspace folders in VS Code.

## Filesystem Provider (`filesystemProvider.ts`)

The `VegContentProvider` is the sole and comprehensive implementation of `vscode.FileSystemProvider` for the `veg://` scheme. It manages command registration, URI translation, and communication with the extension backend API (`http://localhost:2257`).

### Core Data Structures and Types

```typescript
const SERVER_PORT = 2257;
const SERVER_URL = `http://localhost:${SERVER_PORT}`;

type Environ = {
	name?: string
	srcUri?: string
	srcPath?: string
	fromUri?: string
	dstPath?: string
	workdir?: string
}

type Folder = {
	uri: vscode.Uri
	sid: string
	name?: string | undefined
	base?: string | undefined
	session?: any
	environ?: Environ
}
type FolderListing = Array<[string, vscode.FileType]>
```

### URI Translation

The `vsUriToVeg` method is crucial for converting the VS Code internal URI representation (which uses path segments after the authority) into the backend's expected canonical URI format (which embeds the path in a query parameter).

```typescript
private vsUriToVeg(uri: vscode.Uri): vscode.Uri {
	// ... logic to parse 'veg://<auth>/<env>/<path>' into a base OCI-style URI
	// with a 'path' query parameter for the file path
}
```

### Command Registration

The file registers a large set of `veg.explorer.*` commands, primarily for file and environment management in the file explorer context menu:
- `veg.explorer.chat`: Initiates a new session focused on the selected file/directory.
- `veg.explorer.openEnviron`: Prompts the user to open a new directory/repo/image.
- `veg.explorer.showDiff`: Opens a side-by-side diff for all modified paths in the environment.
- `veg.explorer.mergeDiff`: Writes the modified, added, and deleted files from the virtual environment to a specified local `file://` target, utilizing `node:fs/promises`.

### Diff and Merge Implementation

The `showDiff` and `mergeDiff` commands fetch a unified diff payload from the backend (`/fs/diff`) for a given URI.

- `showDiff`: Iterates over `diff.modPaths` and opens multiple `vscode.diff` editors.
- `mergeDiff`: Uses the `diff` payload to perform native filesystem operations (`fs.writeFile`, `fs.rm`) on the destination URI (which must be `file://`).

```typescript
// mergeDiff logic snippet (uses node:fs/promises)
for (var path of diff.addPaths) {
  // ...
  const val = diff.files[path]
  const key = destination.path + path
  await fs.writeFile(key, val)
}
// ... similar for modPaths and delPaths
```
