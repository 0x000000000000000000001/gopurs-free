package Control_Monad_Free

import "gopurs/output/gopurs_runtime"

type FreeObj struct {
	Tag       int // 0 = Pure, 1 = Bind
	ValueOrFa gopurs_runtime.Value
	Binds     any
}

type BindLeaf struct {
	K gopurs_runtime.Value
}

type BindNode struct {
	Left  any
	Right any
}

func PureImpl(a gopurs_runtime.Value) any {
	return &FreeObj{Tag: 0, ValueOrFa: a, Binds: nil}
}

func LiftF(fa gopurs_runtime.Value) any {
	return &FreeObj{Tag: 1, ValueOrFa: fa, Binds: nil}
}

func BindImpl(freeAny gopurs_runtime.Value, k gopurs_runtime.Value) any {
	var free *FreeObj
	if freeAny.Type == gopurs_runtime.TypeAny {
		free = (*(*any)(freeAny.UnsafePtr)).(*FreeObj)
	} else {
		free = (*FreeObj)(freeAny.UnsafePtr) // Try direct unbox
	}

	var newBinds any
	if free.Binds == nil {
		newBinds = &BindLeaf{K: k}
	} else {
		newBinds = &BindNode{Left: free.Binds, Right: &BindLeaf{K: k}}
	}
	return &FreeObj{Tag: free.Tag, ValueOrFa: free.ValueOrFa, Binds: newBinds}
}

func ResumePrime(k gopurs_runtime.Value, j gopurs_runtime.Value, fAny gopurs_runtime.Value) any {
	var f *FreeObj
	if fAny.Type == gopurs_runtime.TypeAny {
		f = (*(*any)(fAny.UnsafePtr)).(*FreeObj)
	} else {
		f = (*FreeObj)(fAny.UnsafePtr)
	}

	for {
		if f.Tag == 0 { // Pure
			curr := f.Binds
			var stack []any
			var first gopurs_runtime.Value
			hasFirst := false

			for curr != nil {
				if leaf, ok := curr.(*BindLeaf); ok {
					first = leaf.K
					hasFirst = true
					break
				} else if node, ok := curr.(*BindNode); ok {
					stack = append(stack, node.Right)
					curr = node.Left
				}
			}

			if !hasFirst {
				return gopurs_runtime.Apply(j, f.ValueOrFa)
			}

			var restBinds any
			for _, s := range stack {
				if restBinds == nil {
					restBinds = s
				} else {
					restBinds = &BindNode{Left: s, Right: restBinds}
				}
			}

			f2Any := gopurs_runtime.Apply(first, f.ValueOrFa)
			var f2 *FreeObj
			if f2Any.Type == gopurs_runtime.TypeAny {
				f2 = (*(*any)(f2Any.UnsafePtr)).(*FreeObj)
			} else {
				f2 = (*FreeObj)(f2Any.UnsafePtr)
			}

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
			cont := func(b gopurs_runtime.Value) gopurs_runtime.Value {
				return gopurs_runtime.Box(&FreeObj{Tag: 0, ValueOrFa: b, Binds: f.Binds})
			}
			kfAny := gopurs_runtime.Apply(k, f.ValueOrFa)
			
			contWrapped := gopurs_runtime.Func(func(p0_0 gopurs_runtime.Value) gopurs_runtime.Value {
				return cont(p0_0)
			})
			return gopurs_runtime.Apply(kfAny, contWrapped)
		}
	}
}

var BindNodeClass any = "BindNode"
var BindLeafClass any = "BindLeaf"
var FreeObjClass any = "FreeObj"
