package boltstorage

import (
	"encoding/json"
	"log"
	"math"

	"github.com/boltdb/bolt"
	"github.com/gofrs/uuid"

	gotell "github.com/ftpsolutions/go-tell"
)

type BoltStore struct {
	bucketName []byte
	logger     *log.Logger
	db         *bolt.DB
}

func Open(filePath string, bucketName string, opts *bolt.Options, logger *log.Logger) (*BoltStore, error) {
	database, err := bolt.Open(filePath, 0600, opts)
	if err != nil {
		return nil, err
	}

	boltStore := &BoltStore{
		bucketName: []byte(bucketName),
		logger:     logger,
		db:         database,
	}

	err = boltStore.resetPendingJobs()
	if err != nil {
		return nil, err
	}

	return boltStore, nil
}

func (s *BoltStore) resetPendingJobs() error {
	return s.write(func(bucket *bolt.Bucket) error {
		c := bucket.Cursor()
		var err error
		// For each job in the store, reset all pending jobs to created.
		for idBytes, data := c.First(); idBytes != nil; idBytes, data = c.Next() {
			job := &gotell.Job{}
			err = json.Unmarshal(data, job)
			if err != nil {
				return err
			}

			if job.Status != gotell.StatusJobPending {
				continue
			}

			job.Status = gotell.StatusJobCreated
			_, updatedJobData, err := getByteDataFromJob(job)
			if err != nil {
				return err
			}
			err = bucket.Put(idBytes, updatedJobData)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *BoltStore) Close() error {
	return s.db.Close()
}

func (s *BoltStore) AddJob(job *gotell.Job) error {
	// Get byte from struct
	id, data, err := getByteDataFromJob(job)
	if err != nil {
		return err
	}

	return s.write(func(bucket *bolt.Bucket) error {
		return bucket.Put(id, data)
	})
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (s *BoltStore) GetJob() (*gotell.Job, error) {
	job := &gotell.Job{}
	noJob := true
	lowestRetryCount := -1
	err := s.read(func(bucket *bolt.Bucket) error {
		c := bucket.Cursor()
		// Find the lowest retry count in the job queue
		for idBytes, data := c.First(); idBytes != nil; idBytes, data = c.Next() {
			job = &gotell.Job{} // Reset job.
			err := json.Unmarshal(data, job)
			if err != nil {
				// Try next job
				s.logger.Println("Unable to unmarshal job from boltDB", err)
				continue
			}
			if job.Status == gotell.StatusJobCreated {
				if lowestRetryCount == -1 {
					lowestRetryCount = job.RetryCount
				} else {
					lowestRetryCount = Min(lowestRetryCount, job.RetryCount)
				}
			}
		}
		// Find the job with the lowest retry count
		for idBytes, data := c.First(); idBytes != nil; idBytes, data = c.Next() {
			job = &gotell.Job{}
			err := json.Unmarshal(data, job)
			if err != nil {
				// Try next job
				s.logger.Println("Unable to unmarshal job from boltDB", err)
				continue
			}
			if job.Status == gotell.StatusJobCreated {
				id, err := uuid.FromBytes(idBytes)
				if err != nil {
					s.logger.Println("Failed to convert id bytes to UUID")
					// Try next job
					continue
				}
				job.ID = id
				noJob = false
				if job.RetryCount == lowestRetryCount {
					break
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if noJob {
		return nil, gotell.ErrorNoJobFound
	}

	// Set the job to being done.
	job.Status = gotell.StatusJobPending
	err = s.UpdateJob(job)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (s *BoltStore) UpdateJob(job *gotell.Job) error {
	// Get byte from struct
	id, data, err := getByteDataFromJob(job)
	if err != nil {
		return err
	}

	return s.write(func(bucket *bolt.Bucket) error {
		// Look through DB for key.
		c := bucket.Cursor()
		key, _ := c.Seek(id)
		if key == nil {
			// TODO should this behaviour be upsert?
			return gotell.ErrorNoJobFound
		}
		return bucket.Put(id, data)
	})
}

func (s *BoltStore) DeleteJob(job *gotell.Job) error {
	return s.write(func(bucket *bolt.Bucket) error {
		return bucket.Delete(getByteIDFromJob(job))
	})
}

type withBucketfunc func(*bolt.Bucket) error

func (s *BoltStore) read(wb withBucketfunc) error {
	return s.db.View(s.with(wb))
}

func (s *BoltStore) write(wb withBucketfunc) error {
	return s.db.Update(s.with(wb))
}

func (s *BoltStore) with(wb withBucketfunc) func(*bolt.Tx) error {
	return func(tx *bolt.Tx) error {
		bucket, err := s.getBucket(tx)
		if err != nil {
			return err
		}
		return wb(bucket)
	}
}

func (s *BoltStore) getBucket(tx *bolt.Tx) (*bolt.Bucket, error) {
	bucket := tx.Bucket(s.bucketName)
	var err error
	if bucket == nil {
		bucket, err = tx.CreateBucket(s.bucketName)
		if err != nil {
			return nil, err
		}
	}
	return bucket, nil
}

func getByteIDFromJob(job *gotell.Job) []byte {
	return job.ID.Bytes()
}

func getByteDataFromJob(job *gotell.Job) ([]byte, []byte, error) {
	data, err := json.Marshal(job)
	if err != nil {
		return nil, nil, err
	}
	id := getByteIDFromJob(job)
	return id, data, nil
}
