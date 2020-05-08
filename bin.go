package binmerge

import (
	"io"
	"os"
	"path/filepath"
)

func createMergedBin(cueMap []cueBinFile, basePath string, destination string, binName string) error {
	_ = os.MkdirAll(destination, os.ModePerm)

	f, err := os.Create(filepath.Join(destination, binName+".bin"))
	if err != nil {
		return err
	}
	defer f.Close()

	const BufferSize = 1024 * 1024
	buffer := make([]byte, BufferSize)

	for _, cueBinFile := range cueMap {
		file, err := os.Open(filepath.Join(basePath, cueBinFile.File))
		if err != nil {
			return err
		}
		defer file.Close()

		for {
			bytesread, err := file.Read(buffer)
			if err != nil {
				if err != io.EOF {
					return err
				}

				break
			}

			f.Write(buffer[:bytesread])
		}
	}

	f.Sync()

	return nil
}
