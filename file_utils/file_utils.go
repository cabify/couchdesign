package file_utils

import "os"

func MissingDir(dir string) bool {
	_, err := os.Stat(dir)
	return err != nil
}
