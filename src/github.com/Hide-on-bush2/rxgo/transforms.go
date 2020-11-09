// Copyright 2018 The SS.SYSU Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rxgo

import (
	"context"
	"reflect"
	"sync"
	"time"
)

var (
	typeAny        = reflect.TypeOf((*interface{})(nil)).Elem()
	typeContext    = reflect.TypeOf((*context.Context)(nil)).Elem()
	typeError      = reflect.TypeOf((*error)(nil)).Elem()
	typeBool       = reflect.TypeOf(true)
	typeObservable = reflect.TypeOf(&Observable{})
)

// transform node implementation of streamOperator
type transOperater struct {
	opFunc func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool)
}

func (tsop transOperater) op(ctx context.Context, o *Observable) {
	// must hold defintion of flow resourcs here, such as chan etc., that is allocated when connected
	// this resurces may be changed when operation routine is running.
	in := o.pred.outflow
	out := o.outflow
	//fmt.Println(o.name, "operator in/out chan ", in, out)
	var wg sync.WaitGroup

	go func() {
		end := false
		for x := range in {
			if end {
				continue
			}
			// can not pass a interface as parameter (pointer) to gorountion for it may change its value outside!
			xv := reflect.ValueOf(x)
			// send an error to stream if the flip not accept error
			if e, ok := x.(error); ok && !o.flip_accept_error {
				o.sendToFlow(ctx, e, out)
				continue
			}
			// scheduler
			switch threading := o.threading; threading {
			case ThreadingDefault:
				if tsop.opFunc(ctx, o, xv, out) {
					end = true
				}
			case ThreadingIO:
				fallthrough
			case ThreadingComputing:
				wg.Add(1)
				go func() {
					defer wg.Done()
					if tsop.opFunc(ctx, o, xv, out) {
						end = true
					}
				}()
			default:
			}
		}

		wg.Wait() //waiting all go-routines completed
		o.closeFlow(out)
	}()
}

func (tsop myOperater) op(ctx context.Context, o *Observable) {
	// must hold defintion of flow resourcs here, such as chan etc., that is allocated when connected
	// this resurces may be changed when operation routine is running.
	in := o.pred.outflow
	out := o.outflow
	//fmt.Println(o.name, "operator in/out chan ", in, out)
	var wg sync.WaitGroup

	go func() {
		end := false
		for x := range in {
			if end {
				continue
			}
			// can not pass a interface as parameter (pointer) to gorountion for it may change its value outside!
			xv := reflect.ValueOf(x)
			// send an error to stream if the flip not accept error
			if e, ok := x.(error); ok && !o.flip_accept_error {
				o.sendToFlow(ctx, e, out)
				continue
			}
			// scheduler
			switch threading := o.threading; threading {
			case ThreadingDefault:
				if tsop.opFunc(ctx, o, xv, out) {
					end = true
				}
			case ThreadingIO:
				fallthrough
			case ThreadingComputing:
				wg.Add(1)
				go func() {
					defer wg.Done()
					if tsop.opFunc(ctx, o, xv, out) {
						end = true
					}
				}()
			default:
			}
		}
		if o.flip != nil {
			buffer, ok := reflect.ValueOf(o.flip).Interface().([]reflect.Value)
			if !ok {
				panic("filter.go(op): buffer is not a slice")
			}
			for _, v := range buffer {
				o.sendToFlow(ctx, v.Interface(), out)
			}
		}

		wg.Wait() //waiting all go-routines completed
		o.closeFlow(out)
	}()
}

type myOperater struct {
	opFunc func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool)
}

func (parent *Observable) TransformOp(tf transformFunc) (o *Observable) {
	o = parent.newTransformObservable("customTransform")
	o.flip_accept_error = true

	o.flip = tf
	o.flip_accept_error = true
	o.operator = transformOperater
	return o
}

var transformOperater = transOperater{func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
	tf := o.flip.(transformFunc)
	send := func(x interface{}) (endSignal bool) {
		endSignal = o.sendToFlow(ctx, x, out)
		return
	}
	tf(ctx, x.Interface(), send)
	return
}}

// Map maps each item in Observable by the function with `func(x anytype) anytype` and
// returns a new Observable with applied items.
func (parent *Observable) Map(f interface{}) (o *Observable) {
	// check validation of f
	fv := reflect.ValueOf(f)
	inType := []reflect.Type{typeAny}
	outType := []reflect.Type{typeAny}
	b, ctx_sup := checkFuncUpcast(fv, inType, outType, true)
	if !b {
		panic(ErrFuncFlip)
	}

	o = parent.newTransformObservable("map")
	o.flip_accept_error = checkFuncAcceptError(fv)

	o.flip_sup_ctx = ctx_sup
	o.flip = fv.Interface()
	o.operator = mapOperater
	return o
}

