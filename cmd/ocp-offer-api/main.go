package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

// Задание: Используя defer и функтор реализовать открытие и закрытие файла в цикле
//
// Принемает слайс путей и возвращает слайс считаных файлов
func FilesReader(files []string) []string {
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
			log.Fatal(err)
		}
		result = append(result, string(fbyte))
	}

	return result
}

func main() {
	fpaths := []string{"./main.go", "../../internal/model/offer.go"}
	files := FilesReader(fpaths)
	fmt.Println(files)

	fmt.Println("Project: ocp-offer-api")

}
