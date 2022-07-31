const express = require('express');
const prettier = require("prettier");

const app = express();

app.use(
  express.urlencoded({
    extended: true,
  }),
);
app.use(express.json());

app.get('/', (req, res) => {
	var info = prettier.getSupportInfo();
	res.write(JSON.stringify(info));
	res.end();
});

app.post('/', (req, res) => {
	// TODO, add debug level (env) and print
	// console.log(req.body);
	//
	var source = req.body.source;
	var config = req.body.config;

	try {
		var fmt = prettier.format(source, config);
		res.write(fmt);
		res.end();
	} catch(error) {
		console.log(error)
		res.status(400).send(error.toString());
	}
});

var port = 3000;
var PORT = process.env.PORT
if (PORT) {
	port = PORT
}

app.listen(port, ()=> {
	console.log(`listening on ${port}`)
})
