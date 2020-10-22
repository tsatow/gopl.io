package tar

import (
	".."
	"archive/tar"
	"io"
	"os"
)

func Reader(in io.Reader) *io.Reader {
	reader := tar.NewReader(in)
	f, _ := os.Open("")
	f.Read()
	return nil
}

func init() {
	archive.RegisterFormat("tar", "ustar\000", Reader)
}
