package reports

import (
	"reflect"
	"testing"
)

func TestDefectStatusTypes(t *testing.T) {

}

func expect(t *testing.T, expect interface{}, actual interface{}, msg string) {
	if !reflect.DeepEqual(expect, actual) {
		t.Errorf("%s - expect: %v, but got: %v", msg, expect, actual)
	}
}
