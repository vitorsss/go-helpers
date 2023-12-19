package test

import (
	"fmt"
	"strings"
)

type StringBuilderHelper interface {
	TestHelper
	Cleanuper

	Content() string
}

type stringBuilderHelper struct {
	sb         *strings.Builder
	clean      bool
	cleanupFns []func()
}

func NewStringBuilderHelper() StringBuilderHelper {
	return &stringBuilderHelper{
		sb:    &strings.Builder{},
		clean: false,
	}
}

func (h *stringBuilderHelper) Errorf(format string, args ...interface{}) {
	h.sb.WriteString(fmt.Sprintf(format, args...))
}

func (h *stringBuilderHelper) Fatalf(format string, args ...interface{}) {
	h.sb.WriteString(fmt.Sprintf(format, args...))
}

func (h *stringBuilderHelper) Helper() {
}

func (h *stringBuilderHelper) Cleanup(fn func()) {
	h.cleanupFns = append(h.cleanupFns, fn)
}

func (h *stringBuilderHelper) Content() string {
	if !h.clean {
		for _, fn := range h.cleanupFns {
			fn()
		}
		h.clean = true
	}
	return h.sb.String()
}
