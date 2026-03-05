package slide

type PartLookup interface {
	Get(partPath string) ([]byte, bool)
}
