package deepmerge

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

var mergeObjectsFunc = &function.Spec{
	Description: "Recursively merges an object and other objects nested directly within it.",
	VarParam: &function.Parameter{
		Name:             "val",
		Type:             cty.DynamicPseudoType,
		AllowDynamicType: true,
		AllowNull:        true,
		AllowUnknown:     true,
	},
	Type: func(args []cty.Value) (cty.Type, error) {
		if len(args) == 0 {
			return cty.DynamicPseudoType, fmt.Errorf("must pass at least one argument")
		}

		// We'll decide the return type dynamically in the implementation function,
		// since if there are any maps we can't know what keys they have until
		// they are known values.
		return cty.DynamicPseudoType, nil
	},
	Impl: mergeAllValues,
}

func mergeAllValues(vals []cty.Value, retTy cty.Type) (cty.Value, error) {
	ret := cty.EmptyObjectVal
	for _, v := range vals {
		ret = mergeValues(ret, v)
	}
	return ret, nil
}

func mergeValues(a, b cty.Value) cty.Value {
	if a == cty.DynamicVal {
		// If we encounter DynamicVal at any point then we
		// can't possibly do any better because we have no
		// idea what we'd be merging into.
		return a
	}
	if b.IsNull() {
		// If b is null then we return null but we erase its
		// type since all nulls are considered equal for our
		// purposes here.
		return cty.NullVal(cty.DynamicPseudoType)
	}
	if a == cty.NilVal || !a.Type().IsObjectType() {
		// If the first value isn't an object type then we'll
		// start with an empty object and try to merge into
		// that instead.
		a = cty.EmptyObjectVal
	}

	// Otherwise what we do here depends on the type of the second value.
	switch {
	case b.Type().IsObjectType():
		// Merging an object into an object is the best case because
		// the type information tells us what attributes we're expecting
		// even if the values are unknown, and so we can always produce
		// a known object result even though the attribute values might
		// end up all being unknown in the worst case.
		attrs := make(map[string]cty.Value)
		for name := range a.Type().AttributeTypes() {
			attrs[name] = a.GetAttr(name)
		}
		for name := range b.Type().AttributeTypes() {
			attrs[name] = mergeValues(attrs[name], b.GetAttr(name))
		}
		return cty.ObjectVal(attrs)
	case b.Type().IsMapType():
		// As a measure of pragmatism we allow merging a map into an
		// object, but we can only do that if the map is already known.
		if !b.IsKnown() {
			// Can't even predict what type the result will be because
			// that depends on what keys are in the map.
			return cty.DynamicVal
		}
		attrs := make(map[string]cty.Value)
		for name := range a.Type().AttributeTypes() {
			attrs[name] = a.GetAttr(name)
		}
		for it := b.ElementIterator(); it.Next(); {
			k, v := it.Element()
			name := k.AsString()
			attrs[name] = mergeValues(attrs[name], v)
		}
		return cty.ObjectVal(attrs)
	default:
		// For any other type the new value just replaces whatever
		// came before it.
		return b
	}
}
