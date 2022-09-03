package direct

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/sys/unix"
)

func appendTo(path string, data []byte) error {

	fd, err := unix.Open(path, unix.O_RDWR|unix.O_DIRECT, 0766)
	defer unix.Close(fd)
	if err != nil {
		return err
	}

	var stat unix.Stat_t
	if err = unix.Fstat(fd, &stat); err != nil {
		return err
	}

	appendBlock := (stat.Size - 1) / BlockSize

	buffer := CallocAlignedBlock(1)

	//fmt.Println("append block is: ", appendBlock)

	if _, readErr := unix.Pread(fd, buffer, appendBlock*4096); readErr != nil {
		return readErr
	}

	terminationIndex := bytes.Index(buffer, []byte{0})
	if terminationIndex < 0 {
		panic("Termination index not found!")
	}
	fmt.Println("Termination index is: ", terminationIndex)

	/*writeBuffer := append(buffer[:terminationIndex], data...)

	blocksToWrite := len(writeBuffer)/BlockSize + 1

	for writeCounter := 0; writeCounter < blocksToWrite; writeCounter++ {

		lowerBound := writeCounter * 4096
		upperBound := lowerBound + 4096
		if upperBound > len(writeBuffer) {
			upperBound = len(writeBuffer)
		}

		buffer = CallocAlignedBlock(1)
		copy(buffer, writeBuffer[lowerBound:upperBound])
		offset := appendBlock*BlockSize + int64(lowerBound)
		//fmt.Println("Offset is: ", offset)
		_, readErr := unix.Pwrite(fd, buffer, offset)
		if readErr != nil {
			return readErr
		}
	}*/

	return nil
}

func compareAndAppend(path string, data []byte, aba string) error {

	return errors.New("method not implemented")
}
