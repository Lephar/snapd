package osutil_test

import (
	"os"
	"path/filepath"
	
	. "gopkg.in/check.v1"
	
	"github.com/snapcore/snapd/osutil"
)

type ShredTestSuite struct{}

var _ = Suite(&ShredTestSuite{})

func (ts *ShredTestSuite) TestShredEmptyFileName(c *C) {
	c.Assert(osutil.Shred(""), NotNil)
}

func (ts *ShredTestSuite) TestShredFileDoesNotExist(c *C) {
	c.Assert(osutil.Shred("/i-do-not-exist"), NotNil)
}

func (ts *ShredTestSuite) TestShredFileExistsZeroSize(c *C) {
	fileName := filepath.Join(c.MkDir(), "emptyfile")
	file, _ := os.OpenFile(fileName, os.O_CREATE | os.O_WRONLY, 0777)
	_ = file.Close()
	
	c.Assert(osutil.Shred(fileName), IsNil)
}

func (ts *ShredTestSuite) TestShredFileExistsNoPermission(c *C) {
	fileName := filepath.Join(c.MkDir(), "readonlyfile")
	file, _ := os.OpenFile(fileName, os.O_CREATE | os.O_WRONLY, 0444)
	_, _ = file.Write([]byte("Lorem ipsum dolor sit amet"))
	_ = file.Close()

	c.Assert(osutil.Shred(fileName), NotNil)
}

func (ts *ShredTestSuite) TestShredFileExistsSimple(c *C) {
	fileName := filepath.Join(c.MkDir(), "textfile")
	file, _ := os.OpenFile(fileName, os.O_CREATE | os.O_WRONLY, 0777)
	_, _ = file.Write([]byte("Lorem ipsum dolor sit amet"))
	_ = file.Close()

	c.Assert(osutil.Shred(fileName), IsNil)
}

