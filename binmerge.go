package binmerge

import (
	"path"
	"path/filepath"
)

// Merge creates a unified bin file and a new cue file and writes it in the given destination
func Merge(cuePath string, destination string) error {
	cueMap, err := cueToCueMap(cuePath)
	if err != nil {
		return err
	}

	cueDir := path.Dir(cuePath)
	cueName := path.Base(cuePath)[0 : len(path.Base(cuePath))-len(filepath.Ext(path.Base(cuePath)))]
	destination = filepath.Join(destination, cueName)

	err = createCuesheet(
		cueName,
		cueDir,
		cueMap,
		destination,
	)
	if err != nil {
		return err
	}

	err = createMergedBin(cueMap, cueDir, destination, cueName)
	if err != nil {
		return err
	}

	return nil
}
