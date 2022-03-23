package main

// #include "xipverifier.h"
import "C"

import (
	"errors"
)

type XARIter struct {
	c C.xar_iter_t
}

func NewXARIter() (*XARIter, error) {
	c := C.xar_iter_new()
	if c == nil {
		return nil, errors.New("error creating xar iterator")
	}

	xi := XARIter{
		c: c,
	}
	return &xi, nil
}

func (xi *XARIter) C() C.xar_iter_t {
	return xi.c
}

func (xi *XARIter) IsInitialized() bool {
	return xi.c != nil
}

func (xi *XARIter) Close() error {
	if !xi.IsInitialized() {
		return errors.New("!xi.IsInitialized()")
	}

	C.xar_iter_free(xi.c)

	xi.c = nil
	return nil
}
