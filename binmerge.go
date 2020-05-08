package binmerge

import (
	"path/filepath"
)

// Merge creates a unified bin file and a new cue file and writes it in the given destination
func Merge(cuePath string, destination string) error {
	cueMap, err := cueToCueMap(cuePath)
	if err != nil {
		return err
	}

	cueDir := filepath.Dir(cuePath)
	cueName := filepath.Base(cuePath)[0 : len(filepath.Base(cuePath))-len(filepath.Ext(filepath.Base(cuePath)))]
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
