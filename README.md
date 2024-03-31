# scaffold-cli

A simple CLI tool to generate project code from your template layout. Inspired by [nunu](https://github.com/go-nunu/nunu).

## Installation

```sh
go install github.com/crappycook/scaffold-cli@latest
```

## Use

```sh
Build new project from your layout

Usage:
  scaffold-cli [command]

Examples:
scaffold-cli new project

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  new         Create a new project.

Flags:
  -h, --help   help for scaffold-cli

Use "scaffold-cli [command] --help" for more information about a command.
```

## Example

```sh
scaffold-cli new $MODULE_NAME -r <repo_url>
```
