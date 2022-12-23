# Duct Tape Resource

Custom Concourse Resource to work as duct tape, meaning it can quickly bring things together that have no native or custom resource.

## Source Configuration

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

TBD:

```yaml
resources:
- name: foobar
  type: duct-tape
  icon: application-variable
  check_every: 20m
  source:
    check:
      run: |
        #!/bin/bash
        
        date +%Y-%m-%d-%H%M
```

## Behavior

### `check`: Check for TBD

TBD

### `in`: TBD

TBD

### `out`: TBD

TBD

## Development

### Prerequisites

* Go is *required* - TBD.
* Docker is *required* - TBD.

### Running the tests

TBD.

### Contributing

Please make all pull requests to the `main` branch and ensure tests pass locally.
