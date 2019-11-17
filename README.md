# krane

krane is a tool to render kubernetes manifests using a simple structure of templates and configuration.

A sample structure could be something like:

```
KRANE_ROOT
├── config
│   ├── clusters
│   │   ├── cluster01.yml
│   │   └── cluster02.yml
│   ├── config.yml
│   ├── production.yml
│   └── staging.yml
├── manifests
├── templates
```

## Configuration files

All configuration is described using yaml.


### config.yml

This file has two keys, `fleets` and `globals`, and it's mandatory.

#### `fleets`

This key describes the grouping of clusters we choose to use. 

This could be the typical staging, production environments, could be clusters living in separate cloud provider accounts, or clusters running workloads for different teams.

#### `globals`

This key will hold values common to all fleets.

### fleet.yml

This file will hold the values to be used by all clusters in that fleet. Each defined fleet must have a configuration file with its name.

Here you can also use the special key `Secrets` to hold all values that will be base64 encoded before being rendered into the secret templates that reference it.

### cluster.yml

This file will hold the overrides to be used in a particular cluster. This file is optional.

Here you can also use the special key `Exclude` that holds a list of strings that will be matched against the path of templates being rendered and determine which templates should be ruled out.

## Rendering process

The tool will look in a determined path for template files. That path will be traversed and all `.yaml` or `.yml` files loaded.
The found structure will be maintained for each of the clusters, meaning that the structure of the templates folder will detemine the structure of the final manifests folder.

### Steps

1. Wipe out the manifests folder to start from a clean slate
   - If the manifests are in a git repository this allows `git diff` to easily show us any changes introduced
1. Load all templates into memory
1. Load all configuration, merging globals and overides into each cluster configuration
1. Loop clusters, rendering the templates and appending the names of both Secrets and ConfigMaps with a hash based on their contents
1. Update all manifests that reference those Secrets and ConfigMaps
1. Write all manifests

## Building

- `make` or `make build` will build the docker image `krane`

### Testing

- `make unit-test` will build a docker image named `tests/krane` that will run all unit tests
- `make test` will run the `krane` image with `CI=true` which will render manifests in the tests folder

## Running

`krane` will run on the directory where is called, provided it's in the `PATH`.

If the actual structure exists elsewhere, please provide it's location via `KRANE_ROOT` environment variable.