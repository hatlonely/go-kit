package strx

import "regexp"

var ReIdentifier = regexp.MustCompile(`^[a-zA-Z]\w+$`)
var ReEmail = regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,4}$`)
var RePhone = regexp.MustCompile(`^1[345789][0-9]{9}$`)
