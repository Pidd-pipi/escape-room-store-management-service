// Package handler implements the HTTP handlers for escape-room-ops.
package handler

import "strconv"

func parseUint(s string) (uint, error) {
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(n), nil
}
