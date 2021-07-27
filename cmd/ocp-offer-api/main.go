package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

// ReadFiles Принимает слайс путей и возвращает слайс считанных файлов
func ReadFiles(files []string) ([]string, error) {
	freader := func(fpath string) ([]byte, error) {
		file, err := os.OpenFile(fpath, os.O_RDONLY, 0)
		if err != nil {
			return nil, err
		}

		defer func() {
			if err = file.Close(); err != nil {
				fmt.Printf("Error closing file on path: \"%s\" \n", fpath)
				log.Fatal(err)
			} else {
				fmt.Printf("File on path \"%s\" has successfully closed \n", fpath)
			}
		}()

		data := new(bytes.Buffer)

		if _, err = data.ReadFrom(file); err != nil {
			return nil, err
		}

		return data.Bytes(), nil
	}

	result := make([]string, 0)

	for _, file := range files {
		fbyte, err := freader(file)
		if err != nil {
			return nil, err
		}
		result = append(result, string(fbyte))
	}

	return result, nil
}

func main() {

	files, _ := ReadFiles([]string{"./test/test-file-1.txt", "./test/test-file-2.txt"})
	fmt.Printf("%v\n", files)

	fmt.Println("Project: ocp-offer-api")
}
