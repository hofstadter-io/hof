// @ts-nocheck
const vscode = require('vscode');
const { encoding_for_model, get_encoding } = require('tiktoken');
const { countTokens, getTokenizer } = require('@anthropic-ai/tokenizer');
const { TextDecoder } = require('util');

const utf8Decoder = new TextDecoder('utf-8');

const CONFIG_SECTION = 'gpt-token-counter-live';
const DEFAULT_EVEN_COLOR = '#B8D4FF';
const DEFAULT_ODD_COLOR = '#FFE0A6';
const DEFAULT_STATUS_TEMPLATE = 'Token Count: {count} ({family})';

const MODEL_FAMILIES = {
    'openai': 'GPT',
    'anthropic': 'Claude',
    'gemini': 'Gemini'
};

const HIGHLIGHT_EVEN_KEY = 'highlightEvenColor';
const HIGHLIGHT_ODD_KEY = 'highlightOddColor';
const NON_HIGHLIGHT_SCHEMES = new Set([
    'output',
    'vscode-output',
    'vscode',
    'walkThrough',
    'walkThroughSnippet',
    'debug',
    'git',
    'gitlens-git'
]);

function sanitizeColorSetting(value, fallback) {
    if (typeof value !== 'string') {
        return fallback;
    }
    const trimmed = value.trim().toUpperCase();
    if (/^#[0-9A-F]{6}$/.test(trimmed)) {
        return trimmed;
    }
    if (/^#[0-9A-F]{8}$/.test(trimmed)) {
        return trimmed;
    }
    return fallback;
}

function stripAlpha(hex) {
    if (typeof hex !== 'string') {
        return '#000000';
    }
    const clean = hex.replace('#', '').toUpperCase();
    if (clean.length === 6) {
        return `#${clean}`;
    }
    if (clean.length === 8) {
        return `#${clean.slice(0, 6)}`;
    }
    return '#000000';
}

function hexToLinearRgb(hex) {
    const noAlpha = stripAlpha(hex).slice(1);
    const r = parseInt(noAlpha.slice(0, 2), 16) / 255;
    const g = parseInt(noAlpha.slice(2, 4), 16) / 255;
    const b = parseInt(noAlpha.slice(4, 6), 16) / 255;

    const toLinear = (channel) => {
        return channel <= 0.03928 ? channel / 12.92 : Math.pow((channel + 0.055) / 1.055, 2.4);
    };

    return {
        r: toLinear(r),
        g: toLinear(g),
        b: toLinear(b)
    };
}

function textColorForBackground(hex) {
    const linear = hexToLinearRgb(hex);
    const luminance = 0.2126 * linear.r + 0.7152 * linear.g + 0.0722 * linear.b;
    return luminance > 0.55 ? '#1F1F1F' : '#FFFFFF';
}

function hexToCssColor(hex) {
    if (typeof hex !== 'string') {
        return hex;
    }
    const clean = hex.replace('#', '').toUpperCase();
    if (clean.length === 6) {
        return `#${clean}`;
    }
    if (clean.length === 8) {
        const r = parseInt(clean.slice(0, 2), 16);
        const g = parseInt(clean.slice(2, 4), 16);
        const b = parseInt(clean.slice(4, 6), 16);
        const a = parseInt(clean.slice(6, 8), 16) / 255;
        return `rgba(${r}, ${g}, ${b}, ${a.toFixed(3)})`;
    }
    return hex;
}

function sanitizeStringSetting(value, fallback) {
    if (typeof value === 'string' && value.trim().length > 0) {
        return value.trim();
    }
    return fallback;
}

function sanitizeTemplateSetting(value, fallback) {
    const sanitized = sanitizeStringSetting(value, fallback);
    if (!sanitized.includes('{count}')) {
        return fallback;
    }
    return sanitized;
}

function sanitizeProviderSetting(value) {
    if (typeof value === 'string') {
        const normalized = value.trim().toLowerCase();
        if (MODEL_FAMILIES[normalized]) {
            return normalized;
        }
    }
    return 'openai';
}

let highlightColors = {
    even: DEFAULT_EVEN_COLOR,
    odd: DEFAULT_ODD_COLOR
};

let statusBarTemplate = DEFAULT_STATUS_TEMPLATE;

