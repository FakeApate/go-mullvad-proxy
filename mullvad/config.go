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
package mullvad

import "time"

// MullvadConfig holds all Mullvad VPN configuration.
type MullvadConfig struct {
	RelayURL       string        `toml:"relay_url"`       // Mullvad wireguard relay list API endpoint
	DataFile       string        `toml:"data_file"`       // path to the cached relay list on disk
	MetaFile       string        `toml:"meta_file"`       // path to the cached relay list metadata on disk
	ProxyPort      int           `toml:"proxy_port"`      // SOCKS5 port on Mullvad relay hosts
	UpdateInterval time.Duration `toml:"update_interval"` // how often to refresh the relay list
}

// DefaultMullvadConfig returns a [MullvadConfig] with sensible defaults.
func DefaultMullvadConfig() MullvadConfig {
	return MullvadConfig{
		RelayURL:       "https://api.mullvad.net/public/relays/wireguard/v2",
		DataFile:       "relays.json",
		MetaFile:       "relays.meta.json",
		ProxyPort:      1080,
		UpdateInterval: time.Hour,
	}
}
