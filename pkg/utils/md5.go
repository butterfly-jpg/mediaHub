package utils

import "crypto/md5"

func MD5(content []byte) []byte {
	m := md5.New()
	m.Write(content)
	bs := m.Sum(nil)
	return bs
}
