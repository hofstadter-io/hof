package k8s

import (
	"github.com/hofstadter-io/cuelm/schema"
)

Install: schema.#List & {
	items: [
		#Site.Ingress,
		#Site.Service,
		#Site.Deployment,
	]
}

Update: schema.#List & {
	items: [
		#Site.Deployment,
	]
}

#Site: {
	_Values: {
		name:      string | *"docs" @tag(name)
		namespace: string           @tag(namespace)

		version:  string | *"manual" @tag(version)
		registry: string             @tag(registry)
		domain:   string             @tag(domain)

		ga_mp_apikey: string | *"" @tag(ga_mp_apikey)

		port: {
			nginx:    80
			server:   3000
			postgres: 5432
		}

		resources: {
			requests: {cpu: "100m", memory: "100Mi"}
			// limits: {cpu: "200m", memory: "200Mi"}
		}

		#metadata: {
			name:      _Values.name
			namespace: _Values.namespace
			labels: {
				app: _Values.name
			}
			...
		}
	}

	Ingress: schema.#Ingress & {
		metadata: _Values.#metadata & {
			annotations: {
				"kubernetes.io/tls-acme":                         "true"
				"kubernetes.io/ingress.class":                    "nginx"
				"nginx.ingress.kubernetes.io/force-ssl-redirect": "true"
				"cert-manager.io/cluster-issuer":                 "letsencrypt-prod"
				"cert-manager.io/issue-temporary-certificate":    "true"
				"acme.cert-manager.io/http01-edit-in-place":      "true"
			}
		} // END Ingress.metadata

		spec: {
			tls: [{
				hosts: [_Values.domain]
				secretName: "\(_Values.name)-tls"
			}]

			rules: [{
				host: _Values.domain
				http: paths: [{
					backend: {
						service: {
							name: Service.metadata.name
							port: "number": Service.spec.ports[0].port
						}
					}
				}]
			}]

		} // END Ingress.spec
	} // END Ingress

	Service: schema.#Service & {
		metadata: _Values.#metadata
		spec: {
			selector: _Values.#metadata.labels
			type:     "NodePort"
			ports: [{
				port:       _Values.port.nginx
				targetPort: _Values.port.nginx
			}]
		}
	}

	Deployment: schema.#Deployment & {
		metadata: _Values.#metadata
		spec: {
			selector: matchLabels: _Values.#metadata.labels

			template: {
				metadata: labels: _Values.#metadata.labels
				spec: {
					containers: [{
						name:            "nginx"
						image:           "\(_Values.registry)/docs-nginx:\(_Values.version)"
						imagePullPolicy: "Always"
						env: [{
							name:  "GA_MP_APIKEY"
							value: _Values.ga_mp_apikey
						}]
						ports: [{
							containerPort: _Values.port.nginx
							protocol:      "TCP"
						}]
						resources: _Values.resources
						_Probes & {_port: _Values.port.nginx}
					}, {
						name:            "server"
						image:           "\(_Values.registry)/docs-server:\(_Values.version)"
						imagePullPolicy: "Always"
						ports: [{
							containerPort: _Values.port.server
							protocol:      "TCP"
						}]
						envFrom: [{
							secretRef: name: "studios"
						}]
						resources: _Values.resources
						_Probes & {_port: _Values.port.server}
					}, {
						name:            "postgres"
						image:           "postgres:14"
						imagePullPolicy: "Always"
						ports: [{
							containerPort: _Values.port.postgres
							protocol:      "TCP"
						}]
						env: [
							{name: "POSTGRES_USER", value:     "supacode"},
							{name: "POSTGRES_PASSWORD", value: "supacode"},
							{name: "POSTGRES_DB", value:       "supacode"},
						]
						resources: _Values.resources
						// _Probes & {_port: _Values.port.postgres}
					}]

					volumes: [{
						name: "env-file"
						secret: secretName: "studios"
					}]
				}
			}
		}
	}

	_Probes: {
		_port: int
		readinessProbe: {
			httpGet: port: _port
			initialDelaySeconds: 6
			failureThreshold:    3
			periodSeconds:       10
		}
		livenessProbe: {
			httpGet: port: _port
			initialDelaySeconds: 6
			failureThreshold:    3
			periodSeconds:       10
		}
	}

}
