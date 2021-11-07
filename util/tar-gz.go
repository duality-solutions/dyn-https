package util

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// ExtractTarGz extracts the tar.gz file passed in path parameter
func ExtractTarGz(src string, dest string, keep []string) ([]string, error) {
	var filenames []string
	//var directories []string
	gzipStream, err := os.Open(src)
	if err != nil {
		Error.Println("ExtractTarGz gzipStream Open err", err)
		return filenames, err
	}
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		Error.Println("ExtractTarGz NewReader err", err)
		return filenames, err
	}

	tarReader := tar.NewReader(uncompressedStream)
	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
		switch header.Typeflag {
		case tar.TypeReg:
			if contains(keep, header.Name) {
				// if in keep array
				file := dest + filepath.Base(header.Name)
				Info.Println("ExtractTarGz file", file)
				outFile, err := os.Create(file)
				if err != nil {
					Error.Println("ExtractTarGz Create error", err)
					continue
				}
				if _, err := io.Copy(outFile, tarReader); err != nil {
					Error.Println("ExtractTarGz Copy error", err)
					continue
				}
				outFile.Close()
				filenames = append(filenames, header.Name)
			}

		default:
			continue
		}
	}
	// Close the file without defer to close before next iteration of loop
	uncompressedStream.Close()
	gzipStream.Close()
	return filenames, nil
}
