package osutil

import (
	"fmt"
	"log"
	"os"
)

func overwriteRandom(file *os.File) error {
	// How many times the file will be overwritten
	const overwriteCount = 10

	randomFile, err := os.Open("/dev/urandom")
	if err != nil {
		return err
	}
	defer func(randomFile *os.File) {
		if err := randomFile.Close(); err != nil {
			log.Println(err)
		}
	}(randomFile)

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	// Get the file size so we can read enough random data from /dev/urandom
	fileSize := stat.Size()
	
	// For performance reasons, read all the necessary random data at once
	// instead of reading just enough data every iteration. So the buffer
	// size must be enough to hold 10 times the file size.
	// If the security is a bigger concern than performance, and storing the
	// random data on the memory can be exploited, reading just before we
	// need may be the better choice. 
	buffer := make([]byte, overwriteCount * fileSize)
	readSize, err := randomFile.Read(buffer)
	if err != nil {
		return err
	} else if readSize != len(buffer) {
		return fmt.Errorf("cannot read %d bytes from '/dev/urandom'", len(buffer))
	}

	for i := int64(0); i < overwriteCount; i++ {
		// Rewind to the beginning of the file every time, so it is overwritten
		seekPos, err := file.Seek(0, 0)
		if err != nil {
			return err
		} else if seekPos != 0 {
			return fmt.Errorf("cannot rewind %q for the next pass", file.Name())
		}

		// Calculate the data slice we need to use to overwrite the file
		// depending on the current iteration
		sliceBegin := i * fileSize
		sliceEnd := sliceBegin + fileSize

		writeSize, err := file.Write(buffer[sliceBegin : sliceEnd])
		if err != nil {
			return err
		} else if int64(writeSize) != fileSize {
			return fmt.Errorf("cannot write %d bytes to %q", fileSize, file.Name())
		}
	}

	return nil
}

// This function opens the file and overwrites it 10 times with random data
// obtained from /dev/urandom device, then deletes it
func Shred(fileName string) error {
	// Open the file with minimum permissions possible
	file, err := os.OpenFile(fileName, os.O_WRONLY, 0200)
	if err != nil {
		return err
	}

	// We cannot defer close here because we need to close the file BEFORE
	// we remove it. Regardless of the errors in overwriteRandom, we close
	// the file, THEN check the errors
	errWrite := overwriteRandom(file)
	err = file.Close()

	if errWrite != nil {
		return errWrite
	}
	if err != nil {
		return err
	}

	// If everything goes as expected, remove the file
	if err = os.Remove(fileName); err != nil {
		return err
	}

	return nil
}
