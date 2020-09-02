package version

// APIVersion increases every time a new call is added to the API. Clients should use this info
// to determine if the server supports specific features.
const APIVersion = "v1.0.0"

// Default build-time variable.
// These values are overridden via ldflags
var (
	Version = "unknown-version"
)
