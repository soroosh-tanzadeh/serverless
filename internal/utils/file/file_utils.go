package file

import (
	"os"
	"path"
	"sync"
	"time"
	"unicode/utf8"
)

func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}

type RWValue struct {
	mutex     sync.RWMutex
	value     interface{}
	timestamp time.Time // time of last set()
}

func (v *RWValue) Set(value interface{}) {
	v.mutex.Lock()
	v.value = value
	v.timestamp = time.Now()
	v.mutex.Unlock()
}

func (v *RWValue) Get() (interface{}, time.Time) {
	v.mutex.RLock()
	defer v.mutex.RUnlock()
	return v.value, v.timestamp
}

func IsText(s []byte) bool {
	const max = 1024
	if len(s) > max {
		s = s[0:max]
	}
	for i, c := range string(s) {
		if i+utf8.UTFMax > len(s) {
			break
		}
		if c == 0xFFFD || c < ' ' && c != '\n' && c != '\t' && c != '\f' {
			return false
		}
	}
	return true
}

var textExt = map[string]bool{
	".css": false,
	".svg": false,
}

func IsTextFile(filename string) bool {
	if isText, found := textExt[path.Ext(filename)]; found {
		return isText
	}

	f, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer f.Close()

	var buf [1024]byte
	n, err := f.Read(buf[0:])
	if err != nil {
		return false
	}

	return IsText(buf[0:n])
}
