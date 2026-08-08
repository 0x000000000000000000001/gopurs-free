package Control_Monad_Free

import "gopurs/output/gopurs_runtime"

type FreeObj struct {
	Tag       int // 0 = Pure, 1 = Bind
	ValueOrFa any
	Binds     any
}

type BindLeaf struct {
	K any
}

type BindNode struct {
	Left  any
	Right any
}

func PureImpl(a any) any {
	return &FreeObj{Tag: 0, ValueOrFa: a, Binds: nil}
}

func LiftF(fa any) any {
	return &FreeObj{Tag: 1, ValueOrFa: fa, Binds: nil}
}

func BindImpl(freeAny any, k any) any {
	free := freeAny.(*FreeObj)
	var newBinds any
	if free.Binds == nil {
		newBinds = &BindLeaf{K: k}
	} else {
		newBinds = &BindNode{Left: free.Binds, Right: &BindLeaf{K: k}}
	}
	return &FreeObj{Tag: free.Tag, ValueOrFa: free.ValueOrFa, Binds: newBinds}
}

func ResumePrime(k any, j any, fAny any) any {
	f := fAny.(*FreeObj)
	for {
		if f.Tag == 0 { // Pure
			curr := f.Binds
			var stack []any
			var first any

			for curr != nil {
				if leaf, ok := curr.(*BindLeaf); ok {
					first = leaf.K
					break
				} else if node, ok := curr.(*BindNode); ok {
					stack = append(stack, node.Right)
					curr = node.Left
				}
			}

			if first == nil {
				return gopurs_runtime.Apply(j, gopurs_runtime.Box(f.ValueOrFa)).Unbox()
			}

			var restBinds any
			for _, s := range stack {
				if restBinds == nil {
					restBinds = s
				} else {
					restBinds = &BindNode{Left: s, Right: restBinds}
				}
			}

			f2Any := gopurs_runtime.Apply(first, gopurs_runtime.Box(f.ValueOrFa)).Unbox()
			f2 := f2Any.(*FreeObj)

			var newBinds any
			if f2.Binds == nil {
				newBinds = restBinds
			} else if restBinds == nil {
				newBinds = f2.Binds
			} else {
				newBinds = &BindNode{Left: f2.Binds, Right: restBinds}
			}

			f = &FreeObj{Tag: f2.Tag, ValueOrFa: f2.ValueOrFa, Binds: newBinds}
		} else { // Lift
			cont := func(b any) any {
				return &FreeObj{Tag: 0, ValueOrFa: b, Binds: f.Binds}
			}
			kfAny := gopurs_runtime.Apply(k, gopurs_runtime.Box(f.ValueOrFa)).Unbox()
			// kf is (b -> Free f a) -> r
			// We pass the continuation 'cont'
			// In gopurs, PureScript functions take one boxed value and return a boxed value
			contWrapped := func(p0_0 any) any {
				return cont(p0_0)
			}
			return gopurs_runtime.Apply(kfAny, gopurs_runtime.Box(contWrapped)).Unbox()
		}
	}
}

var BindNodeClass any = "BindNode"
var BindLeafClass any = "BindLeaf"
var FreeObjClass any = "FreeObj"
