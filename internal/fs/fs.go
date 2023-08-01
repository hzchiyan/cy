package fs

import (
	"archive/zip"
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//By https://github.com/coreybutler/go-fsutil

const (
	// KB represents the size of a kilobyte.
	KB float64 = 1024
	// MB represents the size of a megabyte.
	MB float64 = 1024 * KB
	// GB represents the size of a gigabyte.
	GB float64 = 1024 * MB
	// TB represents the size of a terabyte.
	TB float64 = 1024 * GB
	// PB represents the size of a petabyte.
	PB float64 = 1024 * TB
)

func Md5File(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return "", err
	}
	var size int64 = 1048576 // 1M
	hash := md5.New()
	if fi.Size() < size {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return "", err
		}
		hash.Write(data)
	} else {
		b := make([]byte, size)
		for {
			n, err := f.Read(b)
			if err != nil {
				break
			}
			hash.Write(b[:n])
		}
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// Sha1File sha1_file()
func Sha1File(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	hash := sha1.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// Touch Similar to the touch command on *nix, where the file
// or directory will be created if it does not already exist.
// Returns the absolute path.
// The optional second boolean argument will force
// the method to treat the path as a file instead of a directory
// (useful when the filename has no extension).
// An optional 3rd boolean argument will force the method
// to treat the path as a directory even if a file extension is present.
//
// For example:
// `fs.Touch("./path/to/archive.old", false, true)`
//
// Normally, any file path with an extension is determined
// to be a file. However; the second argument (`false`)
// instructs the command to **not** force a file. The third
// argument (`true`) instructs the command to **treat the path
// like a directory**.
func Touch(path string, flags ...any) string {
	abs := Abs(path)
	if !Exists(path) {
		forceFile := false
		forceDir := false
		if len(flags) > 0 {
			for i, flag := range flags {
				if i == 0 {
					forceFile = flag.(bool)
				} else if i == 1 {
					forceDir = flag.(bool)
				}
			}
		}
		ext := filepath.Ext(abs)
		if !forceDir && (forceFile || len(ext) > 0) {
			Mkdir(filepath.Dir(abs))
			file, err := os.Create(abs)
			if err != nil {
				panic(err)
			}
			file.Close()
		} else {
			Mkdir(abs)
		}
	}
	return abs
}

// Mkdir is the equivalent of [mkdir -p](https://en.wikipedia.org/wiki/Mkdir)
// It will generate the full directory path if it does not already
// exist.
func Mkdir(path string) string {
	path = Abs(path)
	_ = os.MkdirAll(path, os.ModePerm)
	return path
}

// Exists is a helper method to quickly
// determine whether a directory or file exists.
func Exists(path string) bool {
	if len(Abs(path)) == 0 {
		return false
	}
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// Abs Returns the fully resolved path, even if the
// path does not exist.
//
// ```
// fs.Abs("./does/not/exist")
// ```
// If the code above was run within `/home/user`, the
// result would be `/home/user/does/not/exist`.
func Abs(path string) string {
	abs, _ := filepath.Abs(path)
	return abs
}

// Clean will ensure the specified directory exists.
// If the directory already exists, all of contents
// are deleted. If the directory does not exist, it
// is automatically created.
func Clean(path string) {
	path = Abs(path)
	if IsFile(path) {
		path = filepath.Dir(path)
	}
	if Exists(path) {
		_ = os.RemoveAll(path)
	}
	Mkdir(path)
}

// IsFile determines whether the specified path
// represents a file.
func IsFile(path string) bool {
	if !Exists(path) {
		return false
	}
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !stat.IsDir()
}

// IsDirectory determines whether the specified path
// represents a directory.
func IsDirectory(path string) bool {
	if !Exists(path) {
		return false
	}
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

// WriteTextFile writes text to a file (automatically converts string to
// a byte array). If the path does not exist, it will be
// created automatically. This is the equivalent of using
// the Touch() method first, then writing text content to
// the file.
//
// It is also possible to pass a third argument, a custom permission.
// By default, os.ModePerm is used.
func WriteTextFile(path string, content string, args ...any) error {
	path = Touch(path, true)
	perm := os.ModePerm
	if len(args) > 0 {
		perm = args[0].(os.FileMode)
	}
	return ioutil.WriteFile(path, []byte(content), perm)
}

func WriteFile(path string, content string) (int, error) {
	path = Touch(path, true)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	size, err := w.Write([]byte(content))
	_ = w.Flush()
	return size, err
}

// ReadTextFile reads a text file and converts results from bytes
// to a string.
func ReadTextFile(path string) (string, error) {
	data, err := ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(Abs(path))
}

// IsReadable determines whether the file/directory is readable
// for the active system user.
func IsReadable(path string) bool {
	return allowFileAction(path, os.O_RDONLY, 0666)
}

// IsWritable determines whether the file/directory is writable
// for the active system user.
func IsWritable(path string) bool {
	return allowFileAction(path, os.O_WRONLY, 0666)
}

// IsExecutable determines whether the file/directory is executable
// for the active system user.
func IsExecutable(path string) bool {
	path = Abs(path)
	if !Exists(path) {
		return false
	}
	fileInfo, err := os.Lstat("file.txt")
	if err != nil {
		return false
	}
	mode := fileInfo.Mode()
	return mode&0111 != 0
}

func allowFileAction(path string, flag int, perm os.FileMode) bool {
	path = Abs(path)
	if !Exists(path) {
		return false
	}
	file, err := os.OpenFile(path, flag, perm)
	allowed := true
	if err != nil {
		if os.IsPermission(err) {
			allowed = false
		}
	}
	file.Close()
	return allowed
}

type listPath struct {
	Path string
	Stat os.FileInfo
}

func list(directory string, recursive bool, ignore ...string) ([]*listPath, error) {
	directory = Abs(directory)
	response := make([]*listPath, 0)
	var ignored error
	// Walk recursive lists
	if recursive {
		_ = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			ignored = isIgnoredPath(path, ignore...)
			if ignored != nil {
				return ignored
			}
			response = append(response, &listPath{
				Path: path,
				Stat: info,
			})
			return nil
		})
	} else {
		paths, matchErr := filepath.Glob(filepath.Clean(filepath.Join(directory, "/*")))
		if matchErr == nil {
			for _, path := range paths {
				ignored = isIgnoredPath(path, ignore...)
				if ignored == nil {
					stat, _ := os.Stat(path)
					response = append(response, &listPath{
						Path: path,
						Stat: stat,
					})
				}
			}
		} else {
			return make([]*listPath, 0), matchErr
		}
	}
	return response, nil
}

func isIgnoredPath(path string, ignore ...string) error {
	if len(ignore) > 0 {
		for _, pattern := range ignore {
			matched, matchErr := filepath.Match(pattern, path)
			if matchErr != nil {
				return matchErr
			}
			if matched {
				return errors.New("Ignored")
			}
		}
	}
	return nil
}

// List Generate a list of path names for the given directory.
// Optionally provide a list of ignored paths, using
// [glob](https://en.wikipedia.org/wiki/Glob_%28programming%29) syntax.
func List(directory string, recursive bool, ignore ...string) ([]string, error) {
	response, err := list(directory, recursive, ignore...)
	if err != nil {
		return make([]string, 0), err
	}
	paths := make([]string, len(response))
	for i := range response {
		paths[i] = response[i].Path
	}
	return paths, nil
}

// ListDirectories provides absolute paths of directories only, ignoring files.
func ListDirectories(directory string, recursive bool, ignore ...string) ([]string, error) {
	paths := make([]string, 0)
	response, err := list(directory, recursive, ignore...)
	if err != nil {
		return paths, err
	}
	if len(response) == 0 {
		return paths, nil
	}
	for _, item := range response {
		if item.Stat.IsDir() {
			paths = append(paths, item.Path)
		}
	}
	return paths, nil
}

// ListFiles provides absolute paths of files only, ignoring directories.
func ListFiles(directory string, recursive bool, ignore ...string) ([]string, error) {
	paths := make([]string, 0)
	response, err := list(directory, recursive, ignore...)
	if err != nil {
		return paths, err
	}
	if len(response) == 0 {
		return paths, nil
	}
	for _, item := range response {
		if !item.Stat.IsDir() {
			paths = append(paths, item.Path)
		}
	}
	return paths, nil
}

// ByteSize returns the number of bytes (size) of a file/directory.
func ByteSize(path string) (int64, error) {
	path = Abs(path)
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	if err != nil {
		return -1, err
	}
	return size, nil
}

// Size returns a "pretty" version of the size, such as "3.12MB"
func Size(path string, sigfig ...int) (string, error) {
	size, err := ByteSize(path)
	if err != nil {
		return "", err
	}
	return FormatSize(size, sigfig...), nil
}

// FormatSize returns a nicely formatted representation of a number of bytes,
// such as `3.14MB`
func FormatSize(bytesize int64, sigfig ...int) string {
	size := float64(bytesize)
	var sigfigs int
	if len(sigfig) == 0 {
		sigfigs = 2
	} else {
		sigfigs = sigfig[0]
	}
	switch {
	case size >= PB:
		return strconv.FormatFloat(math.Round((size*100)/PB)/100, 'f', sigfigs, 64) + "PB"
	case size >= TB:
		return strconv.FormatFloat(math.Round((size*100)/TB)/100, 'f', sigfigs, 64) + "TB"
	case size >= GB:
		return strconv.FormatFloat(math.Round((size*100)/GB)/100, 'f', sigfigs, 64) + "GB"
	case size >= MB:
		return strconv.FormatFloat(math.Round((size*100)/MB)/100, 'f', sigfigs, 64) + "MB"
	case size >= KB:
		return strconv.FormatFloat(math.Round((size*100)/KB)/100, 'f', sigfigs, 64) + "KB"
	default:
		return strconv.FormatInt(bytesize, 10) + "B"
	}
}

// Symlink creates a symbolic link. This just runs `os.Symlink()`.
func Symlink(target string, name string) error {
	return os.Symlink(target, name)
}

// IsSymlink determines whether the path is a symbolic link.
func IsSymlink(path string) bool {
	info, err := os.Readlink(path)
	return err == nil && len(info) > 0
}

// LastModified identies the last time the path was modified.
func LastModified(path string) (time.Time, error) {
	file, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}
	return file.ModTime(), nil
}

// Move a file/directory to another location
func Move(source string, dest string, ignoreErrors ...bool) error {
	ignore := false
	if len(ignoreErrors) > 0 {
		ignore = ignoreErrors[0]
	}
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		stub := strings.Replace(path, source, "", 1)
		target := filepath.Join(dest, stub)

		if info.IsDir() {
			Touch(target)
		} else if !IsSymlink(path) {
			err := os.Rename(path, target)
			if err != nil && !ignore {
				return err
			}
		}
		return nil
	})
}

