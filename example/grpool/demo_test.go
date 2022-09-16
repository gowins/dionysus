package grpool

import "testing"

func TestGrpoolDemo(t *testing.T) {
	err := grpoolDemo()
	if err != nil {
		t.Fatal(err.Error())
	}
}
