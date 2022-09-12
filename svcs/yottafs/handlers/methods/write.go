package methods

import (
	"errors"
	"github.com/fxamacker/cbor/v2"
	"yottafs/iodriver"
)

func Write(ioReq iodriver.Request, driver iodriver.Interface) ([]byte, error) {

	resp, err := driver.Write(ioReq)
	if err != nil {
		return nil, errors.New("Read failed")
	}

	buff, err := cbor.Marshal(resp)
	if err != nil {
		return nil, errors.New("Read failed")
	}

	return buff, nil
}
