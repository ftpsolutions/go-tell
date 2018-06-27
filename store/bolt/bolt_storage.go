package boltstorage

import (
	"encoding/json"
	"log"

	"github.com/boltdb/bolt"
	"github.com/kithix/go-tell/store"
	uuid "github.com/satori/go.uuid"
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

	return boltStore, nil
}

func (s *BoltStore) Close() error {
	return s.db.Close()
}

func (s *BoltStore) AddJob(job *store.Job) error {
	// Get byte from struct
	id, data, err := getByteDataFromJob(job)
	if err != nil {
		return err
	}

	return s.write(func(bucket *bolt.Bucket) error {
		return bucket.Put(id, data)
	})
}

func (s *BoltStore) GetJob() (*store.Job, error) {
	job := &store.Job{}
	noJob := true
	err := s.read(func(bucket *bolt.Bucket) error {
		c := bucket.Cursor()
		for idBytes, data := c.First(); idBytes != nil; idBytes, data = c.Next() {
			job = &store.Job{} // Reset job.
			err := json.Unmarshal(data, job)
			if err != nil {
				// Try next job
				s.logger.Println("Unable to unmarshal job from boltDB", err)
				continue
			}
			if job.Status == store.StatusJobCreated {
				id, err := uuid.FromBytes(idBytes)
				if err != nil {
					s.logger.Println("Failed to convert id bytes to UUID")
					// Try next job
					continue
				}
				job.ID = id
				noJob = false
				break
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if noJob {
		return nil, store.ErrorNoJobFound
	}

	// Set the job to being done.
	job.Status = store.StatusJobPending
	err = s.UpdateJob(job)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (s *BoltStore) UpdateJob(job *store.Job) error {
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
			return store.ErrorNoJobFound
		}
		return bucket.Put(id, data)
	})
}

func (s *BoltStore) DeleteJob(job *store.Job) error {
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

func getByteIDFromJob(job *store.Job) []byte {
	return job.ID.Bytes()
}

func getByteDataFromJob(job *store.Job) ([]byte, []byte, error) {
	data, err := json.Marshal(job)
	if err != nil {
		return nil, nil, err
	}
	id := getByteIDFromJob(job)
	return id, data, nil
}
