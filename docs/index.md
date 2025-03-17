# "deepmerge" provider

This is a utility provider for OpenTofu which includes a single function that takes one or
more values and performs a simplistic "deep merge" across all of them, returning the
result.

To use this provider in an OpenTofu module you must first declare a dependency on the
provider to bring its functions into scope:

```hcl
terraform {
  required_providers {
    deepmerge = {
      source = "registry.opentofu.org/apparentlymart/deepmerge"
    }
  }
}
```

The local name you choose for the provider ("deepmerge" in the above example) decides
the namespace used to invoke its function. With the configuration above, the
merge function is called as `provider::deepmerge::merge_objects`.
