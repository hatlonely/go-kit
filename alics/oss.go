package alics

import (
	"regexp"
	"unicode/utf8"

	"github.com/pkg/errors"
)

var ossPathRegex = regexp.MustCompile(`^oss://([^/]+)/([\s\S]*)$`)

type OSSUriInfo struct {
	Bucket string
	Object string
}

func ParseOSSUri(uri string) (*OSSUriInfo, error) {
	if !utf8.ValidString(uri) {
		return nil, errors.New("InvalidOSSPath")
	}

	res := ossPathRegex.FindStringSubmatch(uri)
	if len(res) != 3 {
		return nil, errors.New("InvalidOSSPath")
	}

	return &OSSUriInfo{
		Bucket: res[1],
		Object: res[2],
	}, nil
}