var mapOperater = transOperater{func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {

	fv := reflect.ValueOf(o.flip)
	var params = []reflect.Value{x}
	rs, skip, stop, e := userFuncCall(fv, params)

	var item interface{} = rs[0].Interface()
	if stop {
		end = true
		return
	}
	if skip {
		return
	}
	if e != nil {
		item = e
	}
	// send data
	if !end {
		end = o.sendToFlow(ctx, item, out)
	}

	return
}}

// FlatMap maps each item in Observable by the function with `func(x anytype) (o *Observable) ` and
// returns a new Observable with merged observables appling on each items.
func (parent *Observable) FlatMap(f interface{}) (o *Observable) {
	// check validation of f
	fv := reflect.ValueOf(f)
	inType := []reflect.Type{typeAny}
	outType := []reflect.Type{typeObservable}
	b, ctx_sup := checkFuncUpcast(fv, inType, outType, true)
	if !b {
		panic(ErrFuncFlip)
	}

	o = parent.newTransformObservable("flatMap")
	o.flip_accept_error = checkFuncAcceptError(fv)

	o.flip_sup_ctx = ctx_sup
	o.flip = fv.Interface()
	o.operator = flatMapOperater
	return o
}

var flatMapOperater = transOperater{func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {

	fv := reflect.ValueOf(o.flip)
	var params = []reflect.Value{x}
	//fmt.Println("x is ", x)
	rs, skip, stop, e := userFuncCall(fv, params)

	var item = rs[0].Interface().(*Observable)

	if stop {
		end = true
		return
	}
	if skip {
		return
	}
	if e != nil {
		end = o.sendToFlow(ctx, e, out)
		if end {
			return
		}
		return
	}
	// send data
	if !end {
		if item != nil {
			// subscribe ro without any ObserveOn model
			ro := item
			for ; ro.next != nil; ro = ro.next {
			}
			ro.connect(ctx)

			ch := ro.outflow
			for x := range ch {
				end = o.sendToFlow(ctx, x, out)
				if end {
					return
				}
			}
		}
	}
	return
}}

// Filter `func(x anytype) bool` filters items in the original Observable and returns
// a new Observable with the filtered items.
func (parent *Observable) Filter(f interface{}) (o *Observable) {
	// check validation of f
	fv := reflect.ValueOf(f)
	inType := []reflect.Type{typeAny}
	outType := []reflect.Type{typeBool}
	b, ctx_sup := checkFuncUpcast(fv, inType, outType, true)
	if !b {
		panic(ErrFuncFlip)
	}

	o = parent.newTransformObservable("filter")
	o.flip_accept_error = checkFuncAcceptError(fv)

	o.flip_sup_ctx = ctx_sup
	o.flip = fv.Interface()
	o.operator = filterOperater
	return o
}

var filterOperater = transOperater{func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {

	fv := reflect.ValueOf(o.flip)
	var params = []reflect.Value{x}
	rs, skip, stop, e := userFuncCall(fv, params)

	var item interface{} = rs[0].Interface()
	if stop {
		end = true
		return
	}
	if skip {
		return
	}
	if e != nil {
		item = e
	}
	// send data
	if !end {
		if b, ok := item.(bool); ok && b {
			end = o.sendToFlow(ctx, x.Interface(), out)
		}
	}

	return
}}

func (parent *Observable) newTransformObservable(name string) (o *Observable) {
	//new Observable
	o = newObservable()
	o.Name = name

	//chain Observables
	parent.next = o
	o.pred = parent
	o.root = parent.root

	//set options
	o.buf_len = BufferLen
	return o
}

// Howework, Distinct
// var distinct_x map[interface{}]bool

// func distinct(x string) bool {
// 	if _, ok := distinct_x[x]; !ok {
// 		distinct_x[x] = true
// 		return true
// 	} else {
// 		return false
// 	}
// }

func (parent *Observable) Distinct() (o *Observable) {
	// check validation of f
	distinct_x := make(map[interface{}]bool)
	fv := reflect.ValueOf(
		func(x interface{}) bool {
			if _, ok := distinct_x[x]; !ok {
				distinct_x[x] = true
				return true
			} else {
				return false
			}
		})
	inType := []reflect.Type{typeAny}
	outType := []reflect.Type{typeBool}
	b, ctx_sup := checkFuncUpcast(fv, inType, outType, true)
	if !b {
		panic(ErrFuncFlip)
	}

	o = parent.newTransformObservable("distinct")
	o.flip_accept_error = checkFuncAcceptError(fv)

	o.flip_sup_ctx = ctx_sup
	o.flip = fv.Interface()
	o.operator = distinctOperater

	// distinct_x = make(map[interface{}]bool)
	return o
}

