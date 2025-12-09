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
	"net/url"
	"strconv"
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

// isLikelyHostname checks if a string looks like a hostname rather than a URL scheme
func isLikelyHostname(s string) bool {
	// Common URL schemes that should not be treated as hostnames
	knownSchemes := []string{"http", "https", "ws", "wss", "ftp", "ftps", "file", "data", "mailto", "tel"}
	for _, scheme := range knownSchemes {
		if s == scheme {
			return false
		}
	}
	// If it contains a dot, it's likely a hostname (e.g., example.com, 192.168.1.1)
	// Single-word strings (like "localhost") could be hostnames too
	return strings.Contains(s, ".") || (s != "" && !strings.ContainsAny(s, "/:@?#[]"))
}

// isNumericPort checks if a string is a valid port number
func isNumericPort(s string) bool {
	port, err := strconv.Atoi(s)
	return err == nil && port > 0 && port <= 65535
}

func ParseLigoloURL(rawURL string) (*LigoloURL, error) {
	u, err := url.Parse(rawURL)

	if err != nil {
		if urlErr, ok := err.(*url.Error); ok && strings.Contains(urlErr.Err.Error(), "first path segment") {
			u, err := url.Parse("//" + rawURL)
			if err != nil {
				return nil, err
			}
			return &LigoloURL{u}, nil
		}
		return nil, err
	}

	// If parsed with scheme but no Host and has Opaque, it might be host:port mistaken as scheme:opaque
	// Only reparse if the scheme looks like a hostname and opaque looks like a port
	if u.Scheme != "" && u.Host == "" && u.Opaque != "" && isLikelyHostname(u.Scheme) && isNumericPort(u.Opaque) {
		// Try parsing with // prefix
		u, err = url.Parse("//" + rawURL)
		if err != nil {
			return nil, err
		}
	}

	// If no scheme and the Host is empty (bare hostname without port), default to :11601
	if u.Scheme == "" && u.Host == "" && u.Path != "" && !strings.Contains(u.Path, "/") {
		// This is a bare hostname, append :11601
		u, err = url.Parse("//" + rawURL + ":11601")
		if err != nil {
			return nil, err
		}
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
