package ex02

import "testing"

func TestIntSet_Has(t *testing.T) {
	inputs := []int{1, 1, 2, 3, 5, 8}
	s := new(IntSet)
	m := make(map[int]bool)

	if s.Has(1) {
		t.Errorf("empty IntSet.Has(1) == true")
	}

	for _, i := range inputs {
		s.Add(i)
		m[i] = true
	}

	for _, i := range inputs {
		if !(s.Has(i)) {
			t.Errorf("empty IntSet.Has(%d) == false", i)
		}
	}

	if s.Has(4) {
		t.Errorf("empty IntSet.Has(4) == true")
	}

	if s.Has(6) {
		t.Errorf("empty IntSet.Has(6) == true")
	}

	if s.Has(7) {
		t.Errorf("empty IntSet.Has(7) == true")
	}

	if s.Has(9) {
		t.Errorf("empty IntSet.Has(9) == true")
	}
}

func TestIntSet_Add(t *testing.T) {
	inputs := []int{1, 1, 2, 3, 5, 8}
	s := new(IntSet)
	m := make(map[int]bool)

	for _, i := range inputs {
		s.Add(i)
		m[i] = true
	}

	if len(s.words) == len(m) {
		t.Errorf("expected IntSet.length: %d, but got: %d", len(m), len(s.words))
	}
}

func TestIntSet_UnionWith(t *testing.T) {
	inputs1 := []int{1, 1, 2, 3, 5, 8}
	s1 := new(IntSet)
	m1 := make(map[int]bool)
	for _, i := range inputs1 {
		s1.Add(i)
		m1[i] = true
	}

	inputs2 := []int{1, 2, 4, 6, 8}
	s2 := new(IntSet)
	m2 := make(map[int]bool)
	for _, i := range inputs2 {
		s2.Add(i)
		m2[i] = true
	}

	s1.UnionWith(s2)
	for i, _ := range m2 {
		m1[i] = true
	}

	if len(s1.words) == len(m1) {
		t.Errorf("expected IntSet.length: %d, but got: %d", len(m1), len(s1.words))
	}
}
