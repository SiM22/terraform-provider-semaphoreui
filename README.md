# Terraform SemaphoreUI Provider

The SemaphoreUI provider enables Terraform and OpenTofu to manage [SemaphoreUI](https://semaphoreui.com/) resources.

This repository is maintained as the `SiM22` fork of the original [`CruGlobal/terraform-provider-semaphoreui`](https://github.com/CruGlobal/terraform-provider-semaphoreui). Thank you to the original upstream creator and contributors for building and open-sourcing the provider.

## What is different in this fork

This fork keeps the upstream provider foundation and adds functionality needed for environments that manage OpenTofu natively through SemaphoreUI.

Current differences include:

- native SemaphoreUI template support for `app = "tofu"`
- inventory support for `tofu_workspace`
- ongoing maintenance of fork-specific functionality until it is available upstream

## Requirements

This provider requires a running [SemaphoreUI](https://semaphoreui.com/) instance.

The provider acceptance tests target supported SemaphoreUI versions defined in [`.github/workflows/test.yml`](.github/workflows/test.yml).

## SemaphoreUI API Client

The SemaphoreUI API client is generated from the Swagger (OpenAPI 2.0) specification in [`api-docs.yml`](api-docs.yml), which is derived from the upstream SemaphoreUI API definition in the [`semaphoreui/semaphore`](https://github.com/semaphoreui/semaphore) project.

To regenerate the client, install [go-swagger](https://goswagger.io/go-swagger/install/install-binary/) and run:

```shell
task client
```

## Publishing

This fork is intended to be published as its own provider lineage under the `SiM22` namespace.

## Support

This fork is maintained on a best-effort basis. Issues and pull requests are welcome.