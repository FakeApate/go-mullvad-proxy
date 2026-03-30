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
package endpoints

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Relays holds the most recently loaded Mullvad relay list.
// It is nil until the first successful parse completes.
var Relays *MullvadRelays

// StartUpdater performs an initial relay list update and then checks for updates
// at the interval specified by cfg.UpdateInterval.
// It returns immediately after the first update; subsequent checks run in the background.
func StartUpdater(cfg MullvadConfig) {
	update(cfg)
	go func() {
		ticker := time.NewTicker(cfg.UpdateInterval)
		for range ticker.C {
			update(cfg)
		}
	}()
}

func IsConnected() (bool, error) {
	resp, err := http.Get("https://am.i.mullvad.net/json")
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var res AmIConnected
	if err := json.Unmarshal(data, &res); err != nil {
		return false, err
	}

	return res.MullvadExitIp, nil
}

// RelayFilter specifies which relays [SelectProxies] should include.
// Zero values impose no constraint for that field.
// Active and IncludeInCountry are always required to be true and are not exposed.
type RelayFilter struct {
	Location *regexp.Regexp // match relay.Location against pattern; nil = any
	Owned    *bool          // nil = any, true = owned only, false = non-owned only
	Weight   func(int) bool // nil = no weight filter
}

// SelectProxies returns up to limit proxy strings for relays matching filter,
// formatted as socks5://hostname:<port> for use with colly's proxy switcher.
// limit <= 0 means no limit. Reports an error if the relay list is not loaded.
func SelectProxies(cfg MullvadConfig, limit int, filter RelayFilter) ([]string, error) {
	if Relays == nil {
		return nil, errors.New("relay list not loaded")
	}

	port := strconv.Itoa(cfg.ProxyPort)
	var results []string
	for _, relay := range Relays.Wireguard.Relays {
		if !relay.Active || !relay.IncludeInCountry {
			continue
		}
		if filter.Location != nil && !filter.Location.MatchString(relay.Location) {
			continue
		}
		if filter.Owned != nil {
			if props, ok := relay.AdditionalProperties.(map[string]any); ok {
				if owned, ok := props["owned"].(bool); ok && owned != *filter.Owned {
					continue
				}
			}
		}
		if filter.Weight != nil && !filter.Weight(relay.Weight) {
			continue
		}
		hostname := strings.ReplaceAll(relay.Hostname, "-wg-", "-wg-socks5-")
		results = append(results, "socks5://"+hostname+".relays.mullvad.net:"+port)
		if limit > 0 && len(results) == limit {
			break
		}
	}

	return results, nil
}
