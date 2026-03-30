/*
Copyright © 2026 fakeapate <fakeapate@pm.me>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package mullvadProxy

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"
)

// metadata stores the ETag and Last-Modified values from the last successful fetch,
// persisted to disk so freshness checks survive restarts.
type metadata struct {
	ETag         string    `json:"etag"`
	LastModified time.Time `json:"last_modified"`
}

func update(cfg MullvadConfig) error {
	// TODO make AmIConnected api call first to check if connected to mullvad
	// if not log a warning
	resp, err := http.Head(cfg.RelayURL)
	if err != nil {
		return err
	}
	resp.Body.Close()

	remoteETag := resp.Header.Get("ETag")
	remoteLastModified, _ := time.Parse(time.RFC1123, resp.Header.Get("Last-Modified"))

	meta := loadMeta(cfg)

	remoteIsNewer := (remoteETag != "" && remoteETag != meta.ETag) ||
		(!remoteLastModified.IsZero() && remoteLastModified.After(meta.LastModified))

	if remoteIsNewer || Relays == nil {
		if err := fetch(cfg, remoteETag, remoteLastModified); err != nil {
			return err
		}
	} else {
		return nil
	}

	parse(cfg)
	return nil
}

func fetch(cfg MullvadConfig, etag string, lastModified time.Time) error {
	resp, err := http.Get(cfg.RelayURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := os.WriteFile(cfg.DataFile, data, 0644); err != nil {
		return err
	}

	metaData, err := json.Marshal(metadata{ETag: etag, LastModified: lastModified})
	if err != nil {
		return err
	}
	return os.WriteFile(cfg.MetaFile, metaData, 0644)
}

func parse(cfg MullvadConfig) error {
	data, err := os.ReadFile(cfg.DataFile)
	if err != nil {
		return err
	}
	var relays MullvadRelays
	if err := json.Unmarshal(data, &relays); err != nil {
		return err
	}
	Relays = &relays
	return nil
}

func loadMeta(cfg MullvadConfig) metadata {
	data, err := os.ReadFile(cfg.MetaFile)
	if err != nil {
		return metadata{}
	}
	var meta metadata
	json.Unmarshal(data, &meta)
	return meta
}
