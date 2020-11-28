package rpcx

import (
	"errors"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStatusErrorDetail(t *testing.T) {
	Convey("TestStatusErrorDetail", t, func() {
		detail := StatusErrorDetail(errors.New("hello world"), "test-request-id")
		So(detail.Status, ShouldEqual, http.StatusInternalServerError)
		So(detail.RequestID, ShouldEqual, "test-request-id")
		So(detail.Message, ShouldEqual, "hello world")
		So(detail.Code, ShouldEqual, "InternalError")
	})
}
