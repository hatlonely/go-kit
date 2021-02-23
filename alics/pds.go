package alics

import (
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type PDSUriInfo struct {
	Version      string
	DomainID     string
	DriveID      string
	ParentFileID string
	FileID       string
	UserOwnerID  string
	UserRole     string
}

func ParsePDSUri(uri string) (*PDSUriInfo, error) {
	info, err := url.Parse(uri)
	if err != nil {
		return nil, errors.Wrapf(err, "Parse failed, uri: [%v]", uri)
	}
	if info.Scheme != "pdsid" {
		return nil, errors.Errorf("InvalidScheme, scheme: [%v]", info.Scheme)
	}
	version := info.Host
	if version != "v1" {
		return nil, errors.Errorf("InvalidVersion, version: [%v]", version)
	}
	vs := strings.Split(info.Path, "/")
	if len(vs) != 5 {
		return nil, errors.Errorf("len should equal 5, vs: [%v]", vs)
	}

	return &PDSUriInfo{
		Version:      version,
		DomainID:     vs[1],
		DriveID:      vs[2],
		ParentFileID: vs[3],
		FileID:       vs[4],
		UserOwnerID:  info.Query().Get("ownerID"),
		UserRole:     info.Query().Get("role"),
	}, nil
}
