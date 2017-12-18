# Maintainer's Guide

## Releasing
When a pull request is merged, a build is executed on master, but nothing is published. When you are ready
to release a new version of `svcat`:

1. Create a tag: `git tag vX.Y.Z`.
1. Push the tag upstream: `git push upstream vX.Y.Z`.
1. [Add release notes](https://github.com/Azure/service-catalog-cli/releases/new) for the tag on GitHub.

This triggers a build for the tag. Binaries every supported platform are compiled, and then published
them to `https://servicecatalogcli.blob.core.windows.net/cli/VERSION/OS/ARCH/svcat`. The contents of the
release are copied to `https://servicecatalogcli.blob.core.windows.net/cli/latest` so that we have a
permalink to install the latest release.
