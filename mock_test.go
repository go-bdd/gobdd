package gobdd

import "testing"

type mockT struct {
	errored int
}

func (m *mockT) Log(i ...interface{}) {
}

func (m *mockT) Logf(s string, i ...interface{}) {
}

func (m *mockT) Fatal(i ...interface{}) {
}

func (m *mockT) Fatalf(s string, i ...interface{}) {
}

func (m *mockT) Errorf(s string, i ...interface{}) {
	m.errored++
}

func (m *mockT) Error(i ...interface{}) {
	m.errored++
}

func (m *mockT) Parallel() {
}

func (m *mockT) Fail() {
}

func (m *mockT) Run(name string, f func(t *testing.T)) bool {
	return true
}
