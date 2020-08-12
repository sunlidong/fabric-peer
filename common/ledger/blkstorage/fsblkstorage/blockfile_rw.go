/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package fsblkstorage

import (
	"os"

	"github.com/pkg/errors"
)

////  WRITER ////
type blockfileWriter struct {
	filePath string
	file     *os.File
}

// 新建块文件写入器
func newBlockfileWriter(filePath string) (*blockfileWriter, error) {
	writer := &blockfileWriter{filePath: filePath}
	return writer, writer.open()
}

// 截断文件
func (w *blockfileWriter) truncateFile(targetSize int) error {
	fileStat, err := w.file.Stat()
	if err != nil {
		return errors.Wrapf(err, "error truncating the file [%s] to size [%d]", w.filePath, targetSize)
	}
	if fileStat.Size() > int64(targetSize) {
		w.file.Truncate(int64(targetSize))
	}
	return nil
}

// 添加 append
func (w *blockfileWriter) append(b []byte, sync bool) error {
	_, err := w.file.Write(b)
	if err != nil {
		return err
	}
	if sync {
		return w.file.Sync()
	}
	return nil
}

// 打开
func (w *blockfileWriter) open() error {
	file, err := os.OpenFile(w.filePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		return errors.Wrapf(err, "error opening block file writer for file %s", w.filePath)
	}
	w.file = file
	return nil
}

// 关闭
func (w *blockfileWriter) close() error {
	return errors.WithStack(w.file.Close())
}

////  READER ////
type blockfileReader struct {
	file *os.File
}

// 新的块文件读取器
func newBlockfileReader(filePath string) (*blockfileReader, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0600)
	if err != nil {
		return nil, errors.Wrapf(err, "error opening block file reader for file %s", filePath)
	}
	reader := &blockfileReader{file}
	return reader, nil
}

// 读取
func (r *blockfileReader) read(offset int, length int) ([]byte, error) {
	b := make([]byte, length)
	_, err := r.file.ReadAt(b, int64(offset))
	if err != nil {
		return nil, errors.Wrapf(err, "error reading block file for offset %d and length %d", offset, length)
	}
	return b, nil
}

// 关闭
func (r *blockfileReader) close() error {
	return errors.WithStack(r.file.Close())
}
