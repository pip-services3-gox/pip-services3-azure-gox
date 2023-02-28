# <img src="https://uploads-ssl.webflow.com/5ea5d3315186cf5ec60c3ee4/5edf1c94ce4c859f2b188094_logo.svg" alt="Pip.Services Logo" width="200"> <br/> Azure specific components for Golang Changelog

## <a name="1.1.0"></a> 1.1.0 (2023-03-01)

### Breaking changes
* Renamed descriptors for services:
    - "\*:service:azurefunc\*:1.0" -> "\*:service:azurefunc\*:1.0"
    - "\*:service:commandable-azurefunc\*:1.0" -> "\*:service:commandable-azurefunc\*:1.0"

## <a name="1.0.0"></a> 1.0.0 (2022-09-09) 

Initial public release

### Features

- **clients** - client components for working with Azure cloud Functions.
- **connect** - components for installation and connection settings.
- **containers** - contains classes that act as containers to instantiate and run components.
- **services** - contains interfaces and classes used to create services that do operations via the Azure Function protocol.

