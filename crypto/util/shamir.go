package util

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// Share 表示一个秘密分片
type Share struct {
	X, Y *big.Int
}

// Polynomial 表示多项式
type Polynomial struct {
	coefficients []*big.Int
	modulus      *big.Int
}

// NewPolynomial 创建t-1次多项式
func NewPolynomial(degree int, modulus *big.Int) Polynomial {
	coefficients := make([]*big.Int, degree+1)
	for i := 0; i <= degree; i++ {
		coefficients[i], _ = rand.Int(rand.Reader, modulus)
	}
	return Polynomial{coefficients, modulus}
}

// Eval 计算多项式在x处的值
func (p *Polynomial) Eval(x *big.Int) *big.Int {
	result := new(big.Int).Set(p.coefficients[len(p.coefficients)-1])
	for i := len(p.coefficients) - 2; i >= 0; i-- {
		result.Mul(result, x)
		result.Add(result, p.coefficients[i])
		result.Mod(result, p.modulus)
	}
	return result
}

// Distribute 分发秘密分片 (t,n)门限方案
// 增加了坐标x作为函数输入
func Distribute(secret *big.Int, modulus *big.Int, t, n int, x []*big.Int) ([]Share, error) {
	if t > n {
		return nil, fmt.Errorf("threshold t cannot be greater than n")
	}

	p := NewPolynomial(t-1, modulus)
	if secret != nil {
		p.coefficients[0] = secret
	} else {
		secret = p.coefficients[0]
	}

	shares := make([]Share, n)
	for i := 0; i < n; i++ {
		y := p.Eval(x[i])
		shares[i] = Share{x[i], y}
	}
	return shares, nil
}

// Lagrange 拉格朗日插值恢复秘密
func Lagrange(x, modulus *big.Int, shares ...Share) *big.Int {
	result := big.NewInt(0)
	for i, share := range shares {
		numerator := big.NewInt(1)
		denominator := big.NewInt(1)

		for j, other := range shares {
			if i != j {
				// 计算分子: (x - xj)
				tmp := new(big.Int).Sub(x, other.X)
				numerator.Mul(numerator, tmp)

				// 计算分母: (xi - xj)
				tmp = new(big.Int).Sub(share.X, other.X)
				denominator.Mul(denominator, tmp)
			}
		}

		// 计算: yi * (分子/分母)
		tmp := new(big.Int).Mul(share.Y, numerator)
		denominator.ModInverse(denominator, modulus)
		tmp.Mul(tmp, denominator)
		result.Add(result, tmp)
	}
	return result.Mod(result, modulus)
}
