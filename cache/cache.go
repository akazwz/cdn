package cache

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func MakeCache(cachePath string, resp *http.Response) error {
	file, createErr := os.Create(cachePath)
	if createErr != nil {
		return createErr
	}
	defer func() {
		_ = file.Close()
	}()

	// write status code
	_, err := file.Write([]byte(fmt.Sprintf("%d\n", resp.StatusCode)))
	if err != nil {
		return err
	}

	// write headers
	for k, v := range resp.Header {
		_, err = file.Write([]byte(fmt.Sprintf("%s: %s\n", k, strings.Join(v, ""))))
		if err != nil {
			return err
		}
	}
	resp.Header.Set("X-Cache", "MISS")
	// write empty line
	_, err = file.Write([]byte("\r\n"))

	// write body
	var buf bytes.Buffer
	teaReader := io.TeeReader(resp.Body, &buf)
	_, err = io.Copy(file, teaReader)
	if err != nil {
		return err
	}
	resp.Body = io.NopCloser(&buf)
	return nil
}

func ParseCache(file *os.File) (*http.Response, error) {
	reader := bufio.NewReader(file)
	var off int64
	statusLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	off += int64(len(statusLine))
	// parse status code
	statusCode, err := strconv.Atoi(strings.TrimSpace(statusLine))
	if err != nil {
		return nil, err
	}

	headers := make(http.Header)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		off += int64(len(line))
		line = strings.TrimSpace(line)
		// space line: end of headers
		if line == "" {
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, err
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		headers.Add(key, value)
	}

	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	n := fi.Size() - off
	resp := &http.Response{
		StatusCode:    statusCode,
		Header:        headers,
		Body:          io.NopCloser(io.NewSectionReader(file, off, n)),
		ContentLength: n,
	}
	return resp, nil
}
