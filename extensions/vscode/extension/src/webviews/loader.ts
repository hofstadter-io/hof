import * as path from 'path';
import * as vscode from 'vscode';

/**
 * Gets the HTML content for the webview.
 * This is where we load our built React app.
 */
export async function getHtmlForWebview(webview: vscode.Webview, extensionPath: string, name: string): Promise<string> {
  // 1. Get the path
  const buildDir = path.join(extensionPath, 'webviews', name);
  const indexPath = path.join(buildDir, 'index.html');
  const indexUri = vscode.Uri.file(indexPath);

  // 2. Read the file's content
  let html: string;
  try {
    const fileContent = await vscode.workspace.fs.readFile(indexUri);
    html = Buffer.from(fileContent).toString('utf-8');
    // --- NEW: Log the raw HTML ---
    // console.log('Raw index.html content:\n', html.substring(0, 300) + '...');
  } catch (e) {
    console.error('Error reading index.html:', e);
    return `<html><body>Failed to load webview. ${e}</body></html>`;
  }

  // 3. Create a function to replace asset paths
  const getVscodeResourceUri = (filePath: string) => {
    return webview.asWebviewUri(
      vscode.Uri.file(path.join(buildDir, filePath))
    );
  };

  // 4. Use regex to replace all 'src' and 'href' paths
  let replacementCount = 0;
  html = html.replace(
    /(href|src)="(.*?)"/g,
    (match, p1: 'href' | 'src', p2: string) => {
      replacementCount++;
      const assetPath = p2.startsWith('/') ? p2.substring(1) : p2;
      const newSrc = getVscodeResourceUri(assetPath);
      // --- NEW: Log replacements ---
      // console.log(`Replacing: ${p2} \n      WITH: ${newSrc}`);
      return `${p1}="${newSrc}"`;
    }
  );

  // --- NEW: Log final state ---
  // console.log(`Total replacements: ${replacementCount}`);
  if (replacementCount === 0) {
    console.warn('WARNING: No src/href attributes were replaced. Check regex.');
  }

  // 5. Add the nonce (unchanged)
  const nonce = getNonce();
  html = html.replace(
    /<(script|link|style)/g,
    (match, p1) => `<${p1} nonce="${nonce}"`
  );

  // 6. Set the Content Security Policy (unchanged)
  html = html.replace(
    '<head>',
    `<head>
    <meta http-equiv="Content-Security-Policy" content="
      default-src 'none';
      style-src ${webview.cspSource} 'unsafe-inline';
      font-src ${webview.cspSource};
      script-src 'nonce-${nonce}';
      img-src ${webview.cspSource} https://cdn.jsdelivr.net https: data:;
    ">`
  );

  // --- NEW: Log final HTML ---
  // console.log('Final HTML sent to webview:\n', html.substring(0, 500) + '...');
  // console.log('--- End Webview Debug ---');

  return html;
}

// A helper function to generate a random nonce
function getNonce() {
let text = '';
const possible = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz01234s56789';
for (let i = 0; i < 32; i++) {
  text += possible.charAt(Math.floor(Math.random() * possible.length));
}
return text;
}