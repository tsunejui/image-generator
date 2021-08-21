package pkg

import (
	"fmt"
	"os"
)

func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

func ExportFile(path string, data string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	if _, err := file.WriteString(data); err != nil {
		return fmt.Errorf("failed to write data: %v", err)
	}
	return nil
}
