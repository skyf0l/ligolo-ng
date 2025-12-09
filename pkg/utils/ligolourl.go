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
	"fmt"
	"net"
	"net/url"
	"strings"
)

type LigoloURL struct {
	*url.URL
}

func (l *LigoloURL) IsSecure() bool {
	if l.Scheme == "https" || l.Scheme == "wss" {
		return true
	}
	return false
}

func (l *LigoloURL) IsWebsocket() bool {
	if l.Scheme == "http" || l.Scheme == "ws" || l.Scheme == "https" || l.Scheme == "wss" {
		return true
	}
	return false
}

func (l *LigoloURL) IsValid() bool {
	return l.IsWebsocket() || l.Scheme == ""
}

func ParseLigoloURL(rawURL string) (*LigoloURL, error) {
	trimmed := strings.TrimSpace(rawURL)

	// If it's a full URL (has scheme), parse normally and return.
	if strings.Contains(trimmed, "://") {
		u, err := url.Parse(trimmed)
		if err != nil {
			return nil, err
		}
		return &LigoloURL{u}, nil
	}

	// For non-scheme host[:port] form: if no port present, append default 11601.
	if _, _, err := net.SplitHostPort(trimmed); err != nil {
		// Check the specific error to determine if we need to add a port
		errStr := err.Error()
		if strings.Contains(errStr, "missing port") {
			// Plain host without port (e.g., "localhost" or "[::1]")
			// For IPv6 addresses already in brackets, append port directly
			if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
				trimmed = trimmed + ":11601"
			} else {
				trimmed = net.JoinHostPort(trimmed, "11601")
			}
		} else if strings.Contains(errStr, "too many colons") {
			// IPv6 address without brackets (e.g., "::1", "2001:db8::1")
			trimmed = "[" + trimmed + "]:11601"
		}
		// For other errors, proceed with the original trimmed value
	}

	u, err := url.Parse("//" + trimmed)
	if err != nil {
		return nil, err
	}

	return &LigoloURL{u}, nil
}

func printURL(title string, u *url.URL) {
	fmt.Printf("--- %s ---\n", title)
	if u == nil {
		fmt.Println("  URL is nil")
		fmt.Println()
		return
	}
	fmt.Printf("  String() : %s\n", u.String())
	fmt.Printf("  Scheme   : %q\n", u.Scheme)
	fmt.Printf("  Host     : %q\n", u.Host)
	fmt.Printf("  Path     : %q\n", u.Path)
	fmt.Printf("  Opaque   : %q\n", u.Opaque)
	fmt.Println()
}

func main() {
	urlsToTest := []string{
		"wss://foo.bar:8080/path/to",
		"https://foo.bar",
		"ws://127.0.0.1:11601",
		"127.0.0.1:11601",
	}

	for _, raw := range urlsToTest {
		fmt.Printf("Processing: %q\n", raw)
		u, err := ParseLigoloURL(raw)
		if err != nil {
			fmt.Printf("Error: %v\n\n", err)
			continue
		}
		printURL(raw, u.URL)
	}
}
