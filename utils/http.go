package utils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

var (
	userAgent = `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36 Edg/124.0.0.0`
)

var (
	ErrFileSize      = errors.New("file size error")
	ErrRequestFailed = errors.New("request failed")
)

func FetchFileSize(url, referer string) (fileSize int64, err error) {
	var req *http.Request
	req, err = http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", userAgent)
	if len(referer) > 0 {
		req.Header.Set("Referer", referer)
	}

	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fileSize, ErrRequestFailed
	}

	for key, value := range resp.Header {
		if strings.ToLower(key) == "content-length" && len(value) > 0 {
			fileSize, _ = strconv.ParseInt(value[0], 10, 64)
		}
	}

	return fileSize, nil
}

func ParallelDown(fileName, url, referer string, fileSize int64, parallelSize int, readonlyRealtimeSize *int64) (err error) {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	var (
		perSize     = int64(math.Ceil(float64(fileSize) / float64(parallelSize)))
		ctx, cancel = context.WithCancel(context.Background())
		wa          = &sync.WaitGroup{}
	)

	defer func() {
		cancel()

		fi, _ := file.Stat()
		if fi == nil || fi.Size() != fileSize {
			err = ErrFileSize
		}

		_ = file.Close()
		if err != nil {
			_ = os.Remove(fileName)
		}
	}()

	wa.Add(parallelSize)

	for i := 0; i < parallelSize; i++ {
		go func(index int) {
			defer wa.Done()

			offset := int64(index) * perSize
			end := offset + perSize - 1
			mReq, er := http.NewRequest(http.MethodGet, url, nil)
			if er != nil {
				err = er
				cancel()
				return
			}

			mReq.Header.Set("User-Agent", userAgent)
			if len(referer) > 0 {
				mReq.Header.Set("Referer", referer)
			}

			if index == parallelSize-1 {
				mReq.Header.Set("Range", fmt.Sprintf("bytes=%d-", offset))
			} else {
				mReq.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", offset, end))
			}

			mResp, er := http.DefaultClient.Do(mReq)
			if er != nil {
				err = er
				cancel()
				return
			}

			defer mResp.Body.Close()

			if mResp.StatusCode != http.StatusPartialContent && mResp.StatusCode != http.StatusOK {
				err = ErrRequestFailed
				cancel()
				return
			}

			var (
				buf     = make([]byte, 1024*1024)
				readLen int
			)

			for {
				select {
				case <-ctx.Done():
					return
				default:
				}

				readLen, er = mResp.Body.Read(buf)
				if readLen > 0 {
					atomic.AddInt64(readonlyRealtimeSize, int64(readLen))

					_, er = file.WriteAt(buf[:readLen], offset)
					if er != nil {
						err = er
						cancel()
						return
					}

					offset += int64(readLen)
					if offset >= end {
						return
					}
				}

				if er != nil {
					if er == io.EOF {
						return
					}

					err = er
					cancel()
					return
				}
			}
		}(i)
	}

	wa.Wait()

	return
}
