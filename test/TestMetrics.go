package test

import (
	"testing"
)

type MethodCallMetrics struct {
	method        string
	expectedCalls int
	actualCalls   int
}

func (m *MethodCallMetrics) Call() {
	m.actualCalls++
}

func (m *MethodCallMetrics) Report(t *testing.T, index int) {
	if m.expectedCalls != m.actualCalls {
		t.Errorf("%d - Calls for %s: Expected %d -> Actual %d", index, m.method, m.expectedCalls, m.actualCalls)
	}
}
