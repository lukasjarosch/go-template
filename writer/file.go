package writer

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

)

type FileWriter struct {
	opts Options
	path string
}

var ErrNoOp = errors.New("nothing done")

func NewFileWriter(path string, option ...Option) *FileWriter {
	fw := &FileWriter{
		opts: newOptions(option...),
		path: path,
	}

	return fw
}

func (fw *FileWriter) WriteFile(data []byte) error {
	if err := fw.ensurePath(); err != nil {
		return err
	}

	if !fw.PathExists() {
		if err := fw.pathCreate(); err != nil {
		   	return errors.Wrap(err, "unable to create target file")
		}
		if err := fw.write(data); err != nil {
			return err
		}
		return nil
	}

	if fw.opts.Append == true && fw.opts.Append == fw.opts.Overwrite {
		return fmt.Errorf("appendWrite and overwrite flags are set, you have to choose one of them")
	}

	if fw.opts.Overwrite {
		if err := fw.write(data); err != nil {
		   	return err
		}
		return nil
	}

	if fw.opts.Append {
		if err := fw.appendWrite(data); err != nil {
			return err
		}
		return nil
	}

	return ErrNoOp
}

func (fw *FileWriter) write(data []byte) error {
	return ioutil.WriteFile(fw.path, data, 0644)
}

func (fw *FileWriter) appendWrite(data []byte) error {
	f, err := os.OpenFile(fw.path, os.O_APPEND|os.O_WRONLY, 0644)
	defer f.Close()
	if err != nil {
	    return errors.Wrap(err, "unbale to open file for appendWrite-write")
	}
	if _, err := f.Write(data); err != nil {
		return errors.Wrap(err, "unable to appendWrite-write file")
	}
	return nil
}

func (fw *FileWriter) PathExists() bool {
	if _, err := os.Stat(fw.path); os.IsNotExist(err) {
		return false
	}
	return true
}

func (fw *FileWriter) pathCreate() error {
	if _, err := os.Create(fw.path); err != nil {
		return err
	}
	return nil
}

func (fw *FileWriter) ensurePath() error {
	if fw.path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	targetFolder := filepath.Dir(fw.path)
	if _, err := os.Stat(targetFolder); os.IsNotExist(err) {
		if err := os.MkdirAll(targetFolder, 0755); err != nil {
			return errors.Wrap(err, "writer could not ensure target folder")
		}
	}

	// ensure permission for existing
	testFile := filepath.Join(targetFolder, "tester")
	if _, err := os.Create(testFile); os.IsPermission(err) {
		if err := os.Chmod(targetFolder, 0755); err != nil {
			err = errors.Wrap(err, "writer path not writeable and chmod failed")
			return err
		}
	}
	if err := os.Remove(testFile); err != nil {
		return errors.Wrap(err, "failed to remove writer permission tester")
	}

	return nil

	return nil
}

type Option func(*Options)
