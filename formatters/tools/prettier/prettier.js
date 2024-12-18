import express from 'express'
import * as prettier from 'prettier'

const app = express();

app.use(
  express.urlencoded({
    extended: true,
  }),
);
app.use(express.json({limit: "50mb", extended: true, parameterLimit:50000}));

app.get('/', async (req, res) => {
	var info = await prettier.getSupportInfo();
	res.write(JSON.stringify(info));
	res.end();
});

app.post('/', async (req, res) => {
	// TODO, add debug level (env) and print
	// console.log(req.body);
	//
	var source = req.body.source;
	var config = req.body.config;

	try {
		var fmt = await prettier.format(source, config);
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
