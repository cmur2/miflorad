package main

// program version, will be populated on build
var version string

func getVersion() string {
	if version == "" {
		return "dev"
	} else {
		return version
	}
}
