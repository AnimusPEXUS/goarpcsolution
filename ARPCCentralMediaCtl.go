package goarpcsolution

import (
	"sync"
)

type ARPCCentralMediaCtl struct {
	ids      []string
	ids_lock *sync.Mutex
}
