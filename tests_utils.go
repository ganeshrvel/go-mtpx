package mtpx

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func getTestMocksAsset(_filePath string) string {
	currentDir, err := os.Getwd()

	if err != nil {
		log.Panicf("\nunable to fetch the current directory: %s\n", currentDir)
	}

	resultPath := fmt.Sprintf("%s/tests/mtp-test-files/", currentDir)

	resultPath = fmt.Sprintf("%s%s", resultPath, _filePath)

	if exist := existsLocal(resultPath); !exist {
		log.Panicf("\nthe 'mocks' asset not found: %s\n", resultPath)
	}

	return resultPath
}

func newTestMocksAsset(_filePath string) string {
	currentDir, err := os.Getwd()

	if err != nil {
		log.Panicf("\nunable to fetch the current directory: %s\n", currentDir)
	}

	resultPath := fmt.Sprintf("%s/tests/mtp-test-files/", currentDir)

	resultPath = fmt.Sprintf("%s%s", resultPath, _filePath)

	return resultPath
}

func newTempMocksAsset(_filePath string) string {
	currentDir, err := os.Getwd()

	if err != nil {
		log.Panicf("\nunable to fetch the current directory: %s\nerror: %+v\n", currentDir, err.Error())
	}

	resultPath := fmt.Sprintf("%s/tests/mocks-build/", currentDir)

	if exist := isDirLocal(resultPath); !exist {
		_, err := os.Create(resultPath)

		if err != nil {
			log.Panicf("\n'mocks-build' directory not found: %s\nerror: %+v\n", resultPath, err.Error())
		}
	}

	resultPath = fmt.Sprintf("%s%s", resultPath, _filePath)

	return resultPath
}

func newTempMocksDir(_dirPath string, resetDir bool) string {
	currentDir, err := os.Getwd()

	if err != nil {
		log.Panicf("\nunable to fetch the current directory: %s\nerror: %+v\n", currentDir, err.Error())
	}

	resultPath := filepath.Join(currentDir, "tests/mocks-build", _dirPath)

	if resetDir {
		err := os.RemoveAll(resultPath)

		if err != nil {
			log.Panic(err)
		}

		if exist := isDirLocal(resultPath); !exist {
			err = os.MkdirAll(resultPath, os.ModePerm)

			if err != nil {
				log.Panicf("\ntemp mocks directory not found: %s\nerror: %+v\n", resultPath, err.Error())
			}
		}
	}

	if exist := isDirLocal(resultPath); !exist {
		err := os.MkdirAll(resultPath, os.ModePerm)

		if err != nil {
			log.Panicf("\ntemp mocks directory not found: %s\nerror: %+v\n", resultPath, err.Error())
		}
	}

	return resultPath
}
