package mem

import (
	"testing"

	gotell "github.com/ftpsolutions/go-tell"
	"github.com/ftpsolutions/go-tell/store/storetest"
)

func TestMemStorageAgainstStoreSuite(t *testing.T) {
	storetest.StorageSuite(func() (gotell.Storage, func()) {
		return gotell.Storage(Open()), func() {}
	}, t)
}
