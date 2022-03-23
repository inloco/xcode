package main

// #include "xipverifier.h"
import "C"

type XAROpenFlags int32

const (
	XOFRead = XAROpenFlags(iota)
	XOFWrite
)

func (xof XAROpenFlags) C() C.int32_t {
	return C.int32_t(xof)
}
