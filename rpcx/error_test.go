package rpcx

import (
	"testing"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc/codes"
)

func TestNewError(t *testing.T) {
	Convey("TestNewError", t, func() {
		err := NewError(codes.InvalidArgument, "InvalidArgument", "invalid argument", errors.New("error"))
		So(err.code, ShouldEqual, codes.InvalidArgument)
		So(err.err.Error(), ShouldEqual, "error")
		So(err.Detail.Code, ShouldEqual, "InvalidArgument")
		So(err.Detail.Message, ShouldEqual, "invalid argument")
		So(err.Detail.Status, ShouldEqual, 0)
		So(err.Detail.Refer, ShouldEqual, "")
		So(err.Detail.RequestID, ShouldEqual, "")
	})
}
