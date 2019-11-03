package util

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func UploadFileAnDeleteOld(r *http.Request, folder, pattern, oldFile string) (err error, path string) {
	// Parse our multipart form, 10 << 20 specifies a maximum upload of 10 MB files.
	err = r.ParseMultipartForm(10 << 20)

	if err != nil {
		return err, path
	}

	// FormFile returns the first file for the given key `file`
	file, handler, err := r.FormFile("file")
	if err != nil {
		return err, path
	}
	defer file.Close()
	log.Printf("Uploaded File: %v\n", handler.Filename)
	// log.Printf("File Size: %+v\n", handler.Size)
	// log.Printf("MIME Header: %+v\n", handler.Header)

	// Create a file
	tempFile, err := ioutil.TempFile(folder, pattern)
	if err != nil {
		return err, path
	}

	defer tempFile.Close()

	// read all of the contents of our uploaded file into a byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err, path
	}

	// write this byte array to our temporary file
	tempFile.Write(fileBytes)

	// remove old avatar
	os.Remove(oldFile)

	return err, tempFile.Name()
}
