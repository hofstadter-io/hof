---
title: "Data Files"
weight: 50
---

{{<lead>}}
Hof can generate data files during code generation.
Create any JSON, Yaml, XML, Toml and CUE files (note, add all here)
from a single source of truth CUE configuration and schemas.
Use existing data and configurations as part of the input
if you already have a source of truth.
{{</lead>}}

## Ad hoc generation

In CUE you can output a single data file at a time:

```text
cue export config.yaml schema.cue -e app -O app.yaml
cue export config.yaml schema.cue -e helm -O values.yaml
```

However this can become unweildy for large numbers of files
or when you want to do more complex things with CUE.

With hof ad hoc generation, you can output multiple files in one command.

```text
hof gen config.yaml config.cue \
    -T :stg.app=app.yaml \
    -T :stg.helm=values.yaml
```

Here, we have generated the application configuration
and helm input values for the staging deployment of our application.
We use an existing Yaml file, apply a CUE schema,
set defaults, or otherwise validate or transform the inputs.

Input files:

- `config.yaml` - existing or minimal configuration
- `config.cue` - a CUE file where we can apply schemas or transformations

Output files:

- `-T <template>:<cue-path>=<outfile-name>` is a template, or a data file in this case
- the `-T :<...>` leaves off the template, telling hof this is a data file
- `<cue-path>` is any valid CUE path to subselect a field within the larger CUE value
- `<outfile-name>` is any valid data type, hof will infer the format from the extension

[put an example here]


### Placing data

When supplying one more more data files during input,
you can tell hof where to put a data file in the CUE value
by using the `@` notation:

```text
hof gen config.yaml@app.config \
        devops.yaml@app.devops \
        schema.cue \
    -T :stg.app=app.yaml \
    -T :stg.helm=values.yaml
```

You can have local app config combined with organization wide values
to generate the inputs for multiple tools and workflows.
You can intermix existing data with CUE throughout the software stack and life-cycle.


## Module generators

Even hof's commands will get cumbersome when you want to generate many files.
You can use generators to reduce the complexity of the commands.

hof generators can be package as modules,
dependency managed, and give you reusable bundles of CUE and hof.
When authoring a module, you can use the _data file_ features
to capture commands like we just saw into code.
We can then shorten our commands significantly.


```text
hof gen -G app -t env=stg
```

[put an example here]


When you configure a data file for output in your generator configuration,
the fields you set are different.

- `Val` is the input value, instead of `In`
- `DatafileFormat` tells hof to generate data and what format
- `Filepath` for where to write the file

```text
app: gen.Generator & {
  // input values, user provides these
  config: {...}  @embed(config.yaml)     // embed a file local to the user
  devops: {...}  & policies.#Schema      // apply an imported schema
  env:    string | *"dev" @tag(env)      // tag with a default

  // combine, reshape, and validate inputs
  // using all the capabilities of CUE
  helm: {
    // apply some policies to other fields
    devops.policies[env]

    // use local app configuration to set helm values
    app: config[env].resources
    iam: config[env].permissions.iam
    db:  app.db[env]
  }
  
  // output files
  Out: [{
    Filepath: "app.yaml"
    // pass a value
    Val: config[env].settings
  },{
    Filepath: "ci.yaml"
    // combine values inplace
    Val: {
      // apply common values
      job: devops.ci[env]
      // apply local values
      job: steps: config.build.steps
    }
  },{
    Filepath: "values-\(env).yaml"
    // separate out more complex transforms
    Val: helm
  }]
}
```

[Links to real world examples]
