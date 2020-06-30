# wheel timer
[![GoDoc](https://godoc.org/github.com/thinkgos/wheel?status.svg)](https://godoc.org/github.com/thinkgos/wheel)
[![Build Status](https://travis-ci.org/thinkgos/wheel.svg?branch=master)](https://travis-ci.org/thinkgos/wheel)
[![codecov](https://codecov.io/gh/thinkgos/wheel/branch/master/graph/badge.svg)](https://codecov.io/gh/thinkgos/wheel)
![Action Status](https://github.com/thinkgos/wheel/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/thinkgos/wheel)](https://goreportcard.com/report/github.com/thinkgos/wheel)
[![Licence](https://img.shields.io/github/license/thinkgos/wheel)](https://raw.githubusercontent.com/thinkgos/wheel/master/LICENSE)  

 - 五层时间轮: 主级加四个层级
 - 插入,删除,修改时间,扫描超时条目时间复杂度o(1)
 - 默认时间精度为1m.
 - 最大时间受限于时基精度,时间精度1ms最大可定时时间为49.71天,所以可定时最大时间为49.71天*${时基精度(ms)}
