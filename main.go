package main

import (
	"MASA/crypto/component"
)

func main() {
	n := 5                                   // AA 数量
	t := 3                                   // 阈值
	aids := []int64{101, 102, 103, 104, 105} // AA 的 ID 列表

	// 2. 调用 Setup
	component.Setup(n, t, aids)
}
