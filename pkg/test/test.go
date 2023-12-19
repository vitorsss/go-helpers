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
		tr = unwrapTestReporter(nt.TestHelper)
	case *nopTestHelper:
		tr = unwrapTestReporter(nt.TestReporter)
	default:
		// not wrapped
	}
	return tr
}

type cancelReporter struct {
	TestHelper
	cancel func()
}

func (r *cancelReporter) Fatalf(format string, args ...interface{}) {
	defer r.cancel()
	r.TestHelper.Fatalf(format, args...)
}

func AsHelper(t TestReporter) TestHelper {
	h, ok := t.(TestHelper)
	if !ok {
		h = NopTestHelper(t)
	}
	return h
}

func WithContext(ctx context.Context, t TestReporter) (TestHelper, context.Context) {
	h := AsHelper(t)

	ctx, cancel := context.WithCancel(ctx)
	return &cancelReporter{TestHelper: h, cancel: cancel}, ctx
}

type nopTestHelper struct {
	TestReporter
}

func (h nopTestHelper) Helper() {}

func NopTestHelper(t TestReporter) TestHelper {
	return &nopTestHelper{TestReporter: t}
}
