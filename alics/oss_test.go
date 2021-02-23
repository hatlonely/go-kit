package alics

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParseOSSPath(t *testing.T) {
	Convey("TestParseOSSPath", t, func() {
		{
			res, err := ParseOSSUri(`oss://test-res.Bucket/锁定-200403-（开放ISV两方）互联网公交合作协议?(定稿).docx`)
			So(err, ShouldBeNil)
			So(res.Bucket, ShouldEqual, "test-res.Bucket")
			So(res.Object, ShouldEqual, "锁定-200403-（开放ISV两方）互联网公交合作协议?(定稿).docx")
		}
		{
			res, err := ParseOSSUri(`oss://test-res.Bucket/a.docx?q=123`)
			So(err, ShouldBeNil)
			So(res.Bucket, ShouldEqual, "test-res.Bucket")
			So(res.Object, ShouldEqual, "a.docx?q=123")
		}
		{
			res, err := ParseOSSUri(`oss://test-res.Bucket/a/b/c.docx`)
			So(err, ShouldBeNil)
			So(res.Bucket, ShouldEqual, "test-res.Bucket")
			So(res.Object, ShouldEqual, "a/b/c.docx")
		}
		{
			res, err := ParseOSSUri(`oss://test-res.Bucket///a/b/c.docx`)
			So(err, ShouldBeNil)
			So(res.Bucket, ShouldEqual, "test-res.Bucket")
			So(res.Object, ShouldEqual, "//a/b/c.docx")
		}
		{
			res, err := ParseOSSUri(`oss://test-res.Bucket/👌%2F.docx`)
			So(err, ShouldBeNil)
			So(res.Bucket, ShouldEqual, "test-res.Bucket")
			So(res.Object, ShouldEqual, "👌%2F.docx")
		}
		{
			res, err := ParseOSSUri("oss://test-res.Bucket/👌%2F\n.docx")
			So(err, ShouldBeNil)
			So(res.Bucket, ShouldEqual, "test-res.Bucket")
			So(res.Object, ShouldEqual, "👌%2F\n.docx")
		}
	})
}
