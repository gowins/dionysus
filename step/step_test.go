package step

import (
	"fmt"
	"testing"
)

func TestSteps_Run(t *testing.T) {
	steps := New()
	wantStrings := []string{"haha3", "haha2", "haha1"}
	gotStrings := []string{}
	steps.RegActionSteps("haha3", 10001, func() error {
		gotStrings = append(gotStrings, "haha3")
		return nil
	})
	steps.RegActionSteps("haha1", 10003, func() error {
		gotStrings = append(gotStrings, "haha1")
		return nil
	})
	steps.RegActionSteps("haha2", 10002, func() error {
		gotStrings = append(gotStrings, "haha2")
		return nil
	})
	_ = steps.Run()
	for index, str := range wantStrings {
		if str != gotStrings[index] {
			t.Errorf("get %v, want %v\n", gotStrings[index], str)
			return
		}
	}
}

func TestSteps_ActionStepsAppend(t *testing.T) {
	steps := New()
	wantStrings := []string{"haha3", "haha1", "haha2"}
	gotStrings := []string{}
	_ = steps.ActionStepsAppend("haha3", func() error {
		gotStrings = append(gotStrings, "haha3")
		return nil
	})
	err := steps.ActionStepsAppend("haha4", nil)
	if err == nil {
		t.Errorf("want get error, but got nil")
	}
	_ = steps.ActionStepsAppend("haha1", func() error {
		gotStrings = append(gotStrings, "haha1")
		return nil
	})
	_ = steps.ActionStepsAppend("haha2", func() error {
		gotStrings = append(gotStrings, "haha2")
		return nil
	})
	_ = steps.Run()
	for index, str := range wantStrings {
		if str != gotStrings[index] {
			t.Errorf("get %v, want %v\n", gotStrings[index], str)
			return
		}
	}
}

func TestSteps_RunE(t *testing.T) {
	steps := New()
	wantStrings := []string{"haha5", "haha1"}
	gotStrings := []string{}
	steps.RegActionStepsE("haha3", -1, func() error {
		gotStrings = append(gotStrings, "haha3")
		return nil
	})
	steps.RegActionStepsE("haha1", 10003, func() error {
		gotStrings = append(gotStrings, "haha1")
		return fmt.Errorf("121")
	})
	steps.RegActionStepsE("haha2", 1232141, func() error {
		gotStrings = append(gotStrings, "haha2")
		return nil
	})
	steps.RegActionStepsE("haha4", 231, nil)
	steps.RegActionStepsE("haha5", 10002, func() error {
		gotStrings = append(gotStrings, "haha5")
		return nil
	})
	_ = steps.Run()
	for index, str := range wantStrings {
		if str != gotStrings[index] {
			t.Errorf("get %v, want %v\n", gotStrings[index], str)
			return
		}
	}
}
