package files

import (
	"crypto/md5"
	"encoding/hex"
	log "github.com/sirupsen/logrus"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func FileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func EnsurePathDir(path string, mode string) {
	pathDir := filepath.Dir(path)
	m, _ := strconv.ParseUint(mode, 8, 32)
	if _, err := os.Stat(pathDir); os.IsNotExist(err) {
		os.MkdirAll(pathDir, fs.FileMode(m))
	}
	os.Chmod(pathDir, fs.FileMode(m))
}

func EnsureDir(path string, mode string) {
	EnsurePathDir(path+"/a", mode)
}

func GetFileSize(path string) int64 {
	if stat, err := os.Stat(path); err == nil {
		return stat.Size()
	}
	return 0
}

func WaitExists(path string) {
	for {
		if _, err := os.Stat(path); err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func DeleteFile(path string) {
	if err := os.Remove(path); err != nil {
		log.Fatal(err)
	}
}

func DeleteDir(path string) {
	if err := os.RemoveAll(path); err != nil {
		log.Fatal(err)
	}
}

func DeleteIfExists(path string) {
	if FileExists(path) {
		if stat, _ := os.Stat(path); stat.IsDir() {
			DeleteDir(path)
		} else {
			DeleteFile(path)
		}
	}
}

func FileMd5(path string) string {
	file, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return ""
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func Chmod(path string, mode string, defaultMode string) error {
	m, _ := strconv.ParseUint(defaultMode, 8, 32)
	if mode != "" {
		m, _ = strconv.ParseUint(mode, 8, 32)
	}
	// change sock file permission
	if err := os.Chmod(path, fs.FileMode(m)); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func ListFiles(path string) []FileItem {
	var files []FileItem
	filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		files = append(files, FileItem{
			IsDir: info.IsDir(),
			Name:  info.Name(),
			Path:  path,
			Size:  info.Size(),
			Mtime: info.ModTime().Unix(),
		})
		return nil
	})
	return files
}
