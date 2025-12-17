# Browser Tool

TypeScript-based browser automation tool using [Playwright](https://playwright.dev/). This tool launches a Chromium browser instance to interact with websites, scrape content, or perform automated tasks.

## Key Files

- **`src/index.ts`**: The main entry point. Initializes the browser, context, and page.

## Usage

Run the development script:

```bash
npm run dev
```

## Implementation Details

The tool uses Playwright's `chromium` launcher.

```typescript
// src/index.ts
import { chromium } from "playwright"

const browser = await chromium.launch({
  headless: false, // Currently configured to run with UI
})
const context = await browser.newContext({
  userAgent: 'verdverm/veg/chromium',
})
const page = await context.newPage()

// Navigation and extraction
await page.goto("https://news.ycombinator.com")
const title = await page.title()
const content = await page.content()
```

## Tests / Benchmarks

- [Online-Mind2Web](https://github.com/OSU-NLP-Group/Online-Mind2Web)
