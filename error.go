package main

import (
	"errors"
)

var ErrXdrParsFirstFloorErr = errors.New("xdrParsFirstFloorErr")
var ErrXdrParsScdFloorErr = errors.New("xdrParsScdFloorErr")
var ErrXdrParseErr = errors.New("xdrParseErr")
var ErrSuffixErr = errors.New("suffixErr")
var ErrReadFileErr = errors.New("readFileErr")
var ErrPushKafkaErr = errors.New("pushKafkaErr")
var ErrNotConfErr = errors.New("NotConfErr")

func isXdrPkgErr(err error) bool {
	var errPkgMap = map[error]int{
		ErrXdrParsFirstFloorErr: 0,
		ErrXdrParsScdFloorErr:   0,
		ErrXdrParseErr:          0,
	}
	if _, ok := errPkgMap[err]; ok {
		return true
	}
	return false
}
