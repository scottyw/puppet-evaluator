package eval

import (
	"testing"
	"github.com/puppetlabs/go-evaluator/eval"
	"github.com/puppetlabs/go-evaluator/resource"
	"github.com/puppetlabs/go-pspec/pspec"
	"github.com/puppetlabs/go-evaluator/types"
)

func TestPSpecs(t *testing.T) {
	pspec.RunPspecTests(t, `testdata`, func() eval.DefiningLoader {
		eval.NewGoFunction(`get_resource`,
			func(d eval.Dispatch) {
				d.Param(`Variant[Type[Resource],String]`)
				d.Function(func(c eval.Context, args []eval.PValue) eval.PValue {
					ref := types.WrapString(resource.Reference(c, args[0]))
					if r, ok := resource.Resources(c).Get(ref); ok {
						return r
					}
					if node, ok := resource.FindNode(c, ref); ok {
						return node.Resources(c).Get2(ref, eval.UNDEF)
					}
					return eval.UNDEF
				})
			},
		)

		eval.NewGoFunction(`get_edges`,
			func(d eval.Dispatch) {
				d.Param(`Variant[Resource,Type[Resource],String]`)
				d.Function(func(c eval.Context, args []eval.PValue) eval.PValue {
					if from, ok := resource.FindNode(c, args[0]); ok {
						return resource.GetGraph(c).Edges(from);
					}
					return eval.EMPTY_ARRAY
				})
			},
		)

		eval.NewGoFunction(`get_resources`,
			func(d eval.Dispatch) {
				d.Function(func(c eval.Context, args []eval.PValue) eval.PValue {
					return resource.Resources(c)
				})
			},
		)

		eval.NewGoFunction(`nodes_from`,
			func(d eval.Dispatch) {
				d.Param(`Variant[Resource,Type[Resource],String]`)
				d.Function(func(c eval.Context, args []eval.PValue) eval.PValue {
					if from, ok := resource.FindNode(c, args[0]); ok {
						return resource.GetGraph(c).FromNode(from);
					}
					return eval.EMPTY_ARRAY
				})
			},
		)

		eval.NewGoFunction(`nodes_to`,
			func(d eval.Dispatch) {
				d.Param(`Variant[Resource,Type[Resource],String]`)
				d.Function(func(c eval.Context, args []eval.PValue) eval.PValue {
					if to, ok := resource.FindNode(c, args[0]); ok {
						return resource.GetGraph(c).ToNode(to);
					}
					return eval.EMPTY_ARRAY
				})
			},
		)

		return eval.StaticResourceLoader().(eval.DefiningLoader)
	})
}