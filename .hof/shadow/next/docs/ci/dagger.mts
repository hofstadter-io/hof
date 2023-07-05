import util from 'util';
import { exec as origExec } from "child_process"
const exec = util.promisify(origExec)

import Client, { connect, Container, Directory } from "@dagger.io/dagger"

const registry = "us-central1-docker.pkg.dev/hof-io--develop/testing"
const version = "0.0.1"
const cluster = ""
const zone = ""

// initialize Dagger client
connect(
  async (client: Client) => {

		// get our source
    const source = client.host().directory(".", { exclude: [".next", "node_modules/"] })

		const image = makeImage(client, source)
		const nginx = makeNginx(client, source)

		// await exportImage(image, `docs:${version}`)
		await image.publish(`${registry}/docs-server:${version}`)
		await nginx.publish(`${registry}/docs-nginx:${version}`)

		var gcloud = await gcloudImage(client)

		const cuelm = source.file("./ci/k8s/cuelm.cue")

		deploy(gcloud, cluster, zone, cuelm)
			.sync()

  },
  { LogOutput: process.stderr }
)

function deploy(gcloud: Container, cluster: string, zone: string, cuelm: File) {
	return gcloud.withEnvVariable("CACHEBUST", Date.now().toString())
		.withExec(["gcloud", "container", "clusters", "get-credentials", cluster, "--zone", zone])
		.withWorkdir("/work")
		.withFile("/work/cuelm.cue", cuelm)
		.withExec(["hof", "mod", "init", "hof.io/deploy"])
		.withExec(["hof", "mod", "tidy"])
		.withExec(["cue", "export", "cuelm.cue", "-e", "Install", "-f", "-o", "cuelm.yaml", "-t", `version=${version}`])
		.withExec(["cat", "cuelm.yaml"])
		.withExec(["kubectl", "apply", "-f", "cuelm.yaml"])
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
	const ver = "v0.6.0-alpha.2"
	const url = `https://github.com/cue-lang/cue/releases/download/${ver}/cue_${ver}_linux_amd64.tar.gz`
	const targz = client.http(url)

	const cue = untargz(client, targz).file("cue")

	return container.withFile("/usr/local/bin/cue", cue)
}

function addHof(client: Client, container: Container) {
	const ver = "v0.6.8-rc.5"
	const url = `https://github.com/hofstadter-io/hof/releases/download/${ver}/hof_${ver}_Linux_x86_64`
	const hof = client.http(url)

	return container
		.withFile("/usr/local/bin/hof", hof)
		.withExec(["chmod", "+x", "/usr/local/bin/hof"])
}
