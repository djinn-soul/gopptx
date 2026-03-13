//go:build !windows

package protection

func canEncryptAgile() bool {
	return false
}

func encryptAgilePackage(_ []byte, _ string) ([]byte, error) {
	return nil, errorsAgileUnavailable()
}

func errorsAgileUnavailable() error {
	return errAgileUnavailable
}
