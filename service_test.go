package xRPC

import (
	"fmt"
	"reflect"
	"testing"
)

type XKM int

type Args struct{ Num1, Num2 int }

func (x *XKM) Sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

func (x *XKM) sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

func _assert(condition bool, msg string, v ...interface{}) {
	if !condition {
		panic(fmt.Sprintf("assertion failed: ,"+msg, v...))
	}
}

func TestNewService(t *testing.T) {
	var xk XKM
	s := newService(&xk)
	_assert(len(s.method) == 1, "wrong service method, expect 1 but god %d", len(s.method))
	mType := s.method["Sum"]
	_assert(mType != nil, "Wrong method. Sum should not be nil")
}

func TestMethodType_Call(t *testing.T) {
	var xk XKM
	s := newService(&xk)
	mType := s.method["Sum"]

	argv := mType.newArgv()
	reply := mType.newReply()
	argv.Set(reflect.ValueOf(Args{Num1: 1, Num2: 3}))
	err := s.call(mType, argv, reply)
	_assert(err == nil && *reply.Interface().(*int) == 4 && mType.NumCalls() == 1, "failed to call XKM.Sum")
}
