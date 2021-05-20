# patissuer

Azure DevOps Personal Access Token Issuer

## Description

This provides a CLI and package to issue Azure DevOps Personal Access Tokens for the calling user. 
The calling user will be authenticated using MSAL (OIDC) Interactive login.

Usage:

### Listing PATs

```Bash
> patissuer --aad-tenant-id <your-aad-tenant-id> --aad-client-id <your-aad-client-id> --org-url <your Azure DevOps organization URL> list
```

### Issue PATs

```Bash
> patissuer --aad-tenant-id <your-aad-tenant-id> --aad-client-id <your-aad-client-id> --org-url <your Azure DevOps organization URL> issue --token-scope <pat-token-access-scopes> --token-ttl 
```

### Environment variables

All command line flags can be replaced with env variables. The format is: `PATISSUER_<flag-with-underscore>`, so for example `--aad-tenant-id` can be replaced with the `PATISSUER_AAD_TENANT_ID` env variable

### Output formats

Using the `--output` flag, you can control the ouput format. Supported values are:

* raw - plain text format
* json - result formatted as json


## Getting Started with development

### Tools needed

* Visual Studio Code [https://code.visualstudio.com/Download](https://code.visualstudio.com/Download)
* Remote - Containers extension for VS Code: `ms-vscode-remote.remote-containers`
* [windows/amd64] Cascadia Code PL ([Powerline](https://github.com/ryanoasis/powerline-extra-symbols)) font needed for the VS Code terminal ([Download](https://github.com/microsoft/cascadia-code/releases)).  
In VS Code type `Ctrl+,` to bring up `Settings`, then search for `terminal.integrated.fontFamily` and enter `'Cascadia Code PL'` with quotes.

### Open the repo

The complete toolchain and environment needed to develop is defined in a [DevContainer](https://code.visualstudio.com/docs/remote/containers) in the `.devcontainer` folder.

1. Open the repo in VSCode
1. Hit `F1` and type/select `Remote-Containers: Reopen in Container`

**Note:** The first time you open the container (or when it has been updated) it may take a while because the container needs to be built.

### Build and Test

Building is done using [Mage](https://magefile.org/)

#### To build

```Bash
> mage build
```

#### To install tool dependencies (ginkgo, mockgen etc.)

This will be done by default when building.

```Bash
> mage toolinstall
```

#### To run the tests

```Bash
> mage test
```

#### To run the tests in watch mode (rebuild and test on save)

```Bash
> mage watch
```

#### To check other available targets

```Bash
> mage
```

### Structure and guidelines

This repo is structured to conform with the [Go Language Standard Project Layout](https://github.com/golang-standards/project-layout)

Furthermore, linting is required and enabled by default.

## Development Tools

* Visual Studio Code: [Code.VisualStudio.Com/Download](https://code.visualstudio.com/Download)
* Remote - Containers extension for VS Code: `ms-vscode-remote.remote-containers`
* Go: [golang.org](https://golang.org) *- included in dev container*
* Mage (Build system): [magefile.org](https://magefile.org) *- included in dev container*
* Ginkgo/Gomega (Test framework): [onsi.github.io/ginkgo](http://onsi.github.io/ginkgo/)
