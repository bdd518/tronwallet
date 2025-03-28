package main

import (
	"log"
	"sync"

	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/keys"
	"github.com/fbsobreira/gotron-sdk/pkg/mnemonic"
)

func main() {
	log.Println("Search lucky wallet")

	wg := sync.WaitGroup{}
	// 开 8 个携程去跑
	for n := range 8 {

		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			cnt := int64(0)

			for {
				cnt++

				// 生成助记词
				mn := mnemonic.Generate()

				// 助记词生成私钥
				private, _ := keys.FromMnemonicSeedAndPassphrase(mn, "", 0)

				// 私钥到处钱包地址
				aa := address.PubkeyToAddress(private.ToECDSA().PublicKey)

				// 取 后8位
				addr := aa.String()
				n := len(addr)
				last8 := addr[n-8:]

				// 检查是不是 XXXXYYYY
				// I7-11700 CPU 90%+ 跑了8个小时，没跑出一个来
				if checkPattern(last8) {
					log.Printf("Worker(%v)- try(%v) Found addr: %v [%v]\n", workerID, cnt, addr, mn)
				}
			}
		}(n)
	}

	wg.Wait()
}

// checkPattern 检查字符串是否为XXXXYYYY格式
func checkPattern(s string) bool {
	if len(s) < 8 {
		return false
	}
	return s[0] == s[1] && s[0] == s[2] && s[0] == s[3] &&
		s[4] == s[5] && s[4] == s[6] && s[4] == s[7]
}