var distinctOperater = transOperater{func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {

	fv := reflect.ValueOf(o.flip)
	var params = []reflect.Value{x}
	rs, skip, stop, e := userFuncCall(fv, params)

	var item interface{} = rs[0].Interface()
	if stop {
		end = true
		return
	}
	if skip {
		return
	}
	if e != nil {
		item = e
	}

	// send data
	if !end {
		if b, ok := item.(bool); ok && b {
			end = o.sendToFlow(ctx, x.Interface(), out)
		}
		// if item != nil {
		// 	end = o.sendToFlow(ctx, x.Interface(), out)
		// }
		// end = o.sendToFlow(ctx, item, out)
	}

	return
}}

//Homework, ElementAt

func (parent *Observable) ElementAt(i interface{}) (o *Observable) {
	// check validation of f
	fv := reflect.ValueOf(
		func(x interface{}) bool {
			return x == i
		})
	inType := []reflect.Type{typeAny}
	outType := []reflect.Type{typeBool}
	b, ctx_sup := checkFuncUpcast(fv, inType, outType, true)
	if !b {
		panic(ErrFuncFlip)
	}

	o = parent.newTransformObservable("distinct")
	o.flip_accept_error = checkFuncAcceptError(fv)

	o.flip_sup_ctx = ctx_sup
	o.flip = fv.Interface()
	o.operator = elementAtOperater
	return o
}

var elementAtOperater = transOperater{func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {

	fv := reflect.ValueOf(o.flip)
	var params = []reflect.Value{x}
	rs, skip, stop, e := userFuncCall(fv, params)

	var item interface{} = rs[0].Interface()
	if stop {
		end = true
		return
	}
	if skip {
		return
	}
	if e != nil {
		item = e
	}

	// send data
	if !end {
		if b, ok := item.(bool); ok && b {
			end = o.sendToFlow(ctx, x.Interface(), out)
		}
		// if item != nil && item == x.Interface() {
		// 	end = o.sendToFlow(ctx, x.Interface(), out)
		// }
		// end = o.sendToFlow(ctx, item, out)
	}

	return
}}

//Homework, First

func (parent *Observable) First(f interface{}) (o *Observable) {
	// check validation of f
	flag := false
	fv := reflect.ValueOf(
		func(x interface{}) bool {
			if !flag && f.(func(interface{}) bool)(x) {
				flag = true
				return true
			} else {
				return false
			}
		})
	inType := []reflect.Type{typeAny}
	outType := []reflect.Type{typeBool}
	b, ctx_sup := checkFuncUpcast(fv, inType, outType, true)
	if !b {
		panic(ErrFuncFlip)
	}

	o = parent.newTransformObservable("distinct")
	o.flip_accept_error = checkFuncAcceptError(fv)

	o.flip_sup_ctx = ctx_sup
	o.flip = fv.Interface()
	o.operator = firstOperater
	return o
}

var firstOperater = transOperater{func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {

	fv := reflect.ValueOf(o.flip)
	var params = []reflect.Value{x}
	rs, skip, stop, e := userFuncCall(fv, params)

	var item interface{} = rs[0].Interface()
	if stop {
		end = true
		return
	}
	if skip {
		return
	}
	if e != nil {
		item = e
	}

	// send data
	if !end {
		if b, ok := item.(bool); ok && b {
			end = o.sendToFlow(ctx, x.Interface(), out)
		}
		// if item != nil && item == x.Interface() {
		// 	end = o.sendToFlow(ctx, x.Interface(), out)
		// }
		// end = o.sendToFlow(ctx, item, out)
	}

	return
}}

//Homework, Skip

func (parent *Observable) Skip(count int) (o *Observable) {
	// check validation of f
	i := 0
	fv := reflect.ValueOf(
		func(x interface{}) bool {
			i += 1
			return i > count
		})
	inType := []reflect.Type{typeAny}
	outType := []reflect.Type{typeBool}
	b, ctx_sup := checkFuncUpcast(fv, inType, outType, true)
	if !b {
		panic(ErrFuncFlip)
	}

	o = parent.newTransformObservable("distinct")
	o.flip_accept_error = checkFuncAcceptError(fv)

	o.flip_sup_ctx = ctx_sup
	o.flip = fv.Interface()
	o.operator = skipOperater
	return o
}

var skipOperater = transOperater{func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {

	fv := reflect.ValueOf(o.flip)
	var params = []reflect.Value{x}
	rs, skip, stop, e := userFuncCall(fv, params)

	var item interface{} = rs[0].Interface()
	if stop {
		end = true
		return
	}
	if skip {
		return
	}
	if e != nil {
		item = e
	}

	// send data
	if !end {
		if b, ok := item.(bool); ok && b {
			end = o.sendToFlow(ctx, x.Interface(), out)
		}
		// if item != nil && item == x.Interface() {
		// 	end = o.sendToFlow(ctx, x.Interface(), out)
		// }
		// end = o.sendToFlow(ctx, item, out)
	}

	return
}}