function loadHighlightColors(context) {
    let evenStored = context.globalState.get(HIGHLIGHT_EVEN_KEY);
    let oddStored = context.globalState.get(HIGHLIGHT_ODD_KEY);

    const config = vscode.workspace.getConfiguration(CONFIG_SECTION);

    if (!evenStored) {
        const legacyEven = config.get('highlightEvenColor');
        if (legacyEven) {
            evenStored = sanitizeColorSetting(legacyEven, DEFAULT_EVEN_COLOR);
            void context.globalState.update(HIGHLIGHT_EVEN_KEY, evenStored);
        }
    }

    if (!oddStored) {
        const legacyOdd = config.get('highlightOddColor');
        if (legacyOdd) {
            oddStored = sanitizeColorSetting(legacyOdd, DEFAULT_ODD_COLOR);
            void context.globalState.update(HIGHLIGHT_ODD_KEY, oddStored);
        }
    }

    const resolvedEven = evenStored !== undefined && evenStored !== null ? evenStored : DEFAULT_EVEN_COLOR;
    const resolvedOdd = oddStored !== undefined && oddStored !== null ? oddStored : DEFAULT_ODD_COLOR;

    highlightColors = {
        even: sanitizeColorSetting(resolvedEven, DEFAULT_EVEN_COLOR),
        odd: sanitizeColorSetting(resolvedOdd, DEFAULT_ODD_COLOR)
    };
}

function loadStatusBarConfig() {
    const config = vscode.workspace.getConfiguration(CONFIG_SECTION);
    statusBarTemplate = sanitizeTemplateSetting(config.get('statusBarDisplayTemplate'), DEFAULT_STATUS_TEMPLATE);
}

function isHighlightableEditor(editor) {
    if (!editor || !editor.document) {
        return false;
    }
    const scheme = editor.document.uri.scheme || '';
    if (NON_HIGHLIGHT_SCHEMES.has(scheme)) {
        return false;
    }
    return true;
}

function getDefaultProviderFromConfig() {
    const config = vscode.workspace.getConfiguration(CONFIG_SECTION);
    return sanitizeProviderSetting(config.get('defaultModelFamily'));
}

function applyStatusBarTemplate(template, data) {
    const toText = (value) => {
        return value === undefined || value === null ? '' : `${value}`;
    };

    const replacements = {
        count: toText(data.count),
        family: toText(data.family),
        provider: toText(data.provider),
        model: toText(data.family),
        label: 'Token Count'
    };

    const baseTemplate = template || DEFAULT_STATUS_TEMPLATE;
    return baseTemplate.replace(/\{(count|family|provider|model|label)\}/g, (_match, key) => {
        return replacements[key] !== undefined ? replacements[key] : '';
    });
}

function createTokenDecorationTypes() {
    const evenTextColor = textColorForBackground(highlightColors.even);
    const oddTextColor = textColorForBackground(highlightColors.odd);
    const evenBackground = hexToCssColor(highlightColors.even);
    const oddBackground = hexToCssColor(highlightColors.odd);

    return {
        even: vscode.window.createTextEditorDecorationType({
            backgroundColor: evenBackground,
            color: evenTextColor
        }),
        odd: vscode.window.createTextEditorDecorationType({
            backgroundColor: oddBackground,
            color: oddTextColor
        })
    };
}

let tokenizerState = {
    encoder: null,
    supportsHighlight: false,
    requiresNormalization: false
};

/**
 * @param {vscode.ExtensionContext} context
 */
