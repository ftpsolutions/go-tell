package boltstorage

import (
	"testing"
	"time"

	"github.com/Flaque/filet"

	gotell "github.com/ftpsolutions/go-tell"
	"github.com/ftpsolutions/go-tell/store/storetest"
)

// Utility to create a temporary boltstore with a cleanup function
// that should be run once the test is complete.
func createTestBoltStore(t *testing.T, dir string) (*BoltStore, func()) {
	// Create a temporary directory on the system
	if dir == "" {
		dir = filet.TmpDir(t, "")
	}
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

func TestResetPendingJobs(t *testing.T) {
	boltStore, cleanupFunc := createTestBoltStore(t, "")
	defer cleanupFunc()
	boltStore.AddJob(&gotell.Job{
		ID: storetest.BuildJobID(42),
	})
	_, err := boltStore.GetJob()
	if err != nil {
		t.Error(err)
	}
	err = boltStore.resetPendingJobs()
	if err != nil {
		t.Error(err)
	}
	job, err := boltStore.GetJob()
	if err != nil {
		t.Error(err)
	}
	if job.ID != storetest.BuildJobID(42) {
		t.Error("Job ID was not the same, real problems here")
	}
}

func TestGolangTimeIsStoredInBoltCorrectly(t *testing.T) {
	dbPath := filet.TmpDir(t, "") + "test.db"
	defer filet.CleanUp(t)
	t.Log("Using test file at: " + dbPath)

	boltStore, err := Open(dbPath, "testjobs", nil, nil)
	if err != nil {
		t.Error("Error opening boltstore", err)
		return
	}
	// Add a job with a date to see if it saves properly
	boltStore.AddJob(&gotell.Job{
		ID:      storetest.BuildJobID(42),
		Created: time.Date(7, 6, 5, 4, 3, 2, 1, time.UTC),
	})
	// Close the store, reopen and get the job out
	err = boltStore.Close()
	if err != nil {
		t.Error(err)
		return
	}

	boltStore, err = Open(dbPath, "testjobs", nil, nil)
	if err != nil {
		t.Error("Error opening boltstore", err)
		return
	}

	job, err := boltStore.GetJob()
	if err != nil {
		t.Error(err)
	}

	expectedDate := time.Date(7, 6, 5, 4, 3, 2, 1, time.UTC)
	if !job.Created.Equal(expectedDate) {
		t.Error("Expected", expectedDate, "got", job.Created)
	}
}

func TestBoltAgainstStoreSuite(t *testing.T) {
	storetest.StorageSuite(func() (gotell.Storage, func()) {
		s, cleanUpFunc := createTestBoltStore(t, "")
		return gotell.Storage(s), cleanUpFunc
	}, t)
}
