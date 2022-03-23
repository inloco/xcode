package main

// #include "xipverifier.h"
import "C"

import (
	"bytes"
	"errors"
	"unsafe"
)

type XAR struct {
	c C.xar_t
}

func NewXAR(file string, flags XAROpenFlags) (*XAR, error) {
	cFile := C.CString(file)
	defer C.free(unsafe.Pointer(cFile))

	c := C.xar_open(cFile, flags.C())
	if c == nil {
		return nil, errors.New("error opening xar file")
	}

	xar := XAR{
		c: c,
	}
	return &xar, nil
}

func (x *XAR) C() C.xar_t {
	return x.c
}

func (x *XAR) IsInitialized() bool {
	return x.c != nil
}

func (x *XAR) RegisterErrHandler() error {
	if !x.IsInitialized() {
		return errors.New("!x.IsInitialized()")
	}

	return RegisterErrHandler(x)
}

func (x *XAR) GetTOCChecksum() []byte {
	var checksumSize C.size_t
	checksumData := C.xar_get_toc_checksum(x.c, &checksumSize)
	defer C.free(unsafe.Pointer(checksumData))

	checksum := C.GoBytes(unsafe.Pointer(checksumData), C.int(checksumSize))

	return checksum
}

func (x *XAR) GetFiles() ([]XARFile, error) {
	if !x.IsInitialized() {
		return nil, errors.New("!x.IsInitialized()")
	}

	var xfs []XARFile

	iter, err := NewXARIter()
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	c := C.xar_file_first(x.c, iter.C())
	for c != nil {
		xfs = append(xfs, XARFile{
			c: c,
		})

		c = C.xar_file_next(iter.C())
	}

	return xfs, nil
}

func (x *XAR) GetSignatures() ([]XARSignature, error) {
	if !x.IsInitialized() {
		return nil, errors.New("!x.IsInitialized()")
	}

	var xss []XARSignature

	c := C.xar_signature_first(x.c)
	for c != nil {
		xss = append(xss, XARSignature{
			c: c,
		})

		c = C.xar_signature_next(c)
	}

	return xss, nil
}

func (x *XAR) Verify() error {
	if !x.IsInitialized() {
		return errors.New("!x.IsInitialized()")
	}

	if err := x.verifySignatures(); err != nil {
		return err
	}

	if err := x.verifyFiles(); err != nil {
		return err
	}

	return nil
}

func (x *XAR) verifySignatures() error {
	xss, err := x.GetSignatures()
	if err != nil {
		return err
	}

	tocChecksum := x.GetTOCChecksum()
	for _, xs := range xss {
		sigChecksum, err := xs.Verify()
		if err != nil {
			return err
		}

		if !bytes.Equal(sigChecksum, tocChecksum) {
			return errors.New("!bytes.Equal(sigChecksum, tocChecksum)")
		}
	}

	return nil
}

func (x *XAR) verifyFiles() error {
	xfs, err := x.GetFiles()
	if err != nil {
		return err
	}

	for _, xf := range xfs {
		if err := xf.Verify(x); err != nil {
			return err
		}
	}

	return nil
}

func (x *XAR) Close() error {
	if !x.IsInitialized() {
		return errors.New("!x.IsInitialized()")
	}

	if err := C.xar_close(x.c); err != 0 {
		return GetError()
	}

	x.c = nil
	return nil
}
