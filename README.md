# Duct Tape Resource

Custom Concourse Resource to work as duct tape, meaning it can quickly bring things together that have no native or custom resource.

## Source Configuration

* `id`: *Optional* Use a custom identifier, if none is specified, `ref` is used as the default.
* `check`: *Required* The configuration to define what happens for `check`.
  * `env`: *Optional* List of key/value pairs that are set as environment variables in the context of the commands.
  * `before`: *Optional* Command to be executed before the `run` command. Any output of this command will redirected to StdErr.
  * `run`: *Required* Command to be executed.

* `in`: *Required* The configuration to define what happens for `in`.
  * `env`: *Optional* List of key/value pairs that are set as environment variables in the context of the commands.
  * `before`: *Optional* Command to be executed before the `run` command. Any output of this command will redirected to StdErr.
  * `run`: *Optional* Command to be executed. If undefined, it is considered a no-op.

* `out`: *Required* The configuration to define what happens for `out`.
  * `env`: *Optional* List of key/value pairs that are set as environment variables in the context of the commands.
  * `before`: *Optional* Command to be executed before the `run` command. Any output of this command will redirected to StdErr.
  * `run`: *Optional* Command to be executed. If undefined, it is considered a no-op.

### Example

Since it is a custom resource type, it has to be configured once in the pipeline configuration.

```yaml
resource_types:
- name: duct-tape-resource
  type: docker-image
  source:
   repository: ghcr.io/homeport/duct-tape-resource
   tag: latest
```

The following sample shows the main fields that can be used. For example, to define a generic check based on inline shell script, use `env` section to feed in variables (i.e. secrets) to the shell script. The code in `before` is only executed once for the resource. With each iteration of the Concourse check, the `run` section is executed.

```yaml
resources:
- name: foobar
  type: duct-tape-resource
  icon: application-variable
  check_every: 20m
  source:
    check:
      env:
        FOO: ((bar))
        SOME: variable
      before: |
        #!/bin/bash

        echo "run something once before the checks"

      run: |
        #!/bin/bash
        
        date +%Y-%m-%d-%H%M
```

## Behavior

### `check`: Run custom check

When Concourse runs the check of the resource, the `run` section of the resource configuration will be executed. Each line of the output will treated as a value for the resource result JSON. For example, if the `check` `run` code prints `something` to StdOut, the output JSON will be `[{"ref": "something"}]`, with `ref` being the default identifier used if none is explicitly provided.

### `in`: No-op

No-op.

### `out`: No-op

No-op.

## Development

### Prerequisites

* Go is *required* - TBD.
* Docker is *required* - TBD.

### Running the tests

Run the `test` Make target to test the project:

```sh
make test
```

### Contributing

Please make all pull requests to the `main` branch and ensure tests pass locally.
