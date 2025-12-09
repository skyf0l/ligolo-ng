// Ligolo-ng
// Copyright (C) 2025 Nicolas Chatelain (nicocha30)

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package utils

import (
	"testing"
)

func TestParseLigoloURL(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantHost   string
		wantScheme string
		wantErr    bool
	}{
		{
			name:       "bare IP should default to port 11601",
			input:      "192.168.1.1",
			wantHost:   "192.168.1.1:11601",
			wantScheme: "",
			wantErr:    false,
		},
		{
			name:       "bare hostname should default to port 11601",
			input:      "example.com",
			wantHost:   "example.com:11601",
			wantScheme: "",
			wantErr:    false,
		},
		{
			name:       "IP with explicit port should be preserved",
			input:      "192.168.1.1:8080",
			wantHost:   "192.168.1.1:8080",
			wantScheme: "",
			wantErr:    false,
		},
		{
			name:       "hostname with explicit port should be preserved",
			input:      "example.com:8080",
			wantHost:   "example.com:8080",
			wantScheme: "",
			wantErr:    false,
		},
		{
			name:       "https URL should be preserved",
			input:      "https://example.com",
			wantHost:   "example.com",
			wantScheme: "https",
			wantErr:    false,
		},
		{
			name:       "wss URL with port should be preserved",
			input:      "wss://example.com:443",
			wantHost:   "example.com:443",
			wantScheme: "wss",
			wantErr:    false,
		},
		{
			name:       "http URL should be preserved",
			input:      "http://example.com:80",
			wantHost:   "example.com:80",
			wantScheme: "http",
			wantErr:    false,
		},
		{
			name:       "ws URL should be preserved",
			input:      "ws://127.0.0.1:11601",
			wantHost:   "127.0.0.1:11601",
			wantScheme: "ws",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLigoloURL(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLigoloURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if got.Host != tt.wantHost {
				t.Errorf("ParseLigoloURL() Host = %v, want %v", got.Host, tt.wantHost)
			}
			if got.Scheme != tt.wantScheme {
				t.Errorf("ParseLigoloURL() Scheme = %v, want %v", got.Scheme, tt.wantScheme)
			}
		})
	}
}
