package state

import (
	"os"
	"log"
    "github.com/peterbourgon/diskv"
)

var disk *diskv.Diskv

func Load(BasePath string) {

	if _, err := os.Stat(BasePath); os.IsNotExist(err) {

	    if err = os.MkdirAll(BasePath, 0700); err != nil {
	    	log.Fatal(err)
	    }
	}

	flatTransform := func(s string) []string { return []string{} }

	disk = diskv.New(diskv.Options{
		BasePath:     BasePath,
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024,
	})
}

func ReadBool(key string) (bool, error) {
	
	value, err := disk.Read(key)

	if err != nil {
		return false, err
	}

	return value[0] == 1, err
}

func WriteBool(key string, value bool) {

	if value {
		disk.Write(key, []byte{1})
	} else {
		disk.Write(key, []byte{0})
	}
}

func ReadString(key string) (string, error) {

	value, err := disk.Read(key)

	if err != nil {
		return "", err
	}

	return string(value[:len(value)]), err
}

func WriteString(key string, value string) {
	disk.Write(key, []byte(value))
}
