package test

type MethodCallMetrics struct {
	expectedCalls int
	actualCalls   int
}

func (m *MethodCallMetrics) Call() {
	m.actualCalls++
}
