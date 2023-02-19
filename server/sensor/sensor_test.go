package sensor

import "testing"

func TestReadingsTypeAreValid(t *testing.T) {
	for id, rt := range AvailableReadingTypes {
		if rt.Count != len(rt.Labels) {
			t.Errorf("For reading type %v number of labels doesn't match field count", id)
		}
	}
}
