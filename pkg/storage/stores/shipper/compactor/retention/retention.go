package retention

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/cortexproject/cortex/pkg/chunk"
	"github.com/prometheus/common/model"
	"go.etcd.io/bbolt"
)

var (
	bucketName    = []byte("index")
	chunkBucket   = []byte("chunks")
	empty         = []byte("-")
	logMetricName = "logs"
)

//  todo we want to extract interfaces for series iterator and marker

type MarkerTx interface {
	Mark(id []byte) error
}

// type Marker interface {
// 	Begin() MarkerTx
// 	Commit() error
// 	Rollback() error
// }

// todo clean up with interfaces.
// markForDelete delete index entries for expired chunk in `in` and add chunkid to delete in `marker`.
// All of this inside a single transaction.
func markForDelete(in, marker *bbolt.DB, expiration ExpirationChecker, config chunk.PeriodConfig) error {
	return in.Update(func(inTx *bbolt.Tx) error {
		bucket := inTx.Bucket(bucketName)
		if bucket == nil {
			return nil
		}
		return marker.Update(func(outTx *bbolt.Tx) error {
			// deleteChunkBucket, err := outTx.CreateBucket(chunkBucket)
			// if err != nil {
			// 	return err
			// }
			seriesMap := newUserSeriesMap()
			// Phase 1 we mark chunkID that needs to be deleted in marker DB
			c := bucket.Cursor()
			var aliveChunk bool

			// it := newBoltdbChunkIndexIterator(bucket)
			// for it.Next() {
			// 	if it.Err() != nil {
			// 		return it.Err()
			// 	}
			// 	ref := it.Entry()
			// }
			// if err := forAllChunkRef(c, func(ref *ChunkRef) error {
			// 	if expiration.Expired(ref) {
			// 		if err := deleteChunkBucket.Put(ref.ChunkID, empty); err != nil {
			// 			return err
			// 		}
			// 		seriesMap.Add(ref.SeriesID, ref.UserID)
			// 		if err := c.Delete(); err != nil {
			// 			return err
			// 		}
			// 		return nil
			// 	}
			// 	// we found a key that will stay.
			// 	aliveChunk = true
			// 	return nil
			// }); err != nil {
			// 	return err
			// }
			// shortcircuit: no chunks remaining we can delete everything.
			if !aliveChunk {
				return inTx.DeleteBucket(bucketName)
			}
			// Phase 2 verify series that have marked chunks have still other chunks in the index per buckets.
			// If not this means we can delete labels index entries for the given series.
			return seriesMap.ForEach(func(seriesID, userID []byte) error {
				// for all buckets if seek to bucket.hashKey + ":" + string(seriesID) is not nil we still have chunks.
				bucketHashes := allBucketsHashes(config, unsafeGetString(userID))
				for _, bucketHash := range bucketHashes {
					if key, _ := c.Seek([]byte(bucketHash + ":" + string(seriesID))); key != nil {
						return nil
					}
					// this bucketHash doesn't contains the given series. Let's remove it.
					if err := forAllLabelRef(c, bucketHash, seriesID, config, func(_ *LabelIndexRef) error {
						if err := c.Delete(); err != nil {
							return err
						}
						return nil
					}); err != nil {
						return err
					}
				}
				return nil
			})
		})
	})
}

func forAllLabelRef(c *bbolt.Cursor, bucketHash string, seriesID []byte, config chunk.PeriodConfig, callback func(ref *LabelIndexRef) error) error {
	// todo reuse memory and refactor
	var (
		prefix string
	)
	// todo refactor ParseLabelRef. => keyType,SeriesID
	switch config.Schema {
	case "v11":
		shard := binary.BigEndian.Uint32(seriesID) % config.RowShards
		prefix = fmt.Sprintf("%02d:%s:%s", shard, bucketHash, logMetricName)
	default:
		prefix = fmt.Sprintf("%s:%s", bucketHash, logMetricName)
	}
	for k, _ := c.Seek([]byte(prefix)); k != nil; k, _ = c.Next() {
		ref, ok, err := parseLabelIndexRef(decodeKey(k))
		if err != nil {
			return err
		}
		if !ok || !bytes.Equal(seriesID, ref.SeriesID) {
			continue
		}
		if err := callback(ref); err != nil {
			return err
		}
	}

	return nil
}

func allBucketsHashes(config chunk.PeriodConfig, userID string) []string {
	return bucketsHashes(config.From.Time, config.From.Add(config.IndexTables.Period), config, userID)
}

func bucketsHashes(from, through model.Time, config chunk.PeriodConfig, userID string) []string {
	var (
		fromDay    = from.Unix() / int64(config.IndexTables.Period/time.Second)
		throughDay = through.Unix() / int64(config.IndexTables.Period/time.Second)
		result     = []string{}
	)
	for i := fromDay; i <= throughDay; i++ {
		result = append(result, fmt.Sprintf("%s:d%d", userID, i))
	}
	return result
}
