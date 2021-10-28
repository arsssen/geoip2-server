package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func getContent(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("got status %d", resp.StatusCode)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(bodyBytes), nil

}

func downloadFile(filepath string, url string) (err error) {
	log.Printf("downloading from %s", url)
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func extractTarGz(path, filename string) (string, error) {
	gzipStream, err := os.Open(path)
	if err != nil {
		fmt.Println("error")
	}
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		log.Fatal("ExtractTarGz: NewReader failed")
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return "", fmt.Errorf("ExtractTarGz: Next() failed: %w", err)
		}

		switch header.Typeflag {
		case tar.TypeReg:
			if filepath.Base(header.Name) != filename {
				log.Printf("skipping file: %s", header.Name)
				if _, err := io.Copy(io.Discard, tarReader); err != nil {
					return "", fmt.Errorf("ExtractTarGz: dCopy() failed: %w", err)
				}
				continue
			}
			targetName := header.ChangeTime.Format("02_01_2006") + ".mmdb"
			if _, err := os.Stat(targetName); err == nil {
				log.Printf("%s exists, skipping extraction", targetName)
				return targetName, nil
			}
			outFile, err := os.Create(targetName)

			if err != nil {
				return "", fmt.Errorf("ExtractTarGz: Create() failed: %w", err)
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return "", fmt.Errorf("ExtractTarGz: Copy() failed: %w", err)
			}
			outFile.Close()
			return targetName, nil

		default:
			log.Printf("found %s in archive", header.Name)
		}

	}
	return "", nil
}
