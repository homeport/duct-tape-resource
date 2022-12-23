# Duct Tape Resource

Custom Concourse Resource to work as duct tape, meaning it can quickly bring things together that have no native or custom resource.

## Source Configuration

* `check`: *Required.* The configuration to define what happens for `check`.
  * `script`: *Required*. Command to be executed.

* `in`: *Required.* The configuration to define what happens for `in`.

* `out`: *Required.* The configuration to define what happens for `out`.

### Example

TBD:

``` yaml
resources:
- name: foobar
  type: duct-tape
  icon: application-variable
  check_every: 20m
  source:
    check:
      script: |
        #!/bin/bash
        
        date +%Y-%m-%d-%H%M
```

## Behavior

### `check`: Check for TBD

TBD

### `in`: TBD

TBD

#### Parameters

* `something`: *Optional.* Something.

* `else`: *Optional.* Else.

### `out`: TBD

TBD

#### Parameters

* `yet`: *Required.* Yet.

* `another`: *Optional.* Another.

## Development

### Prerequisites

* Go is *required* - TBD.
* Docker is *required* - TBD.

### Running the tests

TBD.

### Contributing

Please make all pull requests to the `main` branch and ensure tests pass locally.
