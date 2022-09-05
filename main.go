package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

type application struct {
	searchedHashes []string
}

func main() {
	start := time.Now()
	var projName string
	var hashesFilename string
	flag.StringVar(&projName, "name", "hashInfo", "name of the hashing project")
	flag.StringVar(&hashesFilename, "hashes", "", "name of the file containing hashes to search")
	flag.Parse()
	var topDirectory string
	if topDirectory = flag.Arg(0); topDirectory == "" {
		if os := runtime.GOOS; os == "windows" {
			topDirectory = `c:\`
		} else {
			topDirectory = "/"
		}
	}
	if hashesFilename != "" {
		f, err := os.Open(hashesFilename)
		checkError("Can't open file", err)
		defer f.Close()
		reader := csv.NewReader(f)
		reader.ReuseRecord = true
		var allRecords []string
		for {
			record, err := reader.Read()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				} else {
					checkError("error reading file", err)
				}
			}
			allRecords = append(allRecords, record...)
		}
	}
	resultName := strings.ReplaceAll(projName+"["+time.Now().Format(time.Stamp)+"].", ":", "")
	var allHashes sort.StringSlice
	allHashes, err := hashAll(topDirectory, resultName)
	checkError("Can't hash files!", err)
	allHashes.Sort()
	resultFile, err := os.Create(resultName + "csv")
	checkError("Cannot create file", err)

	writer := bufio.NewWriter(resultFile)
	defer writer.Flush()

	for _, row := range allHashes {
		_, err = writer.WriteString(row + "\n")
		checkError("can't write to logfile", err)
	}

	fmt.Printf("Hashing took: %v\n", time.Since(start))
}

// hashAll returns slice of files under directory in parameter and its corresponding hashes.
func hashAll(root string, resultName string) ([]string, error) {
	var result []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, errWalk error) error {
		if errWalk != nil {
			f, err := os.OpenFile(resultName+"err", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				log.Fatalf("error opening file: %v", err)
			}

			defer f.Close()

			log.SetOutput(f)
			log.Println("Can't open file:", errWalk)
			return nil
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		file, err := os.OpenFile(path, os.O_RDONLY, 0)
		if err != nil {
			log.Println("Can't open file", err)
			return nil
		}

		defer file.Close()

		result = append(result, calculateBasicHashes(file, path))

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// checkError surprisingly checks for error.
func checkError(message string, err error) {
	if err != nil {
		log.Println(message, err)
	}
}

// calculateBasicHashes returns struct with all hashes of the Reader in parameter.
func calculateBasicHashes(rd io.Reader, path string) string {
	result := "[" + path + "], "
	md5hash := md5.New()
	sha1hash := sha1.New()
	sha256hash := sha256.New()

	// For optimum speed, pagesize contains the underlying system's memory page size.
	pagesize := os.Getpagesize()

	// wraps the Reader object into a new buffered reader to read the files in chunks
	// and buffering them for performance.
	reader := bufio.NewReaderSize(rd, pagesize)

	// creates a multiplexer Writer object that will duplicate all write
	// operations when copying data from source into all different hashing algorithms
	// at the same time
	multiWriter := io.MultiWriter(md5hash, sha1hash, sha256hash)

	// Using a buffered reader, this will write to the writer multiplexer
	// so we only traverse through the file once, and can calculate all hashInfo
	// in a single byte buffered scan pass.
	_, err := io.Copy(multiWriter, reader)
	checkError("Can't copy to multiwriter", err)

	md5sum := hex.EncodeToString(md5hash.Sum(nil))
	sha1sum := hex.EncodeToString(sha1hash.Sum(nil))
	sha256sum := hex.EncodeToString(sha256hash.Sum(nil))
	result += "[" + md5sum + "], "
	result += "[" + sha1sum + "], "
	result += "[" + sha256sum + "]"

	return result
}
