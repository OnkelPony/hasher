package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type hashInfo struct {
	md5Sum    string
	sha1Sum   string
	sha256Sum string
}

var logName = strings.ReplaceAll("hashInfo["+time.Now().Format(time.Stamp)+"].", ":", "")

func main() {
	start := time.Now()
	root := os.Args[1]

	allHashes, err := hashAll(root)
	checkError("Can't hash files!", err)

	resFile, err := os.Create(logName + "csv")
	checkError("Cannot create file", err)

	writer := bufio.NewWriter(resFile)
	defer writer.Flush()

	for file, threeHashes := range allHashes {
		row := fmt.Sprintf("[%v], [%s], [%s], [%s]\n", file, threeHashes.md5Sum, threeHashes.sha1Sum, threeHashes.sha256Sum)
		_, err = writer.WriteString(row)
		checkError("can't write to logfile", err)
	}

	fmt.Printf("Hashing took: %v\n", time.Since(start))
}

// hashAll returns map of files under root directory in parameter and its corresponding hashes.
func hashAll(root string) (map[string]hashInfo, error) {
	m := make(map[string]hashInfo)

	err := filepath.Walk(root, func(path string, info os.FileInfo, errWalk error) error {
		if errWalk != nil {
			f, err := os.OpenFile(logName+"err", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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

		m[path] = calculateBasicHashes(file)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return m, nil
}

// checkError surprisingly checks for error.
func checkError(message string, err error) {
	if err != nil {
		log.Println(message, err)
	}
}

// calculateBasicHashes returns struct with all hashes of the Reader in parameter.
func calculateBasicHashes(rd io.Reader) hashInfo {

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
	//
	_, err := io.Copy(multiWriter, reader)
	checkError("Can't copy to multiwriter", err)

	var info hashInfo

	info.md5Sum = hex.EncodeToString(md5hash.Sum(nil))
	info.sha1Sum = hex.EncodeToString(sha1hash.Sum(nil))
	info.sha256Sum = hex.EncodeToString(sha256hash.Sum(nil))

	return info
}
