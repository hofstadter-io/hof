import util from 'util';
import { exec as origExec } from "child_process";
const exec = util.promisify(origExec);

import { Command } from 'commander';

import Client, { connect, Container, Directory, File } from "@dagger.io/dagger"

// set defaults here
const registry = ""
const version = ""
const cluster = ""
const zone = ""

const cue_version = "v0.6.0-alpha.2"
const hof_version = "v0.6.8-rc.5"

// initialize command interface
const cli = new Command();
// override the default command version reporting flag
cli.version('0.0.1', '-V, --script-version', 'output the current version of this script');

// set the flags
cli
	.option('-v, --version <value>', 'set the version to use for the command', version)
	.option('--registry <value>', 'set the container registry to use', registry)
  .option('--cluster <value>', 'set the gke k8s cluster name', cluster)
	.option('--zone <value>', 'set the gcloud zone', zone)

	// these defaults appear in the cuelm.cue
	.option('--domain <value>', 'the host domain for your application')
	.option('--name <value>', 'the name for the k8s resources')
	.option('--namespace <value>', 'the k8s namespace to use')
	.option('--dry-run', 'print the k8s resources yaml')

// set commands
cli
	.command('build')
  .description('build the images')
  .action(() => {
		run(build)
	})

cli
	.command('publish')
  .description('publish the images')
  .action(() => {
		run(publish)
	})

cli
	.command('deploy')
  .description('deploy the app')
  .action(() => {
		run(deploy)
	})

// parse args (run cli?)
cli.parse(process.argv);


function run(fn: any) {
	// initialize Dagger client
	connect(
		async (client: Client) => {

			// get our source
			const source: Directory = client.host().directory(".", { exclude: [".next", "node_modules/"] })

			await fn(client, source)

		},
		{ LogOutput: process.stderr }
	)
}

async function build(client: Client, source: Directory) {
	const image = makeImage(client, source)
	const nginx = makeNginx(client, source)
	image.sync()
	nginx.sync()
}

async function publish(client: Client, source: Directory) {
	const image = makeImage(client, source)
	const nginx = makeNginx(client, source)

	const opts = cli.opts()

	// await exportImage(image, `docs:${version}`)
	await image.publish(`${opts.registry}/docs-server:${opts.version}`)
	await nginx.publish(`${opts.registry}/docs-nginx:${opts.version}`)
}

async function deploy(client: Client, source: Directory) {
	const gcloud = await gcloudImage(client)
	const cuelm = source.file("./ci/k8s/cuelm.cue")

	const opts = cli.opts()

	var cuecmd = ["cue", "export", "cuelm.cue", "-e", "Install", "-f", "-o", "cuelm.yaml"]
	cuecmd.push("-t", `version=${opts.version}`)
	cuecmd.push("-t", `registry=${opts.registry}`)
	if (opts.domain) {
		cuecmd.push("-t", `domain=${opts.domain}`)
	}
	if (opts.name) {
		cuecmd.push("-t", `name=${opts.name}`)
	}
	if (opts.namespace) {
		cuecmd.push("-t", `namespace=${opts.namespace}`)
	}

	gcloud.withEnvVariable("CACHEBUST", Date.now().toString())
		.withExec(["gcloud", "container", "clusters", "get-credentials", opts.cluster, "--zone", opts.zone])
		.withWorkdir("/work")
		.withFile("/work/cuelm.cue", cuelm)
		.withExec(["hof", "mod", "init", "hof.io/deploy"])
		.withExec(["hof", "mod", "tidy"])
		.withExec(cuecmd)
		.withExec(["cat", "cuelm.yaml"])
		.withExec(["kubectl", "apply", "-f", "cuelm.yaml"])
		.sync()
}

