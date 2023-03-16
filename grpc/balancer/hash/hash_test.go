package hash

import (
	"context"
	"fmt"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/resolver"
	"testing"
)

var connNum int

type testSubconn struct {
	num int
}

func (tc testSubconn) UpdateAddresses([]resolver.Address) {

}
func (tc testSubconn) Connect() {
	connNum = tc.num
}

func (tc testSubconn) GetOrBuildProducer(builder balancer.ProducerBuilder) (p balancer.Producer, close func()) {
	return nil, func() {
	}
}

func Test_hashPicker_Pick(t *testing.T) {
	testSubconn1 := testSubconn{num: 1}
	testSubconn2 := testSubconn{num: 2}
	testSubconn3 := testSubconn{num: 3}
	testSubconn4 := testSubconn{num: 4}
	testSubconn5 := testSubconn{num: 5}
	pickerBuilder := &hashPickerBuilder{}
	pickerBuildInfo := base.PickerBuildInfo{
		ReadySCs: map[balancer.SubConn]base.SubConnInfo{
			testSubconn1: {Address: resolver.Address{Addr: "1"}},
			testSubconn2: {Address: resolver.Address{Addr: "2"}},
			testSubconn3: {Address: resolver.Address{Addr: "3"}},
			testSubconn4: {Address: resolver.Address{Addr: "4"}},
			testSubconn5: {Address: resolver.Address{Addr: "5"}},
		},
	}
	picker := pickerBuilder.Build(pickerBuildInfo)
	pickerInfo := balancer.PickInfo{
		FullMethodName: "hello",
		Ctx:            context.Background(),
	}

	existConn := map[int]bool{}
	for i := 1; i <= 5; i++ {
		res, err := picker.Pick(pickerInfo)
		if err != nil {
			t.Errorf("want error nil, get error %v", err)
			return
		}
		res.SubConn.Connect()
		fmt.Printf("connNum is %v\n", connNum)
		if existConn[connNum] {
			t.Errorf("want not repeated")
			return
		} else {
			existConn[connNum] = true
		}
	}

	mdCtx := metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{"hash_balancer": "222"}))
	pickerInfo = balancer.PickInfo{
		FullMethodName: "hello",
		Ctx:            mdCtx,
	}
	res, err := picker.Pick(pickerInfo)
	if err != nil {
		t.Errorf("want error nil, get error %v", err)
		return
	}
	res.SubConn.Connect()
	oldNum := connNum

	for i := 1; i <= 10; i++ {
		res, err = picker.Pick(pickerInfo)
		if err != nil {
			t.Errorf("want error nil, get error %v", err)
			return
		}
		res.SubConn.Connect()
		if connNum != oldNum {
			t.Errorf("want connNum %v, get %v", oldNum, connNum)
		}
	}
}
