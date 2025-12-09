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
		// net.SplitHostPort doesn't return typed errors, so we check the error message.
		// This is a common pattern when dealing with the standard library.
		// The error messages have been stable across Go versions.
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
		} else {
			// For other unexpected errors, return them
			return nil, err
		}
	}

	u, err := url.Parse("//" + trimmed)
	if err != nil {
		return nil, err
	}

	return &LigoloURL{u}, nil
}
