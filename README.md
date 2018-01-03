# Service Catalog CLI

[![CircleCI](https://circleci.com/gh/Azure/service-catalog-cli.svg?style=svg&circle-token=98d6d64c981e70b76736fb3f05a0b41b4fec47cf)](https://circleci.com/gh/Azure/service-catalog-cli)

This project is a command-line interface (CLI) for interacting with
[Kubernetes Service Catalog](https://github.com/kubernetes-incubator/service-catalog)
resources.

| 🚨  | Our releases follow [semver](https://semver.org) and the project is in **alpha** status. This means that no assurances are made about backwards compatibility or stability. Before we hit v1.0.0, breaking changes are indicated by bumping the minor version. |
|---|---|

`svcat` is a domain-specific tool to make interacting with the Service Catalog easier.
While many of its commands have analogs to `kubectl`, our goal is to streamline and optimize
the operator experience, contributing useful code back upstream to Kubernetes where applicable.

`svcat` communicates with the Service Catalog API through the [aggregated API][agg-api] endpoint on a
Kubernetes cluster.

[agg-api]: https://kubernetes.io/docs/concepts/api-extension/apiserver-aggregation/

# Prerequisites
In order to use `svcat`, you will need:

* A Kubernetes cluster running v1.7+ or higher
* A broker compatible with the Open Service Broker API installed on the cluster, such as:
  * [Open Service Broke for Azure](https://github.com/Azure/helm-charts/tree/master/open-service-broker-azure)
* The [Service Catalog](https://github.com/kubernetes-incubator/service-catalog/blob/master/docs/install.md) installed on the cluster.

# Install
Follow the appropriate instructions for your shell to install `svcat`:

## Bash
```bash
curl -sLO https://servicecatalogcli.blob.core.windows.net/cli/latest/$(uname -s)/$(uname -m)/svcat
chmod +x ./svcat
mv ./svcat /usr/local/bin/
svcat --version
```

## PowerShell

```powershell
iwr 'https://servicecatalogcli.blob.core.windows.net/cli/latest/Windows/x86_64/svcat.exe' -UseBasicParsing -OutFile svcat.exe
mkdir -f ~\bin
$env:PATH += ";${pwd}\bin"
svcat --version
```

The snippet above adds a directory to your PATH for the current session only.
You will need to find a permanent location for it and add it to your PATH.

## Manual
1. Download the appropriate binary for your operating system:
    * macOS: https://servicecatalogcli.blob.core.windows.net/cli/latest/Darwin/x86_64/svcat
    * Windows: https://servicecatalogcli.blob.core.windows.net/cli/latest/Windows/x86_64/svcat.exe
    * Linux: https://servicecatalogcli.blob.core.windows.net/cli/latest/Linux/x86_64/svcat
1. Make the binary executable.
1. Move the binary to a directory on your `PATH`.

# Use

This CLI is called `svcat` on the command line. Run `svcat -h` to see the available
commands.

Below are some common tasks made easy with `svcat`. The example output assumes that [Open Service Broker for Azure](https://github.com/Azure/open-service-broker-azure/) is installed on the cluster.

## Find brokers installed on the cluster

```console
$ svcat get brokers
  NAME                               URL                                STATUS
+------+--------------------------------------------------------------+--------+
  osba   http://osba-open-service-broker-azure.osba.svc.cluster.local   Ready
```

## Trigger a sync of a broker's catalog

```console
$ svcat sync broker osba
Successfully fetched catalog entries from osba broker
```

## List available service classes

```console
$ svcat get classes
            NAME                      DESCRIPTION                             UUID
+--------------------------+--------------------------------+--------------------------------------+
  azure-rediscache           Azure Redis Cache                0346088a-d4b2-4478-aa32-f18e295ec1d9
  azure-storage              Azure Storage                    2e2fc314-37b6-4587-8127-8f9ee8b33fea
  azure-aci                  Azure Container Instance         451d5d19-4575-4d4a-9474-116f705ecc95
  azure-cosmos-document-db   Azure DocumentDB provided by     6330de6f-a561-43ea-a15e-b99f44d183e6
                             CosmosDB and accessible via
                             SQL (DocumentDB), Gremlin
                             (Graph), and Table (Key-Value)
                             APIs
  azure-servicebus           Azure Service Bus                6dc44338-2f13-4bc5-9247-5b1b3c5462d3
  azure-eventhub             Azure Event Hub Service          7bade660-32f1-4fd7-b9e6-d416d975170b
  azure-cosmos-mongo-db      MongoDB on Azure provided by     8797a079-5346-4e84-8018-b7d5ea5c0e3a
                             CosmosDB
  azure-mysqldb              Azure Database for MySQL         997b8372-8dac-40ac-ae65-758b4a5075a5
                             Service
  azure-postgresqldb         Azure Database for PostgreSQL    b43b4bba-5741-4d98-a10b-17dc5cee0175
                             Service
  azuresearch                Azure Search                     c54902aa-3027-4c5c-8e96-5b3d3b452f7f
  azure-keyvault             Azure Key Vault                  d90c881e-c9bb-4e07-a87b-fcfe87e03276
  azure-sqldb                Azure SQL Database Service       fb9bc99e-0aa9-11e6-8a8a-000d3a002ed5
```

## View service plans associated with a class

```console
$ svcat describe class azure-mysqldb --traverse
  Name:          azure-mysqldb
  Description:   Azure Database for MySQL Service
  UUID:          997b8372-8dac-40ac-ae65-758b4a5075a5
  Status:        Active
  Tags:          Azure, MySQL, Database

Plans:
     NAME             DESCRIPTION
+-------------+-------------------------+
  standard800   Standard Tier, 800 DTUs
  basic100      Basic Tier, 100 DTUs
  basic50       Basic Tier, 50 DTUs.
  standard200   Standard Tier, 200 DTUs
  standard400   Standard Tier, 400 DTUs
  standard100   Standard Tier, 100 DTUs
```

## Provision a service

```console
$ svcat provision quickstart-wordpress-mysql-instance \
    --class azure-mysqldb --plan standard800 \
    -p location=eastus -p resourceGroup=default
```

## View all instances of a service plan on the cluster

```console
$ svcat describe plan standard800 --traverse
  Name:          standard800
  Description:   Standard Tier, 800 DTUs
  UUID:          08e4b43a-36bc-447e-a81f-8202b13e339c
  Class:         azure-mysqldb
  Status:        Active
  Free:          false

Instances:
                NAME                  NAMESPACE   STATUS
+-----------------------------------+-----------+--------+
  deepthoughts-ghost-mysql-instance   cms         Ready
  ponycms-wordpress-mysql-instance    cms         Ready
```

## List all service instances in a namespace

```console
$ svcat get bindings --namespace cms
                NAME                 NAMESPACE               INSTANCE                STATUS
+----------------------------------+-----------+-----------------------------------+--------+
  deepthoughts-ghost-mysql-binding   cms         deepthoughts-ghost-mysql-instance   Ready
  ponycms-wordpress-mysql-binding    cms         ponycms-wordpress-mysql-instance    Ready
```

## View the details of a service instance

```console
$ svcat describe instance ponycms-wordpress-mysql-instance --traverse
  Name:        ponycms-wordpress-mysql-instance
  Namespace:   cms
  Status:      Ready - The instance was provisioned successfully @ 2017-11-30 13:11:49 -0600 CST
  Class:       azure-mysqldb
  Plan:        standard800

Class:
  Name:     azure-mysqldb
  UUID:     997b8372-8dac-40ac-ae65-758b4a5075a5
  Status:   Active

Plan:
  Name:     standard800
  UUID:     08e4b43a-36bc-447e-a81f-8202b13e339c
  Status:   Active

Broker:
  Name:     asb
  Status:   Ready
```

## Unbind all applications from an instance
```console
$ svcat unbind ponycms-wordpress-mysql-instance
```

## Unbind a single application from an instance
```console
$ svcat unbind --name ponycms-wordpress-mysql-binding
```

## Delete a service instance
Deprovisioning is the process of preparing an instance to be removed, and then deleting it.

```console
$ svcat deprovision ponycms-wordpress-mysql-instance
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
