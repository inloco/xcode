package main

// #include "xipverifier.h"
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

var (
	errCh chan error
)

func init() {
	errCh = make(chan error, 1)
}

//export ErrHandler
func ErrHandler(sev C.int32_t, err C.int32_t, ctx C.xar_errctx_t, _ unsafe.Pointer) C.int32_t {
	text := "libxar"

	switch sev {
	case 1:
		text = fmt.Sprintf("%s [XAR_SEVERITY_DEBUG]", text)
	case 2:
		text = fmt.Sprintf("%s [XAR_SEVERITY_INFO]", text)
	case 3:
		text = fmt.Sprintf("%s [XAR_SEVERITY_NORMAL]", text)
	case 4:
		text = fmt.Sprintf("%s [XAR_SEVERITY_WARNING]", text)
	case 5:
		text = fmt.Sprintf("%s [XAR_SEVERITY_NONFATAL]", text)
	case 6:
		text = fmt.Sprintf("%s [XAR_SEVERITY_FATAL]", text)
	default:
		text = fmt.Sprintf("%s [%d]", text, int32(sev))
	}

	switch err {
	case 1:
		text = fmt.Sprintf("%s [XAR_ERR_ARCHIVE_CREATION]", text)
	case 2:
		text = fmt.Sprintf("%s [XAR_ERR_ARCHIVE_EXTRACTION]", text)
	default:
		text = fmt.Sprintf("%s [%d]", text, int32(err))
	}

	xeString := C.xar_err_get_string(ctx)
	if xeString != nil {
		str := C.GoString(xeString)
		text = fmt.Sprintf("%s %s", text, string(str))
	}

	xeErrno := C.xar_err_get_errno(ctx)
	if xeErrno != 0 {
		str := C.GoString(C.strerror(xeErrno))
		text = fmt.Sprintf("%s (%d - %s)", text, int(xeErrno), string(str))
	}

	errCh <- errors.New(text)

	return C.int32_t(0)
}

func GetError() error {
	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}

func RegisterErrHandler(x *XAR) error {
	if x == nil {
		return errors.New("x == nil")
	}

	if !x.IsInitialized() {
		return errors.New("!x.IsInitialized()")
	}

	C.xar_register_errhandler(x.C(), C.err_handler(C.ErrHandler), nil)
	return nil
}
