package test

import "context"

type TestReporter interface {
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

type TestHelper interface {
	TestReporter
	Helper()
}

type Cleanuper interface {
	Cleanup(func())
}

func IsCleanuper(t TestReporter) (Cleanuper, bool) {
	tr := unwrapTestReporter(t)
	c, ok := tr.(Cleanuper)
	return c, ok
}

func unwrapTestReporter(t TestReporter) TestReporter {
	tr := t
	switch nt := t.(type) {
	case *cancelReporter:
		tr = unwrapTestReporter(nt.t)
	case *nopTestHelper:
		tr = unwrapTestReporter(nt.t)
	default:
		// not wrapped
	}
	return tr
}

type cancelReporter struct {
	t      TestHelper
	cancel func()
}

func (r *cancelReporter) Errorf(format string, args ...interface{}) {
	r.t.Errorf(format, args...)
}

func (r *cancelReporter) Fatalf(format string, args ...interface{}) {
	defer r.cancel()
	r.t.Fatalf(format, args...)
}

func (r *cancelReporter) Helper() {
	r.t.Helper()
}

func WithContext(ctx context.Context, t TestReporter) (TestHelper, context.Context) {
	h, ok := t.(TestHelper)
	if !ok {
		h = NopTestHelper(t)
	}

	ctx, cancel := context.WithCancel(ctx)
	return &cancelReporter{t: h, cancel: cancel}, ctx
}

type nopTestHelper struct {
	t TestReporter
}

func (h *nopTestHelper) Errorf(format string, args ...interface{}) {
	h.t.Errorf(format, args...)
}

func (h *nopTestHelper) Fatalf(format string, args ...interface{}) {
	h.t.Fatalf(format, args...)
}

func (h nopTestHelper) Helper() {}

func NopTestHelper(t TestReporter) TestHelper {
	return &nopTestHelper{t: t}
}
