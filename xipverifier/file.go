package main

// #include "xipverifier.h"
import "C"

import (
	"errors"
)

type XARFile struct {
	c C.xar_file_t
}

func (xf *XARFile) C() C.xar_file_t {
	return xf.c
}

func (xf *XARFile) IsInitialized() bool {
	return xf.c != nil
}

func (xf *XARFile) Verify(xar *XAR) error {
	if !xf.IsInitialized() {
		return errors.New("!xi.IsInitialized()")
	}

	if xar == nil {
		return errors.New("xar == nil")
	}

	if !xar.IsInitialized() {
		return errors.New("!xar.IsInitialized()")
	}

	if err := C.xar_verify(xar.C(), xf.c); err != 0 {
		return GetError()
	}

	return nil
}
