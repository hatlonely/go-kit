package rpcx

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc/codes"
)

func TestNewError(t *testing.T) {
	Convey("test NewError", t, func() {
		err := NewError(codes.InvalidArgument, "InvalidArgument", "invalid argument", errors.New("error"))
		So(err.code, ShouldEqual, codes.InvalidArgument)
		So(err.err.Error(), ShouldEqual, "error")
		So(err.Detail.Code, ShouldEqual, "InvalidArgument")
		So(err.Detail.Message, ShouldEqual, "invalid argument")
		So(err.Detail.Status, ShouldEqual, 0)
		So(err.Detail.Refer, ShouldEqual, "")
		So(err.Detail.RequestID, ShouldEqual, "")
	})

	Convey("test NewErrorf", t, func() {
		err := NewErrorf(codes.InvalidArgument, "InvalidArgument", "invalid argument field %s", "text")
		So(err.code, ShouldEqual, codes.InvalidArgument)
		So(err.err.Error(), ShouldEqual, "invalid argument field text")
		So(err.Detail.Code, ShouldEqual, "InvalidArgument")
		So(err.Detail.Message, ShouldEqual, "invalid argument field text")
		So(err.Detail.Status, ShouldEqual, 0)
		So(err.Detail.Refer, ShouldEqual, "")
		So(err.Detail.RequestID, ShouldEqual, "")
	})

	Convey("test NewInternalError", t, func() {
		err := NewInternalError(errors.New("mysql error"))
		So(err.code, ShouldEqual, codes.Internal)
		So(err.err.Error(), ShouldEqual, "mysql error")
		So(err.Detail.Code, ShouldEqual, "InternalError")
		So(err.Detail.Message, ShouldEqual, "mysql error")
		So(err.Detail.Status, ShouldEqual, http.StatusInternalServerError)
		So(err.Detail.Refer, ShouldEqual, "")
		So(err.Detail.RequestID, ShouldEqual, "")
		So(err.Error(), ShouldEqual, "mysql error")
	})
}

func TestErrorSet(t *testing.T) {
	Convey("TestErrorSet", t, func() {
		err := NewError(codes.InvalidArgument, "InvalidArgument", "invalid argument", errors.New("error"))

		Convey("SetRequest", func() {
			err = err.SetRequestID("test-request-id")
			So(err.Detail.RequestID, ShouldEqual, "test-request-id")
		})

		Convey("SetRefer", func() {
			err = err.SetRefer("test-refer")
			So(err.Detail.Refer, ShouldEqual, "test-refer")
		})

		Convey("SetStatus", func() {
			err = err.SetStatus(http.StatusForbidden)
			So(err.Detail.Status, ShouldEqual, http.StatusForbidden)
		})

		Convey("SetMessage", func() {
			err = err.SetMessage("hello world")
			So(err.Detail.Message, ShouldEqual, "hello world")
		})
	})
}

func TestError_ToStatus(t *testing.T) {
	Convey("TestError_ToStatus", t, func() {
		err := NewError(codes.InvalidArgument, "InvalidArgument", "invalid argument", errors.New("error"))
		err = err.SetRequestID("test-request-id").SetStatus(http.StatusBadRequest)

		detail := StatusErrorDetail(err.ToStatus().Err(), "test-request-id")
		So(detail.Status, ShouldEqual, http.StatusBadRequest)
		So(detail.RequestID, ShouldEqual, "test-request-id")
		So(detail.Code, ShouldEqual, "InvalidArgument")
		So(detail.Message, ShouldEqual, "invalid argument")
	})
}

func TestError_Format(t *testing.T) {
	fun0 := func() error {
		return errors.New("error0")
	}
	fun1 := func() error {
		return errors.WithMessage(fun0(), "error1")
	}
	fun2 := func() error {
		return NewInternalError(fun1())
	}
	fun3 := func() error {
		return errors.WithMessage(fun2(), "error3")
	}
	fun4 := func() error {
		return NewInternalError(fun3())
	}
	fun5 := func() error {
		return errors.WithMessage(fun4(), "error5")
	}

	Convey("TestError_Format", t, func() {
		fmt.Println(reflect.TypeOf(errors.Cause(fun1())))
		fmt.Println(reflect.TypeOf(errors.Cause(fun2())))
		fmt.Println(reflect.TypeOf(errors.Cause(fun3())))
		fmt.Println(reflect.TypeOf(errors.Cause(fun4())))
		fmt.Println(reflect.TypeOf(errors.Cause(fun5())))
		fmt.Printf("%v\n", fun5())
		fmt.Printf("%+v\n", fun5())
	})
}
