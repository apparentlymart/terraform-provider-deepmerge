# `merge_objects` function

Takes one or more values and performs a simplistic "deep merge" across all of them.

This function starts with an empty object and then visits each argument in turn,
applying the following rules, and then returns the resulting value:

- If the new value is a null value of any type, it completely replaces the previous
  value.
- If the new value is of any non-object and non-map type, it completely replaces
  the previous value.
- If the new value is of an object type:
  1. If the current value isn't already an object, the current value is first replaced by an empty object.
  2. Each attribute of the new object in turn is merged into the current value, recursively applying these rules if the old object has an attribute of the same name.
- If the new value is of a map type:
  1. Convert the new value into an object type, retaining its keys as attribute names and its element values as attribute values.
  2. Apply the merging rules for object types as described in the previous step. The result is always of an object type, and never of a map type.

For example:

```hcl
provider::deepmerge::merge_objects(
  { a = "a1", b = "b1" },
  { a = "a2" },
)
```

The above call would produce the object value `{ a = "a2", b = "b1" }`.
