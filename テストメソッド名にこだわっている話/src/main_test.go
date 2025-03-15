package main

import "testing"

func Test_合計値を計算できる(t *testing.T) {
	got := Add(1, 2)
	want := 3
	if got != want {
		t.Errorf("Add(1, 2) = %d; want %d", got, want)
	}
}
