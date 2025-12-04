package main

import (
	"crypto/ecdsa"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
)

func main() {
	fmt.Println("🔐 BIP39 & BIP44 测试")
	fmt.Println("==========================================")

	// ============================================
	// Part 1: BIP39 - 生成助记词和种子
	// ============================================
	fmt.Println("\n📝 Part 1: BIP39 - 助记词生成")
	fmt.Println("------------------------------------------")

	// 1.1 生成随机熵（128位 = 12个单词，256位 = 24个单词）
	entropy, err := bip39.NewEntropy(128) // 128 bits
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("1. 随机熵 (Hex): %x\n", entropy)

	// 1.2 从熵生成助记词
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("2. 助记词 (12 words):\n   %s\n", mnemonic)

	// 1.3 验证助记词是否有效
	isValid := bip39.IsMnemonicValid(mnemonic)
	fmt.Printf("3. 助记词验证: %v\n", isValid)

	// 1.4 从助记词生成种子（可选密码）
	password := "" // 可以设置为 "my-password" 增加安全性
	seed := bip39.NewSeed(mnemonic, password)
	fmt.Printf("4. 种子 (Seed): %x\n", seed[:32]) // 只显示前32字节

	// ============================================
	// Part 2: BIP44 - 分层确定性钱包
	// ============================================
	fmt.Println("\n🌳 Part 2: BIP44 - HD 钱包路径派生")
	fmt.Println("------------------------------------------")

	// 2.1 从助记词创建 HD 钱包
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Fatal(err)
	}

	// 2.2 标准以太坊路径: m/44'/60'/0'/0/0
	// - 44': BIP44 标准
	// - 60': 以太坊币种类型
	// - 0': 账户 0
	// - 0: 外部链（接收地址）
	// - 0: 地址索引 0
	fmt.Println("\n以太坊标准路径格式:")
	fmt.Println("m / 44' / 60' / 0' / 0 / address_index")
	fmt.Println("     │     │     │    │        │")
	fmt.Println("     │     │     │    │        └─ 地址索引")
	fmt.Println("     │     │     │    └────────── 0=接收 1=找零")
	fmt.Println("     │     │     └─────────────── 账户索引")
	fmt.Println("     │     └───────────────────── 60=以太坊")
	fmt.Println("     └─────────────────────────── BIP44 标准")

	// 2.3 派生多个地址
	fmt.Println("\n派生的以太坊地址:")
	for i := 0; i < 5; i++ {
		// 构造标准路径
		path := hdwallet.MustParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/0/%d", i))

		// 派生账户
		account, err := wallet.Derive(path, false)
		if err != nil {
			log.Fatal(err)
		}

		// 获取私钥
		privateKey, err := wallet.PrivateKey(account)
		if err != nil {
			log.Fatal(err)
		}

		// 获取公钥
		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatal("error casting public key to ECDSA")
		}

		// 从公钥生成地址
		address := crypto.PubkeyToAddress(*publicKeyECDSA)

		fmt.Printf("地址 %d: %s\n", i, address.Hex())
		fmt.Printf("       路径: m/44'/60'/0'/0/%d\n", i)
		if i == 0 {
			// 只显示第一个地址的私钥作为示例
			fmt.Printf("       私钥: %x\n", crypto.FromECDSA(privateKey))
		}
		fmt.Println()
	}

	// ============================================
	// Part 3: 不同账户和链
	// ============================================
	fmt.Println("🔗 Part 3: 不同账户派生")
	fmt.Println("------------------------------------------")

	// 账户 0 和 账户 1 的第一个地址
	accounts := []string{
		"m/44'/60'/0'/0/0", // 账户 0 地址 0
		"m/44'/60'/1'/0/0", // 账户 1 地址 0
		"m/44'/60'/2'/0/0", // 账户 2 地址 0
	}

	for _, pathStr := range accounts {
		path := hdwallet.MustParseDerivationPath(pathStr)
		account, err := wallet.Derive(path, false)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s → %s\n", pathStr, account.Address.Hex())
	}

	// ============================================
	// Part 4: 从已有助记词恢复
	// ============================================
	fmt.Println("\n♻️  Part 4: 从助记词恢复钱包")
	fmt.Println("------------------------------------------")

	// 使用相同的助记词恢复钱包
	recoveredWallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Fatal(err)
	}

	// 验证恢复的地址是否一致
	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	originalAccount, _ := wallet.Derive(path, false)
	recoveredAccount, _ := recoveredWallet.Derive(path, false)

	fmt.Printf("原始地址:   %s\n", originalAccount.Address.Hex())
	fmt.Printf("恢复后地址: %s\n", recoveredAccount.Address.Hex())
	fmt.Printf("地址匹配:   %v ✅\n", originalAccount.Address == recoveredAccount.Address)

	// ============================================
	// Part 5: 实用功能演示
	// ============================================
	fmt.Println("\n🛠️  Part 5: 实用功能")
	fmt.Println("------------------------------------------")

	// 5.1 检查地址是否属于钱包
	path = hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, _ := wallet.Derive(path, false)
	testAddress := account.Address

	fmt.Printf("测试地址: %s\n", testAddress.Hex())

	// 遍历前100个地址查找
	found := false
	for i := 0; i < 100; i++ {
		path := hdwallet.MustParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/0/%d", i))
		acc, _ := wallet.Derive(path, false)
		if acc.Address == testAddress {
			fmt.Printf("✅ 找到匹配地址，路径: %s\n", path.String())
			found = true
			break
		}
	}
	if !found {
		fmt.Println("❌ 未找到匹配地址")
	}

	// 5.2 从私钥获取地址（验证推导正确性）
	privateKey, _ := wallet.PrivateKey(account)
	publicKey := privateKey.Public()
	publicKeyECDSA := publicKey.(*ecdsa.PublicKey)
	derivedAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	fmt.Printf("\n从私钥推导地址: %s\n", derivedAddress.Hex())
	fmt.Printf("地址匹配: %v ✅\n", derivedAddress == testAddress)

	fmt.Println("\n==========================================")
	fmt.Println("✨ 测试完成！")
	fmt.Println("\n💡 提示:")
	fmt.Println("   - 助记词是你的主密钥，妥善保管！")
	fmt.Println("   - 使用 password 可以增加额外安全层")
	fmt.Println("   - 标准路径确保钱包间的兼容性")
	fmt.Println("   - 永远不要在生产环境打印私钥！")
}
