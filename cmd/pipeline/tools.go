package pipeline

// Tools are normally defined in json for models
// to help models understand its new capability.
type Tool struct {
	// Descriptions are normally needed to
	// explain what the tool is how its used.
	Description string `json: description`
}

// Most tools have to be set with a json string defining them.
// In order to not hardcode json strings we should create structs,
// then encode them into json strings during comptime or runtime.
// TODO look into go generate
type Tools []string

func tools() {}
