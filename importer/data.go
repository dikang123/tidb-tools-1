// Copyright 2016 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/pingcap/tidb-tools/pkg/base26"
)

var defaultStep int64 = 1

type datum struct {
	sync.Mutex

	intValue    int64
	minIntValue int64
	maxIntValue int64
	timeValue   time.Time
	step        int64

	init     bool
	useRange bool
}

func newDatum() *datum {
	return &datum{intValue: -1, step: 1}
}

func (d *datum) setInitInt64Value(step int64, min int64, max int64) {
	d.Lock()
	defer d.Unlock()

	if d.init {
		return
	}

	d.step = step

	if min != -1 {
		d.minIntValue = min
		d.intValue = min
	}

	if min < max {
		d.maxIntValue = max
		d.useRange = true
	}

	d.init = true
}

func (d *datum) uniqInt64() int64 {
	d.Lock()
	defer d.Unlock()

	data := d.intValue
	if d.useRange {
		if d.intValue+d.step > d.maxIntValue {
			return data
		}
	}

	d.intValue += d.step
	return data
}

func (d *datum) uniqFloat64() float64 {
	data := d.uniqInt64()
	return float64(data)
}

func (d *datum) uniqString(n int) string {
	d.Lock()
	d.intValue++
	data := d.intValue
	d.Unlock()


	return base26.Encode(data, n)
}

func (d *datum) uniqTime() string {
	d.Lock()
	defer d.Unlock()

	if d.timeValue.IsZero() {
		d.timeValue = time.Now()
	} else {
		d.timeValue = d.timeValue.Add(time.Duration(d.step) * time.Second)
	}

	return fmt.Sprintf("%02d:%02d:%02d", d.timeValue.Hour(), d.timeValue.Minute(), d.timeValue.Second())
}

func (d *datum) uniqDate() string {
	d.Lock()
	defer d.Unlock()

	if d.timeValue.IsZero() {
		d.timeValue = time.Now()
	} else {
		d.timeValue = d.timeValue.AddDate(0, 0, int(d.step))
	}

	return fmt.Sprintf("%04d-%02d-%02d", d.timeValue.Year(), d.timeValue.Month(), d.timeValue.Day())
}

func (d *datum) uniqTimestamp() string {
	d.Lock()
	defer d.Unlock()

	if d.timeValue.IsZero() {
		d.timeValue = time.Now()
	} else {
		d.timeValue = d.timeValue.Add(time.Duration(d.step) * time.Second)
	}

	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d",
		d.timeValue.Year(), d.timeValue.Month(), d.timeValue.Day(),
		d.timeValue.Hour(), d.timeValue.Minute(), d.timeValue.Second())
}

func (d *datum) uniqYear() string {
	d.Lock()
	defer d.Unlock()

	if d.timeValue.IsZero() {
		d.timeValue = time.Now()
	} else {
		d.timeValue = d.timeValue.AddDate(int(d.step), 0, 0)
	}

	return fmt.Sprintf("%04d", d.timeValue.Year())
}
