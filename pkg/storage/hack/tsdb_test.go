package main

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/cortexproject/cortex/pkg/chunk"
	util_log "github.com/cortexproject/cortex/pkg/util/log"
	"github.com/grafana/loki/pkg/storage/tsdb"
	"github.com/grafana/loki/pkg/storage/tsdb/chunkenc"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/stretchr/testify/require"
)

func Test_Blocks(t *testing.T) {
	temp := t.TempDir()

	writer, err := tsdb.NewBlockWriter(util_log.Logger, temp, 10e6)
	if err != nil {
		t.Fatal(err)
	}
	app := writer.Appender(context.Background())
	for i := 0; i < 10e3; i++ {
		_, err = app.Append(0, chunk.BenchmarkLabels, time.Now().UnixNano(), []byte(fmt.Sprintf("log %d", i)))
	}
	if err != nil {
		t.Fatal(err)
	}

	if err := app.Commit(); err != nil {
		t.Fatal(err)
	}
	id, err := writer.Flush(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(temp)
	ids := id.String()
	t.Log(ids)

	bk, err := tsdb.OpenBlock(util_log.Logger, filepath.Join(temp, ids), chunkenc.NewPool())
	if err != nil {
		t.Fatal(err)
	}

	t.Log(bk.LabelNames())

	q, err := tsdb.NewBlockChunkQuerier(bk, bk.Meta().MinTime, bk.Meta().MaxTime)
	if err != nil {
		t.Fatal(err)
	}

	var (
		actual []struct {
			ts int64
			l  string
		}
		total int
	)
	set := q.Select(false, nil, labels.MustNewMatcher(labels.MatchEqual, "container_name", "some-name"))
	for set.Next() {
		cs := set.At()
		cit := cs.Iterator()
		for cit.Next() {
			meta := cit.At()
			it := meta.Chunk.Iterator(nil)
			for it.Next() {
				ts, l := it.At()
				total++
				actual = append(actual, struct {
					ts int64
					l  string
				}{ts, string(l)})
			}
		}
	}
	require.Equal(t, len(actual), 10e3)
	require.Equal(t, total, 10e3)
}
