package deepmerge

import (
	"github.com/apparentlymart/go-tf-func-provider/tffunc"
)

func NewProvider() *tffunc.Provider {
	p := tffunc.NewProvider()
	p.AddFunction("merge_objects", mergeObjectsFunc)
	return p
}
