package main

import (
	"strings"
	"crypto/md5"
	"fmt"
)

type Subs struct {
	subs []string
}

func (s *Subs) Append(sub interface{}) *Subs {
	switch sub.(type) {
	case string:
		s.subs = append(s.subs, sub.(string))
	case []string:
		tmpSlice := sub.([]string)
		for i := 0; i < len(tmpSlice); i++ {
			s.subs = append(s.subs, tmpSlice[i])
		}
	}
	return s
}

func (s *Subs) EntrySet() []string {
	return s.subs
}

func (s *Subs) Hash() string {
	md5Buf := md5.Sum([]byte(strings.Join(s.EntrySet(), ",")))
	return fmt.Sprintf("%x", md5Buf[:])
}
