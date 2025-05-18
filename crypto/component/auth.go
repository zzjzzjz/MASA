package component

import (
	"MASA/crypto/util"
	"crypto/rand"
	"fmt"
	bn128 "github.com/fentec-project/bn256"

	"math/big"
)

type Sk struct {
	alpha *big.Int
	beta  *big.Int
}
type Pk struct {
	Alpha *bn128.GT
	Beta  *bn128.G1
}

// AA
type Auth struct {
	aid     int64              //AA的id
	sijalf  map[int64]*big.Int //记录不同aa发给自己的碎片
	sijbeta map[int64]*big.Int
	sk      Sk
	pk      Pk
}

type AuthSystem struct {
	prime *big.Int
	auths []Auth
	G1    *bn128.G1
	G2    *bn128.G2
	Gt    *bn128.GT
}

func Setup(n int, t int, aids []int64) (bool, *AuthSystem, error) { //n--aa数量，t--aa阀值 aids--aa的aid
	//先检查aids的长度是否等于n的长度
	if len(aids) != n {
		fmt.Println("aids不等于n")
		return false, nil, nil
	}
	if n < t {
		fmt.Println("n应该大于t")
		return false, nil, nil
	}
	//shamir需要的prime
	prime, _ := rand.Prime(rand.Reader, 128)

	//生成各种的sk（alpha，beta）
	alpha := make([]*big.Int, n)
	beta := make([]*big.Int, n)
	auths := make([]Auth, n)
	for i := 0; i < n; i++ {
		alpha[i] = util.RandomInt()
		beta[i] = util.RandomInt()
		auths[i] = Auth{
			aid:     aids[i],                  // 设置 aid
			sijalf:  make(map[int64]*big.Int), // 初始化空 map
			sijbeta: make(map[int64]*big.Int), // 初始化空 map
		}
		auths[i].sk.alpha = new(big.Int)
		auths[i].sk.beta = new(big.Int)

	}
	x := util.ConvertInt64SliceToBigIntSlice(aids)

	for i := 0; i < n; i++ {
		alphashares, _ := util.Distribute(alpha[i], prime, t, n, x)
		betashares, _ := util.Distribute(beta[i], prime, t, n, x)
		for j := 0; j < n; j++ {
			auths[j].sijalf[aids[i]] = alphashares[j].Y
			auths[j].sijbeta[aids[i]] = betashares[j].Y
		}

	}

	//公私钥生成
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			auths[i].sk.alpha = auths[i].sk.alpha.Add(auths[i].sk.alpha, auths[i].sijalf[aids[j]])
		}

	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			auths[i].sk.beta = auths[i].sk.beta.Add(auths[i].sk.beta, auths[i].sijbeta[aids[j]])
		}

	}
	//公钥部分，需要构建系统的g1，g2，gt
	gen1 := new(bn128.G1).ScalarBaseMult(big.NewInt(1))
	gen2 := new(bn128.G2).ScalarBaseMult(big.NewInt(1))
	Gt := bn128.Pair(gen1, gen2)
	for i := 0; i < n; i++ {
		auths[i].pk.Alpha = new(bn128.GT).ScalarMult(Gt, auths[i].sk.alpha)
		auths[i].pk.Beta = new(bn128.G1).ScalarMult(gen1, auths[i].sk.beta)
	}

	authsystem := &AuthSystem{
		prime: prime,
		G1:    gen1,
		G2:    gen2,
		Gt:    Gt,
		auths: auths,
	}
	return true, authsystem, nil

}
func setupPK(pks []Pk, aids []int64) Pk {
	n := len(pks)
	for i := 0; i < n; i++ {

	}
}
