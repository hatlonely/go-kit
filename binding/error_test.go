package binding

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestError(t *testing.T) {
	Convey("TestError", t, func() {
		err := &Error{
			Key:  "key1.key2",
			Code: ErrInvalidFormat,
			Err:  errors.New("invalid format"),
		}
		So(err.Error(), ShouldEqual, "code [InvalidFormat] key [key1.key2] : invalid format")

		fmt.Printf("%+v", err)
		fmt.Printf("%s", err)
	})
}
