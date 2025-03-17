package deepmerge

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/zclconf/go-cty-debug/ctydebug"
	"github.com/zclconf/go-cty/cty"
)

func TestMergeObjectsFunc(t *testing.T) {
	p := NewProvider()
	f := p.CallStub("merge_objects")

	tests := map[string]struct {
		args []cty.Value
		want cty.Value
	}{
		"one argument, not object": {
			[]cty.Value{
				cty.StringVal("hello"),
			},
			cty.StringVal("hello"),
		},
		"one argument, empty object": {
			[]cty.Value{
				cty.EmptyObjectVal,
			},
			cty.EmptyObjectVal,
		},
		"one argument, non-empty object": {
			[]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"a": cty.True,
				}),
			},
			cty.ObjectVal(map[string]cty.Value{
				"a": cty.True,
			}),
		},
		"one argument, null": {
			[]cty.Value{
				cty.NullVal(cty.String),
			},
			cty.NullVal(cty.DynamicPseudoType), // type gets erased for simplicity's sake
		},
		"one argument, completely unknown": {
			[]cty.Value{
				cty.DynamicVal,
			},
			cty.DynamicVal,
		},
		"one argument, unknown value of empty object type": {
			[]cty.Value{
				cty.UnknownVal(cty.EmptyObject),
			},
			// Becomes a known object because we can tell by its type that it has no attributes
			cty.EmptyObjectVal,
		},
		"one argument, unknown value of non-empty object type": {
			[]cty.Value{
				cty.UnknownVal(cty.Object(map[string]cty.Type{"a": cty.String})),
			},
			// Becomes a known object because we can tell by its type what attributes it has
			cty.ObjectVal(map[string]cty.Value{
				"a": cty.UnknownVal(cty.String),
			}),
		},
		"one argument, empty map": {
			[]cty.Value{
				cty.MapValEmpty(cty.String),
			},
			// Becomes an empty object instead, because this function never returns maps
			cty.EmptyObjectVal,
		},
		"one argument, non-empty map": {
			[]cty.Value{
				cty.MapVal(map[string]cty.Value{
					"a": cty.True,
				}),
			},
			// Becomes an object instead, because this function never returns maps
			cty.ObjectVal(map[string]cty.Value{
				"a": cty.True,
			}),
		},

		"two objects, disjoint attributes": {
			[]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"a": cty.StringVal("a value"),
				}),
				cty.ObjectVal(map[string]cty.Value{
					"b": cty.StringVal("b value"),
				}),
			},
			cty.ObjectVal(map[string]cty.Value{
				"a": cty.StringVal("a value"),
				"b": cty.StringVal("b value"),
			}),
		},
		"two objects, additional attributes": {
			[]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"a": cty.StringVal("a value 1"),
				}),
				cty.ObjectVal(map[string]cty.Value{
					"a": cty.StringVal("a value 2"),
					"b": cty.StringVal("b value"),
				}),
			},
			cty.ObjectVal(map[string]cty.Value{
				"a": cty.StringVal("a value 2"),
				"b": cty.StringVal("b value"),
			}),
		},
		"two objects, fewer attributes": {
			[]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"a": cty.StringVal("a value 1"),
					"b": cty.StringVal("b value"),
				}),
				cty.ObjectVal(map[string]cty.Value{
					"a": cty.StringVal("a value 2"),
				}),
			},
			cty.ObjectVal(map[string]cty.Value{
				"a": cty.StringVal("a value 2"),
				"b": cty.StringVal("b value"),
			}),
		},
		"object into map": {
			[]cty.Value{
				cty.MapVal(map[string]cty.Value{
					"a": cty.StringVal("a value 1"),
					"b": cty.StringVal("b value"),
				}),
				cty.ObjectVal(map[string]cty.Value{
					"a": cty.StringVal("a value 2"),
				}),
			},
			cty.ObjectVal(map[string]cty.Value{
				"a": cty.StringVal("a value 2"),
				"b": cty.StringVal("b value"),
			}),
		},
		"map into object": {
			[]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"a": cty.StringVal("a value 1"),
					"b": cty.StringVal("b value"),
				}),
				cty.MapVal(map[string]cty.Value{
					"a": cty.StringVal("a value 2"),
				}),
			},
			cty.ObjectVal(map[string]cty.Value{
				"a": cty.StringVal("a value 2"),
				"b": cty.StringVal("b value"),
			}),
		},
		"map into map": {
			[]cty.Value{
				cty.MapVal(map[string]cty.Value{
					"a": cty.StringVal("a value 1"),
					"b": cty.StringVal("b value"),
				}),
				cty.MapVal(map[string]cty.Value{
					"a": cty.StringVal("a value 2"),
				}),
			},
			// This function never returns maps, even if all of the inputs are maps
			cty.ObjectVal(map[string]cty.Value{
				"a": cty.StringVal("a value 2"),
				"b": cty.StringVal("b value"),
			}),
		},

		"nested objects": {
			[]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"nested": cty.ObjectVal(map[string]cty.Value{
						"a": cty.StringVal("a value 1"),
						"b": cty.StringVal("b value 1"),
					}),
				}),
				cty.ObjectVal(map[string]cty.Value{
					"nested": cty.ObjectVal(map[string]cty.Value{
						"a": cty.StringVal("a value 2"),
					}),
				}),
			},
			cty.ObjectVal(map[string]cty.Value{
				"nested": cty.ObjectVal(map[string]cty.Value{
					"a": cty.StringVal("a value 2"),
					"b": cty.StringVal("b value 1"),
				}),
			}),
		},
		"nested maps": {
			[]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"nested": cty.MapVal(map[string]cty.Value{
						"a": cty.StringVal("a value 1"),
						"b": cty.StringVal("b value 1"),
					}),
				}),
				cty.ObjectVal(map[string]cty.Value{
					"nested": cty.MapVal(map[string]cty.Value{
						"a": cty.StringVal("a value 2"),
					}),
				}),
			},
			cty.ObjectVal(map[string]cty.Value{
				// As usual, this function never produces maps
				"nested": cty.ObjectVal(map[string]cty.Value{
					"a": cty.StringVal("a value 2"),
					"b": cty.StringVal("b value 1"),
				}),
			}),
		},
		"nested objects in maps": {
			[]cty.Value{
				// These maps will both get reinterpreted into objects
				// as part of merging, so it doesn't matter that their
				// element types are incompatible with one another.
				cty.MapVal(map[string]cty.Value{
					"nested": cty.ObjectVal(map[string]cty.Value{
						"a": cty.StringVal("a value 1"),
						"b": cty.StringVal("b value 1"),
					}),
				}),
				cty.MapVal(map[string]cty.Value{
					"nested": cty.ObjectVal(map[string]cty.Value{
						"a": cty.StringVal("a value 2"),
					}),
				}),
			},
			cty.ObjectVal(map[string]cty.Value{
				"nested": cty.ObjectVal(map[string]cty.Value{
					"a": cty.StringVal("a value 2"),
					"b": cty.StringVal("b value 1"),
				}),
			}),
		},

		"two objects, first unknown": {
			[]cty.Value{
				cty.UnknownVal(cty.Object(map[string]cty.Type{
					"a": cty.String,
					"b": cty.String,
				})),
				cty.ObjectVal(map[string]cty.Value{
					"b": cty.StringVal("b value"),
				}),
			},
			cty.ObjectVal(map[string]cty.Value{
				"a": cty.UnknownVal(cty.String),
				"b": cty.StringVal("b value"),
			}),
		},
		"two objects, second unknown": {
			[]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"a": cty.StringVal("a value"),
					"b": cty.StringVal("b value"),
				}),
				cty.UnknownVal(cty.Object(map[string]cty.Type{
					"a": cty.String,
				})),
			},
			cty.ObjectVal(map[string]cty.Value{
				"a": cty.UnknownVal(cty.String),
				"b": cty.StringVal("b value"),
			}),
		},
		"unknown map into object": {
			[]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"a": cty.StringVal("a value"),
					"b": cty.StringVal("b value"),
				}),
				cty.UnknownVal(cty.Map(cty.String)),
			},
			cty.DynamicVal, // can't predict anything about the result when there's an unknown map
		},
		"object into unknown map": {
			[]cty.Value{
				cty.UnknownVal(cty.Map(cty.String)),
				cty.ObjectVal(map[string]cty.Value{
					"a": cty.StringVal("a value"),
					"b": cty.StringVal("b value"),
				}),
			},
			cty.DynamicVal, // can't predict anything about the result when there's an unknown map
		},
		"nested unknown map into nested object": {
			[]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					"nested": cty.ObjectVal(map[string]cty.Value{
						"a": cty.StringVal("a value 1"),
						"b": cty.StringVal("b value 1"),
					}),
				}),
				cty.ObjectVal(map[string]cty.Value{
					"nested": cty.UnknownVal(cty.Map(cty.String)),
				}),
			},
			// The outer object is still predictable, but its
			// nested attribute is not.
			cty.ObjectVal(map[string]cty.Value{
				"nested": cty.DynamicVal,
			}),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := f(test.args...)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(test.want, got, ctydebug.CmpOptions); diff != "" {
				t.Error("wrong result\n" + diff)
			}
		})
	}
}
