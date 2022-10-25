# go-lib-common

## Introduction ##
This is a common golang library for standardizing the usage needs of each repo
## Prerequisites ##
1. Golang 
2. Make

## Folder Structure ##
    .
    ├── cast                      # Contains a casting variables.
    ├── client                    # Contains a request configuration to access resources on the server
    │   ├── aws                   # Contains a configuration request to aws
    │   │   └── mocks             # Contains a mocks aws with mockery generate
    │   ├── moladin_evo           # Contains a configuration request to moladin evo
    │   │   └── mocks             # Contains a mocks moladin_evo with mockery generate
    │   └── notification          # Contains a file requests for notifications-related
    │       └── slack             # Contains a configuration request to slack
    │           └── mocks         # Contains a mocks slack with mockery generate
    ├── config                    # Contains a configuration the control plane
    ├── constant                  # Contains a frequently used constant variable
    ├── context                   # Contains a functions related to contexts
    ├── data_source               # Contains a configuration database
    ├── errors                    # Contains a errors handler
    ├── logger                    # Contains a logging files
    ├── middleware                # Contains a middlewares configuration
    │   └── gin                   # Contains a middleware related to web framework gin
    │       ├── auth              # Contains a auth middleware
    │       │   └── mocks         # Contains a mocks auth with mockery generate
    │       ├── panic_recovery    # Contains a panic recovery handler
    │       │   └── mocks         # Contains a mocks panic recovery with mockery generate
    │       └── tracer            # Contains a sentry and logging
    │           └── mocks         # Contains a mocks tracer with mockery generate
    ├── response                  # Contains a response format
    ├── sentry                    # Contains a sentry configurations 
    │   └── mocks                 # Contains a mocks sentry with mockery generate
    ├── strings                   # Contains a converter string to other type
    ├── time                      # Contains a times configurations
    │   └── mocks                 # Contains a mocks time with mockery generate
    └── validator                 # Contains a validator configurations

## How to Import ##
1. You must read to configure your computer's environment [RFC](https://moladin.atlassian.net/wiki/spaces/TEC/pages/380993667/Setting+up+Your+Workstation+for+Backend+Engineers)
2. Run this command in your repo 
    ```bash
    $ go get module bitbucket.org/moladinTech/go-lib-common@v0.0.11
    ```
3. Once successful, you're ready to use the common library

## How to Run Test ##
1. Install [Mockery](https://github.com/vektra/mockery)
2. Run generates files mock use the following command:
    ```bash
    $ go generate ./...
    ```
3. Run the testing
    ```bash
    $ make test-unit
    ```

## Linter

## Contribution guidelines
- Add/update new code
- Writing tests
- Create Pull Request (Checklist all the fields in the description)
- Code review
- Merge to branch `main`
- Create tag if needed using **Semantic Versioning**

## Who do I talk to?
- [Contact Internal Tools Team](https://moladin.atlassian.net/wiki/spaces/FC/overview)
- Repo owner or admin