package boltstorage

import (
	"testing"
	"time"

	"github.com/Flaque/filet"
	"github.com/boltdb/bolt"

	gotell "github.com/ftpsolutions/go-tell"
	"github.com/ftpsolutions/go-tell/store/storetest"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetJobWithRetryCount(t *testing.T) {
	// Setup: Mock a BoltDB database and seed it with jobs
	dbPath := "test.db"
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Create some jobs with different retry counts
	job1 := &gotell.Job{
		ID:         uuid.Must(uuid.NewV4()),
		RetryCount: 0,
		Status:     gotell.StatusJobCreated,
	}
	job2 := &gotell.Job{
		ID:         uuid.Must(uuid.NewV4()),
		RetryCount: 5,
		Status:     gotell.StatusJobCreated,
	}
	job3 := &gotell.Job{
		ID:         uuid.Must(uuid.NewV4()),
		RetryCount: 3,
		Status:     gotell.StatusJobCreated,
	}

	store := &BoltStore{
		db:         db,
		bucketName: []byte("testBucket"),
	}

	// Seed the database with jobs
	store.AddJob(job1)
	store.AddJob(job2)
	store.AddJob(job3)

	// Get the first job; should be job1 since it has the lowest retry count
	retrievedJob, err := store.GetJob()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, job1.ID, retrievedJob.ID, "Expected the job with the lowest retry count")

	// Now let's say job1 fails and its retry count increases
	job1.RetryCount = 6
	store.UpdateJob(job1)

	// Get the job again; this time it should be job3 since job1's retry count is increased
	retrievedJob, err = store.GetJob()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, job3.ID, retrievedJob.ID, "Expected the job with the lowest retry count")
}

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
