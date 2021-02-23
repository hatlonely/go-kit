package alics

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParsePDSUri(t *testing.T) {
	Convey("TestParsePDSUri", t, func() {
		Convey("case 1", func() {
			info, err := ParsePDSUri("pdsid://v1/hz6/101/root/5f801d6f99da67296a44405facda9a003abc1b4e?ownerID=1023210024677934&role=aliyunRole")
			So(err, ShouldBeNil)
			So(info.Version, ShouldEqual, "v1")
			So(info.DomainID, ShouldEqual, "hz6")
			So(info.DriveID, ShouldEqual, "101")
			So(info.FileID, ShouldEqual, "5f801d6f99da67296a44405facda9a003abc1b4e")
			So(info.UserOwnerID, ShouldEqual, "1023210024677934")
			So(info.UserRole, ShouldEqual, "aliyunRole")
		})
		Convey("case 2", func() {
			info, err := ParsePDSUri("pdsid://v1/hz6/101/root/5f801d6f99da67296a44405facda9a003abc1b4e")
			So(err, ShouldBeNil)
			So(info.Version, ShouldEqual, "v1")
			So(info.DomainID, ShouldEqual, "hz6")
			So(info.DriveID, ShouldEqual, "101")
			So(info.FileID, ShouldEqual, "5f801d6f99da67296a44405facda9a003abc1b4e")
			So(info.UserOwnerID, ShouldEqual, "")
			So(info.UserRole, ShouldEqual, "")
		})
	})
}