function activate(context) {
    if (typeof context.globalState.setKeysForSync === 'function') {
        context.globalState.setKeysForSync([HIGHLIGHT_EVEN_KEY, HIGHLIGHT_ODD_KEY]);
    }

    loadHighlightColors(context);
    loadStatusBarConfig();

    const statusBar = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Right, 100);
    statusBar.command = 'gpt-token-counter-live.changeModel';
    statusBar.name = 'LLM Token Counter';

    const highlightStatusBar = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Right, 99);
    highlightStatusBar.command = 'gpt-token-counter-live.toggleHighlight';
    highlightStatusBar.name = 'Token Highlight Toggle';
    highlightStatusBar.accessibilityInformation = {
        label: 'Toggle token highlighting',
        role: 'button'
    };
    highlightStatusBar.text = '$(symbol-color)';

    let tokenDecorations = createTokenDecorationTypes();

    let currentProvider = getDefaultProviderFromConfig();
    let currentFamilyName = MODEL_FAMILIES[currentProvider];
    let highlightEnabled = false;

    context.subscriptions.push(
        statusBar,
        highlightStatusBar,
        {
            dispose: () => {
                if (tokenDecorations) {
                    tokenDecorations.even.dispose();
                    tokenDecorations.odd.dispose();
                }
            }
        }
    );

    function clearTokenHighlights(editor = vscode.window.activeTextEditor) {
        if (!editor) {
            return;
        }
        if (tokenDecorations) {
            editor.setDecorations(tokenDecorations.even, []);
            editor.setDecorations(tokenDecorations.odd, []);
        }
    }

    function refreshTokenDecorations() {
        if (!tokenDecorations) {
            tokenDecorations = createTokenDecorationTypes();
            return;
        }

        vscode.window.visibleTextEditors.forEach(editor => {
            editor.setDecorations(tokenDecorations.even, []);
            editor.setDecorations(tokenDecorations.odd, []);
        });

        tokenDecorations.even.dispose();
        tokenDecorations.odd.dispose();
        tokenDecorations = createTokenDecorationTypes();

        vscode.window.visibleTextEditors.forEach(editor => {
            if (!isHighlightableEditor(editor)) {
                editor.setDecorations(tokenDecorations.even, []);
                editor.setDecorations(tokenDecorations.odd, []);
            }
        });
    }

    function applyTokenHighlights(editor, sourceText, baseOffset, tokenizationResult) {
        if (!highlightEnabled || !tokenizerState.supportsHighlight || !tokenizationResult || !editor) {
            clearTokenHighlights(editor);
            return;
        }

        if (!isHighlightableEditor(editor)) {
            clearTokenHighlights(editor);
            return;
        }

        if (tokenizationResult.normalizationChanged) {
            clearTokenHighlights(editor);
            return;
        }

        const evenRanges = [];
        const oddRanges = [];
        const tokens = tokenizationResult.tokens;
        const document = editor.document;
        let searchIndex = 0;

        for (let i = 0; i < tokens.length; i++) {
            const tokenId = tokens[i];
            let tokenString;

            try {
                const bytes = tokenizerState.encoder.decode_single_token_bytes(tokenId);
                tokenString = utf8Decoder.decode(bytes);
            } catch (error) {
                clearTokenHighlights(editor);
                return;
            }

            if (!tokenString.length) {
                continue;
            }

            const foundIndex = sourceText.indexOf(tokenString, searchIndex);
            if (foundIndex === -1) {
                clearTokenHighlights(editor);
                return;
            }

            const startOffset = baseOffset + foundIndex;
            const endOffset = startOffset + tokenString.length;
            const range = new vscode.Range(document.positionAt(startOffset), document.positionAt(endOffset));

            if (i % 2 === 0) {
                evenRanges.push(range);
            } else {
                oddRanges.push(range);
            }

            searchIndex = foundIndex + tokenString.length;
        }

        editor.setDecorations(tokenDecorations.even, evenRanges);
        editor.setDecorations(tokenDecorations.odd, oddRanges);
    }

    function updateHighlightStatusBar() {
        const activeEditor = vscode.window.activeTextEditor;
        if (!activeEditor || !isHighlightableEditor(activeEditor)) {
            highlightStatusBar.hide();
            return;
        }

        const activeForeground = new vscode.ThemeColor('statusBarItem.activeForeground');
        const activeBackground = new vscode.ThemeColor('statusBarItem.activeBackground');
        const inactiveForeground = new vscode.ThemeColor('statusBarItem.inactiveForeground');
        const inactiveBackground = new vscode.ThemeColor('statusBarItem.inactiveBackground');
        const defaultForeground = new vscode.ThemeColor('statusBarItem.foreground');
        const defaultBackground = new vscode.ThemeColor('statusBarItem.background');
        const unavailableForeground = new vscode.ThemeColor('statusBarItem.errorForeground');
        const unavailableBackground = new vscode.ThemeColor('statusBarItem.errorBackground');

        const applyAppearance = (text, tooltip, foreground, background, fallbackForeground, fallbackBackground) => {
            highlightStatusBar.text = text;
            highlightStatusBar.tooltip = tooltip;
            highlightStatusBar.color = foreground || fallbackForeground;
            highlightStatusBar.backgroundColor = background || fallbackBackground;
        };

        if (!tokenizerState.supportsHighlight) {
            applyAppearance(
                '$(circle-slash)',
                'Token highlighting is unavailable for this model family.',
                unavailableForeground,
                unavailableBackground,
                defaultForeground,
                defaultBackground
            );
        } else if (highlightEnabled) {
            applyAppearance(
                '$(paintcan)',
                'Click to disable token highlighting.',
                activeForeground || new vscode.ThemeColor('statusBarItem.prominentForeground'),
                activeBackground || new vscode.ThemeColor('statusBarItem.prominentBackground'),
                defaultForeground,
                defaultBackground
            );
        } else {
            applyAppearance(
                '$(symbol-color)',
                'Click to enable token highlighting.',
                inactiveForeground,
                inactiveBackground,
                defaultForeground,
                undefined
            );
        }
        highlightStatusBar.show();
    }

    function resetTokenizerState() {
        if (tokenizerState.encoder) {
            tokenizerState.encoder.free();
        }
        tokenizerState = {
            encoder: null,
            supportsHighlight: false,
            requiresNormalization: false
        };
    }

    // Function to initialize the encoder for the selected family
    function initializeEncoderForFamily(provider) {
        resetTokenizerState();

        if (provider === 'openai') {
            let candidate = null;
            try {
                candidate = encoding_for_model('gpt-5');
            } catch (e0) {
                try {
                    candidate = get_encoding('o200k_base');
                } catch (e1) {
                    try {
                        candidate = get_encoding('cl100k_base');
                    } catch (e2) {
                        candidate = null;
                    }
                }
            }
            if (candidate) {
                tokenizerState.encoder = candidate;
                tokenizerState.supportsHighlight = true;
            }
        } else if (provider === 'anthropic') {
            try {
                tokenizerState.encoder = getTokenizer();
                tokenizerState.supportsHighlight = true;
                tokenizerState.requiresNormalization = true;
            } catch (error) {
                vscode.window.showErrorMessage(`Failed to initialize tokenizer for ${MODEL_FAMILIES[provider]}: ${error.message}`);
            }
        } else if (provider === 'gemini') {
            try {
                tokenizerState.encoder = get_encoding('o200k_base');
            } catch (e1) {
                try {
                    tokenizerState.encoder = get_encoding('cl100k_base');
                } catch (e2) {
                    tokenizerState.encoder = null;
                }
            }
        }

        if (highlightEnabled && !tokenizerState.supportsHighlight) {
            highlightEnabled = false;
            clearTokenHighlights();
            vscode.window.showInformationMessage('Token highlighting disabled because the selected model family does not expose precise token boundaries.');
        }

        updateHighlightStatusBar();
    }

    function fallbackTokenCount(text) {
        if (currentProvider === 'anthropic') {
            return countTokens(text);
        }

        if (currentProvider === 'gemini') {
            // Fallback approximation: ~4 characters per token
            return Math.ceil(text.length / 4);
        }

        return text.length;
    }

    function computeTokenization(text) {
        if (tokenizerState.encoder) {
            const processedText = tokenizerState.requiresNormalization ? text.normalize('NFKC') : text;
            try {
                const encoded = tokenizerState.encoder.encode(processedText, 'all');
                return {
                    tokenCount: encoded.length,
                    tokenizationResult: {
                        tokens: encoded,
                        normalizationChanged: tokenizerState.requiresNormalization && processedText !== text
                    }
                };
            } catch (error) {
                vscode.window.showErrorMessage(`Tokenization failed for ${currentFamilyName}: ${error.message}`);
                resetTokenizerState();
            }
        }

        return {
            tokenCount: fallbackTokenCount(text),
            tokenizationResult: null
        };
    }

    let updateTokenCount = () => {
        const editor = vscode.window.activeTextEditor;
        if (!editor) {
            statusBar.hide();
            clearTokenHighlights();
            updateHighlightStatusBar();
            return; // No open text editor
        }

        if (!isHighlightableEditor(editor)) {
            statusBar.hide();
            clearTokenHighlights(editor);
            highlightStatusBar.hide();
            return;
        }

        const document = editor.document;
        const selection = editor.selection;
        const text = selection.isEmpty ? document.getText() : document.getText(selection);
        const baseOffset = selection.isEmpty ? 0 : document.offsetAt(selection.start);

        const { tokenCount, tokenizationResult } = computeTokenization(text);

        statusBar.text = applyStatusBarTemplate(statusBarTemplate, {
            count: tokenCount,
            family: currentFamilyName,
            provider: currentProvider
        });
        statusBar.show();

        if (highlightEnabled) {
            if (!tokenizerState.supportsHighlight) {
                clearTokenHighlights(editor);
            } else {
                applyTokenHighlights(editor, text, baseOffset, tokenizationResult);
            }
        } else {
            clearTokenHighlights(editor);
        }

        updateHighlightStatusBar();
    };

    vscode.window.onDidChangeTextEditorSelection(updateTokenCount, null, context.subscriptions);
    vscode.window.onDidChangeActiveTextEditor(updateTokenCount, null, context.subscriptions);
    vscode.workspace.onDidChangeTextDocument(updateTokenCount, null, context.subscriptions);

    vscode.workspace.onDidChangeConfiguration(event => {
        if (event.affectsConfiguration(`${CONFIG_SECTION}.statusBarDisplayTemplate`)) {
            loadStatusBarConfig();
            updateTokenCount();
        }

        if (event.affectsConfiguration(`${CONFIG_SECTION}.defaultModelFamily`)) {
            const desiredProvider = getDefaultProviderFromConfig();
            if (desiredProvider !== currentProvider) {
                currentProvider = desiredProvider;
                currentFamilyName = MODEL_FAMILIES[currentProvider];

                try {
                    initializeEncoderForFamily(currentProvider);
                } catch (error) {
                    vscode.window.showErrorMessage(`Failed to initialize tokenizer for ${currentFamilyName}: ${error.message}`);
                }

                updateTokenCount();
            }
        }
    }, null, context.subscriptions);

    let disposable = vscode.commands.registerCommand('gpt-token-counter-live.changeModel', async function () {
        /** @type {(vscode.QuickPickItem & { provider?: string, family?: string, command?: string })[]} */
        const familyItems = Object.entries(MODEL_FAMILIES).map(([provider, family]) => ({
            label: `${family} (${provider})`,
            description: provider === currentProvider ? 'Currently active' : undefined,
            detail: provider === 'gemini'
                ? 'Approximate tokenizer (highlighting unavailable)'
                : 'Precise tokenizer with highlighting',
            provider,
            family,
            picked: provider === currentProvider
        }));

        /** @type {(vscode.QuickPickItem & { provider?: string, family?: string, command?: string })[]} */
        const quickPickItems = [
            ...familyItems,
            { kind: vscode.QuickPickItemKind.Separator },
            {
                label: '$(gear) Extension Settings',
                description: 'Open settings for Live LLM Token Counter',
                command: 'openSettings'
            },
            {
                label: '$(symbol-color) Configure Highlight Colors',
                description: 'Preview and adjust token highlight colors',
                command: 'configureHighlights'
            }
        ];

        const selection = await vscode.window.showQuickPick(quickPickItems, {
            placeHolder: 'Select a Model Family',
            matchOnDescription: true,
            ignoreFocusOut: true
        });

        if (selection) {
            if (selection.command === 'openSettings') {
                await vscode.commands.executeCommand('workbench.action.openSettings', CONFIG_SECTION);
                return;
            } else if (selection.command === 'configureHighlights') {
                await vscode.commands.executeCommand('gpt-token-counter-live.configureHighlights');
                return;
            }

            if (selection.provider && MODEL_FAMILIES[selection.provider]) {
                currentProvider = selection.provider;
                currentFamilyName = MODEL_FAMILIES[currentProvider];

                try {
                    initializeEncoderForFamily(currentProvider);
                } catch (error) {
                    vscode.window.showErrorMessage(`Failed to initialize tokenizer for ${currentFamilyName}: ${error.message}`);
                    // Continue with approximation where applicable
                }

                updateTokenCount();
            }
        }
    });

    context.subscriptions.push(disposable);

    const toggleHighlight = vscode.commands.registerCommand('gpt-token-counter-live.toggleHighlight', () => {
        const nextState = !highlightEnabled;

        if (nextState && !tokenizerState.supportsHighlight) {
            vscode.window.showInformationMessage('Token highlighting is only available for GPT and Claude tokenizers.');
            highlightEnabled = false;
        } else if (nextState) {
            highlightEnabled = true;
            vscode.window.showInformationMessage('Token highlighting enabled.');
        } else {
            highlightEnabled = false;
            vscode.window.showInformationMessage('Token highlighting disabled.');
        }

        if (!highlightEnabled) {
            clearTokenHighlights();
        }

        updateHighlightStatusBar();
        updateTokenCount();
    });

    context.subscriptions.push(toggleHighlight);

    const configureHighlights = vscode.commands.registerCommand('gpt-token-counter-live.configureHighlights', async () => {
        loadHighlightColors(context);
        const panel = vscode.window.createWebviewPanel(
            'tokenHighlightConfigurator',
            'Token Highlight Colors',
            vscode.ViewColumn.Active,
            {
                enableScripts: true,
                retainContextWhenHidden: true
            }
        );

        const evenColor = sanitizeColorSetting(highlightColors.even, DEFAULT_EVEN_COLOR);
        const oddColor = sanitizeColorSetting(highlightColors.odd, DEFAULT_ODD_COLOR);

        const getHtml = (evenHex, oddHex) => {
            const nonce = String(Date.now());
            const toHex6 = (value) => {
                if (typeof value !== 'string') {
                    return '#000000';
                }
                const match = value.match(/^#([0-9A-Fa-f]{6})/);
                if (match) {
                    return `#${match[1].toUpperCase()}`;
                }
                const match8 = value.match(/^#([0-9A-Fa-f]{6})([0-9A-Fa-f]{2})$/);
                if (match8) {
                    return `#${match8[1].toUpperCase()}`;
                }
                return '#000000';
            };

            const toAlphaPercent = (value) => {
                if (typeof value === 'string' && value.length === 9) {
                    return Math.round((parseInt(value.slice(7), 16) / 255) * 100);
                }
                return 100;
            };

            const normalizedEven = sanitizeColorSetting(evenHex, DEFAULT_EVEN_COLOR);
            const normalizedOdd = sanitizeColorSetting(oddHex, DEFAULT_ODD_COLOR);

            const evenHex6 = toHex6(normalizedEven).toLowerCase();
            const oddHex6 = toHex6(normalizedOdd).toLowerCase();
            const evenAlphaPercent = toAlphaPercent(normalizedEven);
            const oddAlphaPercent = toAlphaPercent(normalizedOdd);

            return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="Content-Security-Policy" content="default-src 'none'; style-src 'unsafe-inline'; script-src 'nonce-${nonce}';">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Token Highlight Colors</title>
    <style>
        body { font-family: var(--vscode-font-family); padding: 16px; color: var(--vscode-foreground); background: var(--vscode-editor-background); }
        h1 { font-size: 1.3rem; margin-bottom: 12px; }
        .picker-row { display: flex; align-items: center; gap: 12px; margin-bottom: 6px; }
        .alpha-row { display: flex; align-items: center; gap: 12px; margin-bottom: 16px; }
        label { width: 160px; }
        input[type="color"] { appearance: none; -webkit-appearance: none; border: none; width: 40px; height: 40px; padding: 0; background: transparent; cursor: pointer; }
        input[type="range"] { flex: 1; }
        .preview { border: 1px solid var(--vscode-editorWidget-border); border-radius: 6px; overflow: hidden; box-shadow: 0 1px 3px rgba(0,0,0,0.15); }
        .preview-row { display: flex; height: 36px; }
        .preview-row span { flex: 1; display: flex; align-items: center; justify-content: center; font-size: 0.85rem; }
        .helper { font-size: 0.85rem; color: var(--vscode-descriptionForeground); margin-bottom: 16px; }
        .value-chip { font-family: var(--vscode-editor-font-family, monospace); font-size: 0.8rem; padding: 4px 6px; border-radius: 4px; background: var(--vscode-editorWidget-background); border: 1px solid var(--vscode-editorWidget-border); min-width: 110px; text-align: center; }
    </style>
</head>
<body>
    <h1>Token Highlight Colors</h1>
    <p class="helper">Adjust the alternating highlight colors. Changes apply immediately.</p>
    <div class="picker-row">
        <label for="evenColor">Even tokens</label>
        <input id="evenColor" type="color" value="${evenHex6}" aria-label="Even token color" />
        <span class="value-chip" data-label="even">${normalizedEven}</span>
    </div>
    <div class="alpha-row">
        <label for="evenAlpha">Opacity</label>
        <input id="evenAlpha" type="range" min="0" max="100" step="1" value="${evenAlphaPercent}" aria-label="Even token opacity" />
        <span class="value-chip" data-alpha-label="even">${evenAlphaPercent}%</span>
    </div>
    <div class="picker-row">
        <label for="oddColor">Odd tokens</label>
        <input id="oddColor" type="color" value="${oddHex6}" aria-label="Odd token color" />
        <span class="value-chip" data-label="odd">${normalizedOdd}</span>
    </div>
    <div class="alpha-row">
        <label for="oddAlpha">Opacity</label>
        <input id="oddAlpha" type="range" min="0" max="100" step="1" value="${oddAlphaPercent}" aria-label="Odd token opacity" />
        <span class="value-chip" data-alpha-label="odd">${oddAlphaPercent}%</span>
    </div>
    <div class="preview">
        <div class="preview-row" data-preview="even" style="background:${normalizedEven}; color: var(--vscode-editor-foreground);">
            <span>Even token preview (${normalizedEven})</span>
        </div>
        <div class="preview-row" data-preview="odd" style="background:${normalizedOdd}; color: var(--vscode-editor-foreground);">
            <span>Odd token preview (${normalizedOdd})</span>
        </div>
    </div>
    <script nonce="${nonce}">
        const vscodeApi = acquireVsCodeApi();
        const evenInput = document.getElementById('evenColor');
        const oddInput = document.getElementById('oddColor');
        const evenAlpha = document.getElementById('evenAlpha');
        const oddAlpha = document.getElementById('oddAlpha');
        const evenPreview = document.querySelector('[data-preview="even"]');
        const oddPreview = document.querySelector('[data-preview="odd"]');
        const evenLabel = document.querySelector('[data-label="even"]');
        const oddLabel = document.querySelector('[data-label="odd"]');
        const evenAlphaLabel = document.querySelector('[data-alpha-label="even"]');
        const oddAlphaLabel = document.querySelector('[data-alpha-label="odd"]');

        const clamp = (value, min, max) => Math.min(Math.max(value, min), max);

        const parseRgb = (value) => {
            const match = value.match(/rgba?\((\d+),\s*(\d+),\s*(\d+)(?:,\s*([\d.]+))?\)/);
            if (!match) {
                return { r: 30, g: 30, b: 30, a: 1 };
            }
            return {
                r: Number(match[1]),
                g: Number(match[2]),
                b: Number(match[3]),
                a: match[4] !== undefined ? Number(match[4]) : 1
            };
        };

        const bodyBackground = parseRgb(getComputedStyle(document.body).backgroundColor);

        const hexToRgba = (hex) => {
            const clean = hex.replace('#', '');
            const r = parseInt(clean.slice(0, 2), 16);
            const g = parseInt(clean.slice(2, 4), 16);
            const b = parseInt(clean.slice(4, 6), 16);
            const a = clean.length === 8 ? parseInt(clean.slice(6, 8), 16) / 255 : 1;
            return { r, g, b, a };
        };

        const blendOnBackground = (foreground) => {
            const alpha = clamp(foreground.a, 0, 1);
            return {
                r: Math.round(foreground.r * alpha + bodyBackground.r * (1 - alpha)),
                g: Math.round(foreground.g * alpha + bodyBackground.g * (1 - alpha)),
                b: Math.round(foreground.b * alpha + bodyBackground.b * (1 - alpha))
            };
        };

        const luminance = ({ r, g, b }) => {
            const srgb = [r, g, b].map(component => {
                const scaled = component / 255;
                return scaled <= 0.03928 ? scaled / 12.92 : Math.pow((scaled + 0.055) / 1.055, 2.4);
            });
            return 0.2126 * srgb[0] + 0.7152 * srgb[1] + 0.0722 * srgb[2];
        };

        const pickTextColor = (hex) => {
            const fg = hexToRgba(hex);
            const blended = blendOnBackground(fg);
            const lum = luminance(blended);
            return lum > 0.55 ? '#1f1f1f' : '#f5f5f5';
        };

        const toAlphaHex = (percent) => {
            const clampedPercent = clamp(Number(percent), 0, 100);
            const alphaValue = Math.round((clampedPercent / 100) * 255);
            return alphaValue >= 255 ? '' : alphaValue.toString(16).padStart(2, '0').toUpperCase();
        };

        const composeColor = (baseHex, percent) => {
            const normalizedBase = baseHex.toUpperCase();
            const alphaHex = toAlphaHex(percent);
            return alphaHex ? normalizedBase + alphaHex : normalizedBase;
        };

        const updateAlphaLabel = (element, percent) => {
            if (element) {
                element.textContent = Math.round(percent) + '%';
            }
        };

        const applyPreview = (element, value, label) => {
            if (!element) {
                return;
            }
            element.style.background = value;
            const textColor = pickTextColor(value.toUpperCase());
            element.style.color = textColor;
            const span = element.querySelector('span');
            if (span) {
                span.textContent = label + ' preview (' + value.toUpperCase() + ')';
                span.style.color = textColor;
            }
        };

        const applyLabel = (element, value) => {
            if (element) {
                element.textContent = value.toUpperCase();
            }
        };

        const emitColorChange = (key, value) => {
            vscodeApi.postMessage({ type: 'colorChange', key, value });
        };

        const updateEven = () => {
            const percent = Number(evenAlpha.value);
            updateAlphaLabel(evenAlphaLabel, percent);
            const finalHex = composeColor(evenInput.value, percent);
            applyPreview(evenPreview, finalHex, 'Even token');
            applyLabel(evenLabel, finalHex);
            emitColorChange('${HIGHLIGHT_EVEN_KEY}', finalHex);
        };

        const updateOdd = () => {
            const percent = Number(oddAlpha.value);
            updateAlphaLabel(oddAlphaLabel, percent);
            const finalHex = composeColor(oddInput.value, percent);
            applyPreview(oddPreview, finalHex, 'Odd token');
            applyLabel(oddLabel, finalHex);
            emitColorChange('${HIGHLIGHT_ODD_KEY}', finalHex);
        };

        evenInput.addEventListener('input', updateEven);
        evenAlpha.addEventListener('input', updateEven);
        oddInput.addEventListener('input', updateOdd);
        oddAlpha.addEventListener('input', updateOdd);

        window.addEventListener('message', (event) => {
            const message = event.data;
            if (message && message.type === 'colorUpdate') {
                const applyIncomingColor = (input, alphaControl, alphaLabelEl, previewEl, labelEl, labelText) => {
                    const upperValue = message.value.toUpperCase();
                    const base = upperValue.slice(0, 7);
                    const percent = upperValue.length === 9 ? Math.round((parseInt(upperValue.slice(7), 16) / 255) * 100) : 100;
                    input.value = base.toLowerCase();
                    alphaControl.value = percent;
                    updateAlphaLabel(alphaLabelEl, percent);
                    applyPreview(previewEl, upperValue, labelText);
                    applyLabel(labelEl, upperValue);
                };

                if (message.key === '${HIGHLIGHT_EVEN_KEY}') {
                    applyIncomingColor(evenInput, evenAlpha, evenAlphaLabel, evenPreview, evenLabel, 'Even token');
                } else if (message.key === '${HIGHLIGHT_ODD_KEY}') {
                    applyIncomingColor(oddInput, oddAlpha, oddAlphaLabel, oddPreview, oddLabel, 'Odd token');
                }
            }
        });

        const initialize = () => {
            updateAlphaLabel(evenAlphaLabel, Number(evenAlpha.value));
            updateAlphaLabel(oddAlphaLabel, Number(oddAlpha.value));

            const initialEvenHex = composeColor(evenInput.value, Number(evenAlpha.value));
            const initialOddHex = composeColor(oddInput.value, Number(oddAlpha.value));

            applyPreview(evenPreview, initialEvenHex, 'Even token');
            applyPreview(oddPreview, initialOddHex, 'Odd token');
            applyLabel(evenLabel, initialEvenHex);
            applyLabel(oddLabel, initialOddHex);
        };

        initialize();
    </script>
</body>
</html>`;
        };

        panel.webview.html = getHtml(evenColor, oddColor);

        const updateColorSetting = async (key, value) => {
            const targetKey = key === HIGHLIGHT_EVEN_KEY ? 'even' : key === HIGHLIGHT_ODD_KEY ? 'odd' : null;
            if (!targetKey) {
                return;
            }

            const hexValue = sanitizeColorSetting(value, targetKey === 'even' ? DEFAULT_EVEN_COLOR : DEFAULT_ODD_COLOR);
            highlightColors[targetKey] = hexValue;
            await context.globalState.update(key, hexValue);
            loadHighlightColors(context);

            refreshTokenDecorations();
            updateTokenCount();
            panel.webview.postMessage({ type: 'colorUpdate', key, value: hexValue });
        };

        const messageDisposable = panel.webview.onDidReceiveMessage(async (message) => {
            if (message.type === 'colorChange' && message.key) {
                await updateColorSetting(message.key, message.value);
            }
        });

        panel.onDidDispose(() => {
            messageDisposable.dispose();
        }, null, context.subscriptions);
    });

    context.subscriptions.push(configureHighlights);

    // Initial update
    initializeEncoderForFamily(currentProvider);
    updateTokenCount();
    updateHighlightStatusBar();
}

function deactivate() {
    if (tokenizerState.encoder) {
        tokenizerState.encoder.free();
        tokenizerState.encoder = null;
    }
}

module.exports = {
    activate,
    deactivate
}