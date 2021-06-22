package bs

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"
)

func Getwd() string {
	dir, err := os.Getwd()
	if err != nil {
		fnErrorHandler(err)
	}
	return dir
}

func Chdir(dir string) {
	Verbosef("Chdir: %s", dir)
	if err := os.Chdir(dir); err != nil {
		fnErrorHandler(err)
	}
}

func MkdirAll(dir string) {
	Verbosef("MkdirAll: %s", dir)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fnErrorHandler(err)
	}
}

// TODO: if file doesn't exist, don't err
func Remove(dir string) {
	Verbosef("Remove: %s", dir)
	if err := os.Remove(dir); err != nil {
		fnErrorHandler(err)
	}
}

func RemoveAll(dir string) {
	Verbosef("RemoveAll: %s", dir)
	if err := os.RemoveAll(dir); err != nil {
		fnErrorHandler(err)
	}
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			fnErrorHandler(err)
		}
		return false
	}
	return true
}

func IsFile(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			fnErrorHandler(err)
		}
		return false
	}
	return !fi.IsDir()
}

func IsDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			fnErrorHandler(err)
		}
		return false
	}
	return fi.IsDir()
}

func Stat(path string) fs.FileInfo {
	Verbosef("Stat: %s", path)
	fi, err := os.Stat(path)
	if err != nil {
		fnErrorHandler(err)
	}
	return fi
}

// Write file (create or truncate)

func Write(path string, contents string) {
	if err := writeImpl(path, contents, nil, false); err != nil {
		fnErrorHandler(err)
	}
}

func Writef(path string, format string, args ...interface{}) {
	if err := writeImpl(path, fmt.Sprintf(format, args...), nil, false); err != nil {
		fnErrorHandler(err)
	}
}

func WriteErr(path string, contents string) error {
	return writeImpl(path, contents, nil, false)
}

func WriteBytes(path string, b []byte) {
	if err := writeImpl(path, "", b, false); err != nil {
		fnErrorHandler(err)
	}
}

func WriteBytesErr(path string, b []byte) error {
	return writeImpl(path, "", b, false)
}

// Append file

func Append(path string, contents string) {
	if err := writeImpl(path, contents, nil, true); err != nil {
		fnErrorHandler(err)
	}
}

func Appendf(path string, format string, args ...interface{}) {
	if err := writeImpl(path, fmt.Sprintf(format, args...), nil, true); err != nil {
		fnErrorHandler(err)
	}
}

func AppendErr(path string, contents string) error {
	return writeImpl(path, contents, nil, true)
}

func AppendBytes(path string, b []byte) {
	if err := writeImpl(path, "", b, true); err != nil {
		fnErrorHandler(err)
	}
}

func AppendBytesErr(path string, b []byte) error {
	return writeImpl(path, "", b, true)
}

func writeImpl(path string, str string, b []byte, append bool) error {
	if len(str) > 0 && len(b) > 0 {
		return fmt.Errorf("this should never happen: writeImpl has both string and []byte")
	}
	var f *os.File
	var err error
	if append {
		Verbosef("Append to file: %s", path)
		f, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	} else {
		Verbosef("Write to file: %s", path)
		f, err = os.Create(path)
	}
	if err != nil {
		return err
	}
	defer f.Close()
	if len(str) > 0 {
		_, err = io.Copy(f, strings.NewReader(str))
	} else {
		_, err = io.Copy(f, bytes.NewReader(b))
	}
	if err != nil {
		return err
	}
	return nil
}

// Read file

func Read(path string) string {
	str, err := ReadErr(path)
	if err != nil {
		fnErrorHandler(err)
	}
	return str
}

func ReadErr(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	// TODO: this cast from []byte to string involves an allocation and copy.
	// Is there a way to skip that work and read straight into a string?
	return string(b), nil
}

func ReadFile(path string) []byte {
	Verbosef("Read from file: %s", path)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		fnErrorHandler(err)
	}
	return b
}
