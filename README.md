# Service Catalog CLI

[![CircleCI](https://circleci.com/gh/Azure/service-catalog-cli.svg?style=svg&circle-token=98d6d64c981e70b76736fb3f05a0b41b4fec47cf)](https://circleci.com/gh/Azure/service-catalog-cli)

This project is a command line interface (CLI) for interacting with 
[Kubernetes Service Catalog](https://github.com/kubernetes-incubator/service-catalog)
resources.

The goal of the CLI is to provide an easy user interface for Service Catalog users
to inspect and modify the state of the resources in the system.

This CLI is called `svc-cat` on the command line. See below for a description 
of the commands that `svc-cat` offers.

# Commands for `ClusterServiceBroker`s

To list all the brokers in the cluster:

```console
svc-cat get brokers
```

To get the status of an individual broker:

```console
svc-cat get broker <broker name>
```

# Commands for `ClusterServiceClass`es

To get a list of all the `ClusterServiceClass`es in the cluster (i.e. the catalog):

```console
svc-cat get classes
```

# Commands for `ClusterServicePlan`s

To get a list of all the `ClusterServicePlan`s in the cluster (i.e. the catalog):

```console
svc-cat get plans
```

# Commands for `ServiceInstance`s

To get a list of all the `ServiceInstance`s in a namespace:

```console
svc-cat get instances -n <namespace>
```

# Commands for `ServiceBinding`s

To get a list of all the `ServiceBinding`s in a namespace:

```console
svc-cat get bindings -n <namespace>
```

# Contributing

For details on how to contribute to this project, please see 
[contributing.md](./docs/contributing.md).

This project welcomes contributions and suggestions.  Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit https://cla.microsoft.com.

When you submit a pull request, a CLA-bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., label, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.
