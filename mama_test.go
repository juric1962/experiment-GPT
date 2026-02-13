package main

import (
	"testing"
)

func TestDefaultMama(t *testing.T) {
	mama := defaultMama()

	if len(mama) != 1 {
		t.Errorf("expected 1 mama, got %d", len(mama))
	}

	if mama["1"].Name != "Mama" {
		t.Errorf("expected name 'Mama', got '%s'", mama["1"].Name)
	}

	if mama["1"].Age != 50 {
		t.Errorf("expected age 50, got %d", mama["1"].Age)
	}

	if mama["1"].Work != "Housewife" {
		t.Errorf("expected work 'Housewife', got '%s'", mama["1"].Work)
	}
}

func TestSnapshotMama(t *testing.T) {
	snapshot := snapshotMama()

	if len(snapshot) != 1 {
		t.Errorf("expected 1 mama in snapshot, got %d", len(snapshot))
	}

	if snapshot["1"].Name != "Mama" {
		t.Errorf("expected snapshot name 'Mama', got '%s'", snapshot["1"].Name)
	}

	if snapshot["1"].Age != 50 {
		t.Errorf("expected snapshot age 50, got %d", snapshot["1"].Age)
	}
}

func TestMamaStruct(t *testing.T) {
	mama := Mama{
		Name: "TestMama",
		Age:  45,
		Work: "Engineer",
	}

	if mama.Name != "TestMama" {
		t.Errorf("expected name 'TestMama', got '%s'", mama.Name)
	}

	if mama.Age != 45 {
		t.Errorf("expected age 45, got %d", mama.Age)
	}

	if mama.Work != "Engineer" {
		t.Errorf("expected work 'Engineer', got '%s'", mama.Work)
	}
}
