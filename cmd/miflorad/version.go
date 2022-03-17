package main

import (
	"runtime/debug"
)

func getVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}

	for _, kv := range info.Settings {
		switch kv.Key {
		case "vcs.revision":
			return kv.Value[0:8]
		}
	}

	return "unknown"
}
