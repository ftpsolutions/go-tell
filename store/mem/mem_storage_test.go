package memstorage

import (
	"testing"

	"github.com/ftpsolutions/go-tell/store"
	"github.com/ftpsolutions/go-tell/store/storetest"
)

func TestStoreSuite(t *testing.T) {
	storetest.StorageSuite(func() (store.Storage, func()) {
		return store.Storage(Open()), func() {}
	}, t)
}
