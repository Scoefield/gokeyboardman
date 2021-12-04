package main

import "errors"

// bad
func parseBad(lenControlByUser int, data[] byte) {
	size := lenControlByUser
	//对外部传入的size，进行长度判断以免导致panic
	buffer := make([]byte, size)
	copy(buffer, data)
}

// good
func parseGood(lenControlByUser int, data[] byte) ([]byte, error){
	size := lenControlByUser
	//限制外部可控的长度大小范围
	if size > 64*1024*1024 {
		return nil, errors.New("value too large")
	}
	buffer := make([]byte, size)
	copy(buffer, data)
	return buffer, nil
}
