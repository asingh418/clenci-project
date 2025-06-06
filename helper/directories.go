package helper

import (
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

// MkDirsIfNotExist creates a directory named path,
// along with any necessary parents, and returns true, or else returns false.
// The dirs have permission bits 511 (before umask) are used for all directories that this function creates.
// If path is already a directory, the function does nothing and returns false.
func MkDirsIfNotExist(name string) bool {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		logrus.Infof("creating directory %s", name)
		err = os.MkdirAll(name, os.ModePerm)
		if err != nil {
			logrus.Errorf("unable to create %s directory", name)
			return false
		}
		logrus.Infof("directory %s created", name)
		return true
	}

	return false
}

// CreateDirectoryNamedPath creates a directory named path, along with any necessary parents, and returns nil, or else returns an error.
// The permission bits perm (before umask) are used for all directories. If path is already a directory does nothing and returns nil.
func CreateDirectoryNamedPath(path string) (string, error) {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		logrus.Fatal("Unable to create directory (and its parents)", err)
	}

	return path, err
}

// CreateTempDir creates a new temporary directory in the directory dir.
// The directory name is generated by taking pattern and applying a random string to the end
func CreateTempDir(dir string, pattern string) string {
	dir, err := ioutil.TempDir(dir, pattern)
	if err != nil {
		logrus.Fatal(err)
	}

	return dir
}
