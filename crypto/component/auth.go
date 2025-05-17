package component

import (
	"MASA/crypto/util"
	"crypto/rand"
	"fmt"

	"math/big"
)

// AA
type Auth struct {
	aid     int64              //AA的id
	sijalf  map[int64]*big.Int //记录不同aa发给自己的碎片
	sijbeta map[int64]*big.Int
}

type AuthSystem struct {
	prime *big.Int
	auths []Auth
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
	authsystem := &AuthSystem{
		prime: prime,
		auths: auths,
	}
	return true, authsystem, nil

	// 设置门限参数 (t,n)

	// 分发秘密分片

}
