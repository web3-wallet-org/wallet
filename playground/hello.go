package main

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func main() {
	fmt.Println("🎮 Welcome to Playground!")
	fmt.Println("================================")

	// 示例 1: 基础 Hello World
	fmt.Println("\n1️⃣ Basic Hello:")
	fmt.Println("   Hello, Wallet Developer!")

	// 示例 2: 测试以太坊地址
	fmt.Println("\n2️⃣ Ethereum Address:")
	addr := common.HexToAddress("0x1234567890123456789012345678901234567890")
	fmt.Printf("   Address: %s\n", addr.Hex())
	fmt.Printf("   Checksum: %s\n", addr.String())

	// 示例 3: 大数计算
	fmt.Println("\n3️⃣ Big Number Calculation:")
	oneEther := new(big.Int).SetUint64(1e18) // 1 ETH = 10^18 wei
	fmt.Printf("   1 ETH = %s wei\n", oneEther.String())

	gasPrice := new(big.Int).SetUint64(20e9) // 20 gwei
	gasLimit := uint64(21000)
	fee := new(big.Int).Mul(gasPrice, big.NewInt(int64(gasLimit)))
	fmt.Printf("   Gas Fee (21000 * 20 gwei) = %s wei\n", fee.String())

	// 示例 4: 计算以太单位
	feeInGwei := new(big.Int).Div(fee, big.NewInt(1e9))
	fmt.Printf("   Gas Fee = %s gwei\n", feeInGwei.String())

	fmt.Println("\n================================")
	fmt.Println("✨ 在这里快速测试你的代码吧！")
}
