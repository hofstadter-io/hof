const shiki = require('shiki')
const fs = require('fs');
var data = fs.readFileSync(0, 'utf-8');

// https://miguelpiedrafita.com/vscode-highlighting
// https://github.com/shikijs/shiki/pull/236

// const t = shiki.loadTheme('./theme.json')

shiki
  .getHighlighter({
    theme: 'github-light'
  })
  .then(highlighter => {
    console.log(highlighter.codeToHtml(data, 'cue'))
  })
