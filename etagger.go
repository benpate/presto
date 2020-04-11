package presto

// ETagger interface wraps the ETag function, which tells presto whether or not an object
// supports ETags.  Presto uses ETags to automatically support optimistic locking of files,
// as well as saving time and bandwidth using 304: "Not Modified" responses when possible.
type ETagger interface {

	// ETag returns a version-unique string that helps determine if an object has changed or not.
	ETag() string
}
