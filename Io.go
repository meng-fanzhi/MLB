package main

import (
	"fmt"
	"net"
	"public/loginit"
	"strconv"
	"time"
)

//func LimitReader(r Reader, n int64) Reader { return &LimitedReader{r, n} }

func copyBuffer(dst net.Conn, src net.Conn, idla_timeout int, buf []byte) (written int64, err error) {
	if buf == nil {
		size := buffer * 1024
		buf = make([]byte, size)
		//fmt.Print(buf)
	}
	for {
		if idla_timeout != 0 {
			ideaTime, _ := time.ParseDuration(strconv.Itoa(idla_timeout) + "ms")
			err = src.SetReadDeadline(time.Now().Add(ideaTime))
		}
		if err != nil {
			loginit.Warning.Println(err)
		}

		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			//fmt.Printf("written:%v\n",written)
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != EOF {
				err = er
			}
			break
		}
	}
	return written, err
}

// 数据计算
func formatFileSize(fileSize int64) (size string) {
	if fileSize < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}