function makeImage(client: Client, source: Directory) {
	// build up base image
	const base = client.container()
		.from("node:18-alpine").pipeline("base")
		.withExec(["apk", "add", "--no-cache", "libc6-compat"])
		.withEnvVariable("NEXT_TELEMETRY_DISABLED", "1")
		.withWorkdir("/app")

	// fetches dependencies
	const deps = base.pipeline("deps")
		.withDirectory("/app", source, { include: ["package.json", "package-lock.json"] })
		.withExec(["npm", "install", "--production"])
		.directory("/app/node_modules")

	// builds next production output
	const next = base.pipeline("build")
		.withDirectory("/app", source)
		.withDirectory("/app/node_modules", deps)
		.withExec(["npm", "run", "build"])
		.directory("/app/.next")

	const runner = base.pipeline("runner")
		// runner setup
		.withEnvVariable("NODE_ENV", "production")
		.withExec(["addgroup", "--system", "--gid", "1001", "nodejs"])
		.withExec(["adduser", "--system", "--uid", "1001", "nextjs"])

		// code and stuff
		.withDirectory("/app", source, { include: ["package.json", "package-lock.json"] })
		.withDirectory("/app/node_modules", deps)
		.withDirectory("/app/.next", next, { owner: "nextjs:nodejs" })

		// runtime settings
		.withUser("nextjs")
		.withEnvVariable("PORT", "3000")
		.withExposedPort(3000)
		.withEntrypoint(["sh", "-c", "npm start"])

	return runner
}

function makeNginx(client: Client, source: Directory) {
	return client.pipeline("nginx").container()
		.from("nginx:1.25")
		.withFile("/etc/nginx/nginx.conf", source.file("./ci/nginx/nginx.conf"))
		.withFile("/etc/nginx/templates/server.conf.template", source.file("./ci/nginx/server.conf.template"))		
}

async function exportImage(image: Container, name: string) {
	const tarfn = "dagger-export.tar"

	await image.pipeline("export").export(tarfn)

	try {
		const { stdout, stderr } = await exec(`docker load -i ${tarfn}`)
		console.log(stdout)
		console.log(stderr)

		const parts = stdout.split(" ");
		const hash = parts[parts.length-1].trim();

		const { stdout: tagout, stderr: tagerr } = await exec(`docker tag ${hash} ${name}`)
		console.log(tagout)
		console.log(tagerr)
		await exec(`rm -f ${tarfn}`)

	} catch(e) {
		console.error(e);
	}
}

async function gcloudImage(client: Client) {
	const { stdout } = await exec("gcloud info --format='value(config. paths. global_config_dir)'")
	const cfg = stdout.trim()

	const d = client.host().directory(cfg);

	var c = client.container()
		.from("google/cloud-sdk").pipeline("gcloud")
		// mount local user config, need service account in CI? or just a step to auth?
		.withEnvVariable("CLOUDSDK_CONFIG", "/gcloud/config")
		.withMountedDirectory("/gcloud/config", d)

	c = addCue(client, c)
	c = addHof(client, c)
	return c
}

function untargz(client: Client, targz: File) {
	return client.container()
		.from("google/cloud-sdk").pipeline("gcloud")
		.withWorkdir("/tmp")
		.withFile("/tmp/file.tar.gz", targz)
		.withExec(["tar", "-xf", "file.tar.gz"])
		.directory("/tmp")
}

function addCue(client: Client, container: Container) {
	const url = `https://github.com/cue-lang/cue/releases/download/${cue_version}/cue_${cue_version}_linux_amd64.tar.gz`
	const targz = client.http(url)

	const cue = untargz(client, targz).file("cue")

	return container.withFile("/usr/local/bin/cue", cue)
}

function addHof(client: Client, container: Container) {
	const url = `https://github.com/hofstadter-io/hof/releases/download/${hof_version}/hof_${hof_version}_Linux_x86_64`
	const hof = client.http(url)

	return container
		.withFile("/usr/local/bin/hof", hof)
		.withExec(["chmod", "+x", "/usr/local/bin/hof"])
}
