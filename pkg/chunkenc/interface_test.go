package chunkenc

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/tsdb"
)

func TestParseEncoding(t *testing.T) {
	tests := []struct {
		enc     string
		want    Encoding
		wantErr bool
	}{
		{"gzip", EncGZIP, false},
		{"bad", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.enc, func(t *testing.T) {
			got, err := ParseEncoding(tt.enc)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseEncoding() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseEncoding() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_TSDB(t *testing.T) {
	db, err := tsdb.Open(t.TempDir(), log.NewJSONLogger(os.Stdout), prometheus.DefaultRegisterer, tsdb.DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	db.DisableCompactions()

	app := db.Appender(context.Background())
	if _, err := app.Append(0, labels.FromMap(map[string]string{"foo": "bar"}), 10, 1); err != nil {
		t.Fatal(app)
	}
	if err := app.Commit(); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	if err := filepath.Walk(t.TempDir(),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			log.NewJSONLogger(os.Stdout).Log("path", path, "size", info.Size())
			return nil
		}); err != nil {
		t.Fatal(err)
	}
}
