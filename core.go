package main

import (
	"strings"
	"crypto/md5"
	"fmt"
)

type Subs struct {
	subs []string
}

func (s *Subs) Append(sub string) *Subs {
	s.subs = append(s.subs, sub)
	return s
}

func (s *Subs) EntrySet() []string {
	return s.subs
}

func (s *Subs) Hash() string {
	md5Buf := md5.Sum([]byte(strings.Join(s.EntrySet(), ",")))
	return fmt.Sprintf("%x", md5Buf[:])
}
