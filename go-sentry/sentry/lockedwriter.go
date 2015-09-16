package sentry

import (
	"bytes"
	"fmt"
	"sync"
)

type lockedWriter struct {
	buf *bytes.Buffer
	lck sync.Mutex
}

func newLockedWriter() *lockedWriter {
	return &lockedWriter{
		buf: new(bytes.Buffer),
		lck: sync.Mutex{},
	}
}

func (l *lockedWriter) Printf(format string, args ...interface{}) {
	l.lck.Lock()
	defer l.lck.Unlock()
	l.buf.WriteString(fmt.Sprintf(format, args...))
}

func (l *lockedWriter) String() string {
	l.lck.Lock()
	defer l.lck.Unlock()
	return l.buf.String()
}
