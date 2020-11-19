package equal

import "testing"

func TestEqual(t *testing.T) {
	if !Equal([]int{1, 2, 3}, []int{1, 2, 3}) {
		t.Errorf("Equal([]int{1, 2, 3}, []int{1, 2, 3}) returns false")
	}
	if Equal([]string{"foo"}, []string{"bar"}) {
		t.Errorf(`Equal([]string{"foo"}, []string{"bar"}) returns true`)
	}
	if !Equal([]string(nil), []string{}) {
		t.Errorf(`Equal([]string(nil), []string{}) returns false`)
	}
}