package option

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestNew(t *testing.T) {
	convey.Convey("test new option", t, func() {
		o := New("testOption", 1)
		o.Name()
		o.Value()
	})
}
