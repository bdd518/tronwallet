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
		x := n
		wg.Add(1)
		go func() {
			defer wg.Done()

			i := int64(0)

			for {
				i++

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
				if last8[0] == last8[1] &&
					last8[0] == last8[2] &&
					last8[0] == last8[3] &&
					last8[4] == last8[5] &&
					last8[4] == last8[6] &&
					last8[4] == last8[7] {
					log.Printf("%v-%v Found: %v[%v]\n", x, i, addr, mn)
				}
			}
		}()
	}

	wg.Wait()
}
