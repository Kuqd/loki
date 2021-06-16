package main

import (
	"context"
	"testing"
	"time"

	"github.com/cortexproject/cortex/pkg/chunk"
	util_log "github.com/cortexproject/cortex/pkg/util/log"
	"github.com/grafana/loki/pkg/storage/tsdb"
)

func Test_Blocks(t *testing.T) {
	temp := t.TempDir()

	writer, err := tsdb.NewBlockWriter(util_log.Logger, temp, 10e6)
	if err != nil {
		t.Fatal(err)
	}
	app := writer.Appender(context.Background())
	for i := 0; i < 10e6; i++ {
		_, err = app.AppendLog(0, chunk.BenchmarkLabels, time.Now().UnixNano(), []byte(`foo`))
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
	t.Log(id)
}
