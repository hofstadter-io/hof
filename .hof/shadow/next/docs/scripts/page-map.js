// script.js

const fs = require('fs');
const path = require('path');

// Assuming your pages are in a directory named 'pages'
const pagesDirectory = path.join(__dirname, '..', 'src/app');

function getPages(directory, pages = []) {
  // Read all files in the current directory
  fs.readdirSync(directory).forEach(file => {
    const absolutePath = path.join(directory, file);
		const relativeToPages = absolutePath
			.replace(pagesDirectory + "/", "")
			.replace(/\\+/g, '/');

    // If the file is a directory, recursively get pages inside it
    if (fs.statSync(absolutePath).isDirectory()) {
      pages = getPages(absolutePath, pages);
    } else if (
      // Only add .js, .jsx, .ts, or .tsx files and ignore _app.js, _document.js, and _error.js
      /^(page.tsx)/.test(file)
    ) {

      pages.push(relativeToPages);
    }
  });

  return pages;
}

const pages = getPages(pagesDirectory);

console.log(pages);

pages.forEach(page => {
	var dir = path.dirname(page)
	if (dir === ".") {
		dir = ""
	}
	var Page = {
		fn: page,
		path: "/" + dir,
	}

	const mdx = path.join(dir, "content.mdx")
	const mfn = path.join(pagesDirectory, mdx)
	console.log(mfn)
	const content = require(Page.fn);
	console.log(content)
	fs.readFile(mfn, 'utf8', (err, data) => {
	  if (err) { 
			// ignore
	    console.error(err);
			console.log(Page)
	    // return;
	  }
		Page.mdx = mdx;
		Page.content = data;
		console.log(Page)

	  // console.log(data);
	});

})