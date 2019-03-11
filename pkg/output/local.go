package output

import (
	"bytes"
	"os"

	"github.com/pkg/errors"
)

type LocalConfig struct {
	FilePath string `yaml:"filepath"`
}

type Local struct {
	config *LocalConfig
}

func NewLocalClient(config *LocalConfig) (*Local, error) {
	return &Local{
		config: config,
	}, nil
}

// Publish is Local file write
func (l *Local) Publish(buf *bytes.Buffer) error {
	file, err := os.OpenFile(l.config.FilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return errors.Wrapf(err, "failed to open file")
	}
	defer file.Close()

	file.Write(buf.Bytes())
	file.Write([]byte("\n"))

	return nil
}
