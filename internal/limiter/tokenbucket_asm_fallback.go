// Copyright (c) BeduSec. All rights reserved.
// +build !amd64

package limiter

func takeToken(counter *uint64) uint64 {
	if *counter > 0 {
		*counter--
		return *counter
	}
	return 0
}