// Copy a file/directory
func Copy(source string, dest string, ignoreErrors ...bool) error {
	ignore := false
	if len(ignoreErrors) > 0 {
		ignore = ignoreErrors[0]
	}
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		stub := strings.Replace(path, source, "", 1)
		target := filepath.Join(dest, stub)
		if info.IsDir() {
			Touch(target)
		} else if !IsSymlink(path) {
			input, err := ioutil.ReadFile(path)
			if err != nil && !ignore {
				return err
			}
			err = ioutil.WriteFile(target, input, 0644)
			if err != nil && !ignore {
				return err
			}
		}
		return nil
	})
}

// Unzip a file
func Unzip(src string, dest string) error {
	src = Abs(src)
	if !Exists(src) {
		return errors.New(src + " does not exist")
	}
	dest = Abs(dest)
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()
	_ = os.MkdirAll(dest, 0775)
	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()
		path := filepath.Join(dest, f.Name)
		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}
		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(path, f.Mode())
		} else {
			_ = os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()
			if _, err = io.Copy(f, rc); err != nil {
				return err
			}
		}
		return nil
	}
	for _, f := range r.File {
		if err := extractAndWriteFile(f); err != nil {
			return err
		}
	}
	return nil
}

// Zip a file or directory. Does not follow symlinks.
func Zip(src string, target ...string) error {
	dest := strings.Replace(filepath.Base(src), filepath.Ext(src), "", 1) + ".zip"
	if len(target) > 0 {
		dest = target[0]
	}
	dest = Abs(dest)
	newZipFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer newZipFile.Close()
	writer := zip.NewWriter(newZipFile)
	defer writer.Close()
	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && !IsSymlink(path) {
			input, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			localPath := strings.Replace(path, src, "", 1)
			if strings.HasPrefix(localPath, "/") {
				localPath = strings.Replace(localPath, "/", "", 1)
			} else if strings.HasPrefix(localPath, "\\") {
				localPath = strings.Replace(localPath, "\\", "", 1)
			}
			if err = addToZipArchive(writer, localPath, input); err != nil {
				return err
			}
		}
		return nil
	})
}

func addToZipArchive(archive *zip.Writer, filename string, content []byte) error {
	f, err := archive.Create(filename)
	if err != nil {
		return err
	}
	_, err = f.Write(content)
	if err != nil {
		return err
	}
	return nil
}
