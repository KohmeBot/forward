package forward

import (
	"testing"
)

func TestSplit(t *testing.T) {

	val, ats := parseAtInfo("@456@123，好")
	t.Log(val, ats)
}
