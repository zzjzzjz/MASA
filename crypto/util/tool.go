package util

import (
	bn128 "github.com/fentec-project/bn256"
	"github.com/fentec-project/gofe/data"
	"github.com/fentec-project/gofe/sample"
	"math/big"
)

// 随机大数
func RandomInt() *big.Int {
	v, _ := data.NewRandomVector(1, sample.NewUniform(bn128.Order))
	return v[0]
}
func ConvertInt64SliceToBigIntSlice(aids []int64) []*big.Int {
	x := make([]*big.Int, len(aids)) // 创建与 aids 长度相同的切片
	for i, v := range aids {
		x[i] = big.NewInt(v) // 将每个 int64 转换为 *big.Int
	}
	return x
}
