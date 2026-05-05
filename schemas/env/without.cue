package env

WithoutDefaultArgs: {
	$kind: "withoutDefaultArgs"
}

WithoutDirectory: {
	$kind: "withoutDirectory"
	path:  string
	expand: bool | *false
}

WithoutEntrypoint: {
	$kind: "withoutEntrypoint"
	keepDefaultArgs: bool | *false
}

WithoutEnvVariable: {
	$kind: "withoutEnvVariable"
	name:  string
}

WithoutExposedPort: {
	$kind: "withoutExposedPort"
	port:     int
	protocol: string | *"TCP"
}

WithoutFile: {
	$kind: "withoutFile"
	path:  string
	expand: bool | *false
}

WithoutFiles: {
	$kind: "withoutFiles"
	paths: [...string]
	expand: bool | *false
}

WithoutLabel: {
	$kind: "withoutLabel"
	name:  string
}

WithoutMount: {
	$kind: "withoutMount"
	path:  string
	expand: bool | *false
}

WithoutRegistryAuth: {
	$kind: "withoutRegistryAuth"
	address: string
}

WithoutSecretVariable: {
	$kind: "withoutSecretVariable"
	name:  string
}

WithoutUnixSocket: {
	$kind: "withoutUnixSocket"
	path:  string
	expand: bool | *false
}

WithoutUser: {
	$kind: "withoutUser"
}

WithoutWorkdir: {
	$kind: "withoutWorkdir"
}
