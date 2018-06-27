package boltstorage

import (
	"testing"

	"github.com/Flaque/filet"
	"github.com/kithix/go-tell/store"
	"github.com/kithix/go-tell/store/storetest"
)

// Utility to create a temporary boltstore with a cleanup function
// that should be run once the test is complete.
func createTestBoltStore(t *testing.T) (*BoltStore, func()) {
	// Create a temporary directory on the system
	dir := filet.TmpDir(t, "")
	// Use our temporary file for testing
	dbPath := dir + "test.db"
	t.Log("Using test file at: " + dbPath)
	// TODO implement test logger utility?
	boltStore, err := Open(dbPath, "testjobs", nil, nil)
	cleanUpFunc := func() {
		filet.CleanUp(t)
		if err == nil {
			boltErr := boltStore.Close()
			if boltErr != nil {
				t.Error(boltErr)
			}
		}
	}
	if err != nil {
		t.Error(err)
		return nil, cleanUpFunc
	}
	return boltStore, cleanUpFunc
}

func TestStoreSuite(t *testing.T) {
	storetest.StorageSuite(func() (store.Storage, func()) {
		s, cleanUpFunc := createTestBoltStore(t)
		return store.Storage(s), cleanUpFunc
	}, t)
}
