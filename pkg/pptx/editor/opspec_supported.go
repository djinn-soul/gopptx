package editor

// SupportedOps returns the full list of operation codes handled by the bridge.
func SupportedOps() []string {
	ops := supportedSlideAndMetaOps()
	ops = append(ops, supportedContentOps()...)
	return ops
}
