package pipeline

// Most tools have to be set with a json string defining them.
// In order to not hardcode json strings we should create structs, 
// then encode them into json strings during comptime or runtime.
// TODO look into go generate
type Tools []string

func tools() {

}
