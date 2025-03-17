# OpenTofu "deepmerge" provider

This is a utility provider for OpenTofu which includes a single function that takes one or
more values and performs a simplistic "deep merge" across all of them, returning the
result.

This provider is published in the OpenTofu Registry as `registry.opentofu.org/apparentlymart/deepmerge`.

## Contributing

This provider is considered feature complete, and so no feature requests will be accepted.

Bug reports are welcome, but any behavior with input that the function already accepts and successfully returns a result is frozen for backward-compatibility, even if it is not the behavior you'd prefer.

You are welcome to fork this codebase and publish your own version of it if you have a different opinion about what "deep merge" should mean.

## Copyright and License

This provider was developed under sponsorship of the OpenTofu project and is therefore copyright The OpenTofu Authors, but it is not an official product of the OpenTofu core team.

You may distribute and derive from the source code of this project under the terms of the Mozilla Public License, as described in [LICENSE](./LICENSE).
