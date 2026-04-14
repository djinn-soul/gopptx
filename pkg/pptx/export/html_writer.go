package export

import "io"

func writeString(w io.Writer, s string) error {
	_, err := io.WriteString(w, s)
	return err
}
