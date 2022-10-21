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
	if len(gotStrings) != 3 {
		t.Errorf("want get 3 strings, get %v strings", len(gotStrings))
		return
	}
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

func TestStepsDefault(t *testing.T) {
	steps := New()
	wantStrings := []string{"haha1", "haha2", "haha3", "haha4", "haha5",
		"haha6", "haha7", "haha8", "haha9", "haha10"}
	gotStrings := []string{}
	steps.RegSysTenthSteps(InstanceStep{"haha10", func() error {
		gotStrings = append(gotStrings, "haha10")
		return nil
	}})
	steps.RegSysNinethSteps(InstanceStep{"haha9", func() error {
		gotStrings = append(gotStrings, "haha9")
		return nil
	}})
	steps.RegSysFirstSteps(InstanceStep{"haha1", func() error {
		gotStrings = append(gotStrings, "haha1")
		return nil
	}})
	steps.RegSysSecondSteps(InstanceStep{"haha2", func() error {
		gotStrings = append(gotStrings, "haha2")
		return nil
	}})
	steps.RegSysFourthSteps(InstanceStep{"haha4", func() error {
		gotStrings = append(gotStrings, "haha4")
		return nil
	}})
	steps.RegSysFifthSteps(InstanceStep{"haha5", func() error {
		gotStrings = append(gotStrings, "haha5")
		return nil
	}})
	steps.RegSysThirdSteps(InstanceStep{"haha3", func() error {
		gotStrings = append(gotStrings, "haha3")
		return nil
	}})
	steps.RegSysSixthSteps(InstanceStep{"haha6", func() error {
		gotStrings = append(gotStrings, "haha6")
		return nil
	}})
	steps.RegSysSeventhSteps(InstanceStep{"haha7", func() error {
		gotStrings = append(gotStrings, "haha7")
		return nil
	}})
	steps.RegSysEighthSteps(InstanceStep{"haha8", func() error {
		gotStrings = append(gotStrings, "haha8")
		return nil
	}})
	_ = steps.Run()
	for index, str := range wantStrings {
		if str != gotStrings[index] {
			t.Errorf("get %v, want %v\n", gotStrings[index], str)
			return
		}
	}
}