//Homework, Skip

func (parent *Observable) Take(count int) (o *Observable) {
	// check validation of f
	i := 0
	fv := reflect.ValueOf(
		func(x interface{}) bool {
			i += 1
			return i <= count
		})
	inType := []reflect.Type{typeAny}
	outType := []reflect.Type{typeBool}
	b, ctx_sup := checkFuncUpcast(fv, inType, outType, true)
	if !b {
		panic(ErrFuncFlip)
	}

	o = parent.newTransformObservable("distinct")
	o.flip_accept_error = checkFuncAcceptError(fv)

	o.flip_sup_ctx = ctx_sup
	o.flip = fv.Interface()
	o.operator = takeOperater
	return o
}

var takeOperater = transOperater{func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {

	fv := reflect.ValueOf(o.flip)
	var params = []reflect.Value{x}
	rs, skip, stop, e := userFuncCall(fv, params)

	var item interface{} = rs[0].Interface()
	if stop {
		end = true
		return
	}
	if skip {
		return
	}
	if e != nil {
		item = e
	}

	// send data
	if !end {
		if b, ok := item.(bool); ok && b {
			end = o.sendToFlow(ctx, x.Interface(), out)
		}
		// if item != nil && item == x.Interface() {
		// 	end = o.sendToFlow(ctx, x.Interface(), out)
		// }
		// end = o.sendToFlow(ctx, item, out)
	}

	return
}}

//Homework, Debounce

func (parent *Observable) Debounce(wait_time time.Duration) (o *Observable) {
	// check validation of f
	fv := reflect.ValueOf(
		func(x interface{}) bool {
			return true
		})
	inType := []reflect.Type{typeAny}
	outType := []reflect.Type{typeBool}
	b, ctx_sup := checkFuncUpcast(fv, inType, outType, true)
	if !b {
		panic(ErrFuncFlip)
	}

	o = parent.newTransformObservable("Debounce")
	o.flip_accept_error = checkFuncAcceptError(fv)
	o.threading = ThreadingComputing
	o.flip_sup_ctx = ctx_sup
	o.flip = fv.Interface()
	i := 0
	o.operator = transOperater{
		func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool) {
			i++
			var tempCount = i
			time.Sleep(wait_time)
			// time.Sleep(time.Second)
			if tempCount == i {
				end = o.sendToFlow(ctx, item.Interface(), out)
			}
			return
		}}
	return o
}

//Homework IgnoreElements

func (parent *Observable) IgnoreElements() (o *Observable) {
	// check validation of f
	fv := reflect.ValueOf(
		func(x interface{}) bool {
			return false
		})
	inType := []reflect.Type{typeAny}
	outType := []reflect.Type{typeBool}
	b, ctx_sup := checkFuncUpcast(fv, inType, outType, true)
	if !b {
		panic(ErrFuncFlip)
	}

	o = parent.newTransformObservable("IgnoreElements")
	o.flip_accept_error = checkFuncAcceptError(fv)

	o.flip_sup_ctx = ctx_sup
	o.flip = fv.Interface()
	o.operator = ignoreElementsOperater
	return o
}

var ignoreElementsOperater = transOperater{func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
	return false
}}

func (parent *Observable) Last() (o *Observable) {
	o = parent.newTransformObservable("last")
	o.flip_accept_error = false
	o.flip_sup_ctx = false
	o.flip = nil
	temp := make([]reflect.Value, 1)
	o.operator = myOperater{
		func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool) {
			temp[0] = item
			o.flip = temp
			return false
		}}
	return o
}

func (parent *Observable) Skiplast(n int) (o *Observable) {
	o = parent.newTransformObservable("Skiplast")
	o.flip_accept_error = false
	o.flip_sup_ctx = false
	o.flip = nil
	var temp []reflect.Value
	var tempcount = 0
	o.operator = myOperater{
		func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool) {
			if tempcount == n {
				end = o.sendToFlow(ctx, temp[0].Interface(), out)
				temp = temp[1:]
			} else {
				tempcount++
			}
			temp = append(temp, item)
			return
		}}
	return o
}

func (parent *Observable) Takelast(n int) (o *Observable) {
	o = parent.newTransformObservable("Takelast")
	o.flip_accept_error = false
	o.flip_sup_ctx = false
	o.flip = nil
	var temp []reflect.Value
	var tempcount = 0
	o.operator = myOperater{
		func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool) {
			if tempcount == n {
				temp = temp[1:]
			} else {
				tempcount++
			}
			temp = append(temp, item)
			o.flip = temp
			return
		}}
	return o
}
