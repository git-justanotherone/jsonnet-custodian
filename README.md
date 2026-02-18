# jsonnet-custodian

**jsonnet-custodian** is a package manager for [Jsonnet](https://jsonnet.org/) with support for dependencies via Git or local directories. The project allows you to manage multiple versions of packages for each dependency, making it easy to use different versions as needed.

## Key Features

- **Jsonnet package management**: Install, update, and remove dependencies for Jsonnet projects. Each dependency is treated as a module and can have its own dependencies in different versions. Custodian will handle building the dependency tree, downloading all dependencies, and resolving them correctly during execution.

- **SOPS Encryption**: We know that in some cases it is necessary to add secrets to projects. To do this securely, custodian integrates the [SOPS](https://github.com/getsops/sops) library and automatically decrypts files encrypted by the library during execution, as long as the encryption credentials are accessible and the file has a `.sops.*` extension.

- **Jsonnet-compatible subcommand**: Includes a `custodian jsonnet` subcommand with behavior nearly identical to the original, enabling seamless integration. This subcommand is necessary because the original binary cannot resolve imports correctly. There is a possibility that a compatibility mode will be added in the near future, but it is not a priority.

### Planned
    
- Support for modules in subdirectories
- Compatibility with jsonnet-bundler modules
- Compatibility mode with the jsonnet command

## Installation

Clone the repository and install dependencies:

```bash
git clone <repo-url>
make build
```

## Basic Usage

```bash
# Adding a package via Git
custodian mod get github.com/user/jsonnet-package@v1.2.3

# Adding a local package
custodian mod get ./my-local-package

# Running Jsonnet with managed dependencies
custodian jsonnet <file.jsonnet>
```

## Importing Dependencies

### Git Dependencies

To add Git dependencies, simply use a single command with a project identifier in the following format:

```bash
host_fqdn/owner/repo[/branch]@version
# e.g.
github.com/git-justanotherone/jsonnet-custodian@v0.1.0 # specifies a particular tag
# or
github.com/grafana/jsonnet-libs@5a573cd6b179 # specifies a particular commit
```

By default, each dependency has a name and an identifier. For Git dependencies, the name is the repository name, but this can be adjusted in the `custodian.json` file.

To import a dependency in a Jsonnet file, just use the dependency name as the first element of the path, for example:

```jsonnet
local common = import 'jsonnet-libs/common-lib/common/main.libsonnet';
```

## Configuration

Dependencies are defined in a configuration file (`custodian.json`), where you can specify versions and sources for each package.

## Contribution

Pull requests and suggestions are welcome! The project is in its early stages, but I will soon provide more guidelines.

## License

This project is licensed under the MIT license.
