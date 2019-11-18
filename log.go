package dlog

import (
	"github.com/drep-project/DREP-Chain/common/fileutil"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"io"
	"os"
)

var DEBUG = false

var (
	ostream Handler
	glogger *GlogHandler
)

func init() {
	usecolor := (isatty.IsTerminal(os.Stderr.Fd()) || isatty.IsCygwinTerminal(os.Stderr.Fd())) && os.Getenv("TERM") != "dumb"
	output := io.Writer(os.Stderr)
	if usecolor {
		output = colorable.NewColorableStderr()
	}
	ostream = StreamHandler(output, TerminalFormat(usecolor))
	glogger = NewGlogHandler(ostream)
}

func SetUp(dataDir string, loglevel int, vmModule string, backtraceAt string) error {
	if dataDir != "" {
		if !fileutil.IsDirExists(dataDir) {
			err := os.MkdirAll(dataDir, 0777)
			if err != nil {
				return err
			}
		}

		rfh, err := SyncRotatingFileHandler(
			dataDir,
			262144,
			JSONFormatOrderedEx(false, true),
		)
		if err != nil {
			return err
		}
		glogger.SetHandler(MultiHandler(ostream, rfh))
	}
	glogger.Verbosity(Lvl(loglevel))
	glogger.Vmodule(vmModule)
	glogger.BacktraceAt(backtraceAt)
	Root().SetHandler(glogger)
	return nil
}

func SetVerbosity(loglevel Lvl)  {
	 glogger.Verbosity(loglevel)
}

func SetVmodule(vmModule string) error {
	return glogger.Vmodule(vmModule)
}

func SetBacktraceAt(backtraceAt string) error {
	return glogger.BacktraceAt(backtraceAt)
}