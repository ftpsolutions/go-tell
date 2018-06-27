package memstorage

import (
	"testing"

	"github.com/kithix/go-tell/store"
	"github.com/kithix/go-tell/store/storetest"
)

func TestStoreSuite(t *testing.T) {
	storetest.StorageSuite(func() (store.Storage, func()) {
		return store.Storage(Open()), func() {}
	}, t)
}
