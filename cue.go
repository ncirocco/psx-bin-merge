package binmerge

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type cueBinFile struct {
	File   string
	Tracks []track
}

type track struct {
	ID        int
	TrackType string
	Indexes   []index
}

type index struct {
	ID     int
	Offset int
}

func cueToCueMap(cuePath string) ([]cueBinFile, error) {
	f, err := os.Open(cuePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var binFiles []cueBinFile
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		switch fields[0] {
		case "FILE":
			var binFile cueBinFile
			binFile.File = strings.Replace(strings.Join(fields[1:len(fields)-1], " "), "\"", "", -1)
			binFiles = append(binFiles, binFile)
		case "TRACK":
			var track track
			track.ID, err = strconv.Atoi(fields[1])
			if err != nil {
				return nil, err
			}
			track.TrackType = fields[2]
			lastBinFile := &binFiles[len(binFiles)-1]
			lastBinFile.Tracks = append(lastBinFile.Tracks, track)
		case "INDEX":
			var index index
			index.ID, err = strconv.Atoi(fields[1])
			if err != nil {
				return nil, err
			}

			index.Offset, err = stampToSectors(fields[2])
			if err != nil {
				return nil, err
			}

			lastTrack := &binFiles[len(binFiles)-1].Tracks[len(binFiles[len(binFiles)-1].Tracks)-1]
			lastTrack.Indexes = append(
				lastTrack.Indexes,
				index,
			)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return binFiles, nil
}

func createCuesheet(name string, basePath string, cueMap []cueBinFile, destination string) error {
	blockSize, err := getBlockSize(cueMap[0].Tracks[0].TrackType)
	if err != nil {
		return err
	}

	var sectorPosition int64
	cue := fmt.Sprintf("FILE \"%s.bin\" BINARY\n", name)
	for _, cueBinFile := range cueMap {
		for _, track := range cueBinFile.Tracks {
			cue += fmt.Sprintf("  TRACK %02d %s\n", track.ID, track.TrackType)
			for _, index := range track.Indexes {
				cue += fmt.Sprintf("    INDEX %02d %s\n", index.ID, sectorsToStamp(sectorPosition+int64(index.Offset)))
			}
		}
		fi, err := os.Stat(filepath.Join(basePath, cueBinFile.File))
		if err != nil {
			return err
		}

		sectorPosition += fi.Size() / int64(blockSize)
	}

	_ = os.MkdirAll(destination, os.ModePerm)

	path := filepath.Join(destination, name+".cue")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString(cue)

	return nil
}
