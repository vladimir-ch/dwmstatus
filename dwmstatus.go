// Copyright (c) 2015 Vladimir Chalupecky
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Based on
//  https://github.com/oniichaNj/go-dwmstatus
// and
//  https://github.com/phacops/gods

package main

// #cgo LDFLAGS: -lX11
// #include <X11/Xlib.h>
import "C"

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var dpy = C.XOpenDisplay(nil)

func setStatus(status string) {
	C.XStoreName(dpy, C.XDefaultRootWindow(dpy), C.CString(status))
	C.XSync(dpy, 1)
}

func cpuLoad() string {
	loadavg, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return "CPUERR"
	}
	fields := strings.Split(string(loadavg), " ")
	return "CPU " + strings.Join(fields[:3], " ")
}

func usedMemory() string {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return "MEMERR"
	}
	defer file.Close()

	var (
		total, free, cached, slab, buffers int

		prop string
		v    int
	)
	meminfo := bufio.NewScanner(file)
	for meminfo.Scan() {
		_, err := fmt.Sscanf(meminfo.Text(), "%s %d", &prop, &v)
		if err != nil {
			return "MEMERR"
		}
		switch prop {
		case "MemTotal:":
			total = v
		case "MemFree:":
			free = v
		case "Cached:":
			cached = v
		case "Slab:":
			slab = v
		case "Buffers:":
			buffers = v
		}
	}
	used := total - free - cached - slab - buffers
	if used < 0 {
		used = total - free
	}
	return fmt.Sprintf("MEM %d (%d%%)", used, 100*used/total)
}

func datetime() string {
	return time.Now().Local().Format("Mon 2006/01/02 15:04:05")
}

func main() {
	for {
		status := [...]string{
			cpuLoad(),
			usedMemory(),
			datetime(),
		}
		setStatus(strings.Join(status[:], " :: "))
		now := time.Now()
		time.Sleep(now.Truncate(time.Second).Add(time.Second).Sub(now))
	}
}
