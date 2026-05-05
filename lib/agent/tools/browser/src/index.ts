import { chromium } from "playwright"

console.log("hallo!")

const browser = await chromium.launch({
  headless: false,
})
const context = await browser.newContext({
  userAgent: 'verdverm/veg/chromium',
})
const page = await context.newPage()
// await page.goto("https://verdverm.com")
await page.goto("https://news.ycombinator.com")

const title = await page.title()
console.log("title:", title)

const links = await page.$$eval('a', elements => elements.map(element => element.href));
console.log("links:", links)

// actual content
const c = await page.content()
console.log("content\n", c)

// cleanup before exit
browser.close()