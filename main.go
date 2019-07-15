package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tyler-smith/go-bip39"
)

// go get -u github.com/tyler-smith/go-bip39
func test_mnemonic() {
	//Entropy 生成
	b, err := bip39.NewEntropy(128)
	if err != nil {
		log.Panic("failed to NewEntropy:", err, b)
	}

	fmt.Println(b)

	//生成助记词
	nm, err := bip39.NewMnemonic(b)
	if err != nil {
		log.Panic("failed to NewMnemonic:", err)
	}
	fmt.Println(nm)
}

//测试助记词有效
func test_ganache() {
	nm := "august human human affair mechanic night verb metal embark marine orient million"

	// 助记词转化为种子 - > 账户地址
	// 先推导路径，再获得钱包
	path := MustParseDerivationPath("m/44'/60'/0'/0/0")

	wallet, err := NewFromMnemonic(nm, "")

	if err != nil {
		log.Panic("failed to NewFromMnemonic:", err)
	}

	account, err := wallet.Derive(path, false)
	if err != nil {
		log.Panic("failed to Derive:", err)
	}

	fmt.Println(account.Address.Hex())

	path = MustParseDerivationPath("m/44'/60'/0'/0/2")

	account, err = wallet.Derive(path, false)
	if err != nil {
		log.Panic("failed to Derive:", err)
	}

	fmt.Println(account.Address.Hex())

}

//测试转账
func test_transfer() {

	// 助记词转化为种子 - > 账户地址
	// 先推导路径，再获得钱包
	mnemonic := "august human human affair mechanic night verb metal embark marine orient million"
	wallet, err := NewFromMnemonic(mnemonic, "")

	if err != nil {
		log.Panic("failed to NewFromMnemonic:", err)
	}
	// path := MustParseDerivationPath("m/44'/60'/0'/0/1")

	// account, err := wallet.Derive(path, true)
	// if err != nil {
	// 	log.Panic("failed to Derive:", err)
	// }

	// mnemonic := "august human human affair mechanic night verb metal embark marine orient million"
	// wallet, err := NewFromMnemonic(mnemonic, "")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	path := MustParseDerivationPath("m/44'/60'/0'/0/1") //第2个账户地址
	account, err := wallet.Derive(path, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(account.Address.Hex())

	// pkey, err := wallet.derivePrivateKey(path)

	// if err != nil {
	// 	log.Panic("failed to derivePrivateKey:", err)
	// }

	// fmt.Println(*pkey)

	/*

	   func NewTransaction(nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) *Transaction {
	   	return newTransaction(nonce, &to, amount, gasLimit, gasPrice, data)
	   }
	*/
	//创建一笔交易
	nonce := uint64(0)
	amount := big.NewInt(3000000000000000000)
	toAccount := common.HexToAddress("0x318490047EF28b8a04F457dd5449FA76f3f5bC0a")
	gasLimit := uint64(21000)
	gasPrice := big.NewInt(21000000000)
	var data []byte
	tx := types.NewTransaction(nonce, toAccount, amount, gasLimit, gasPrice, data)
	//签名
	fmt.Println(account.Address.Hex())
	stx, err := wallet.SignTx(account, tx, nil)
	if err != nil {
		log.Panic("failed to SignTx:", err)
	}
	//发送
	cli, err := ethclient.Dial("HTTP://127.0.0.1:7545")
	if err != nil {
		log.Panic("failed to Dial:", err)
	}
	defer cli.Close()

	err = cli.SendTransaction(context.Background(), stx)
	if err != nil {
		log.Panic("failed to SendTransaction:", err)
	}

}

func test_keystore() {
	nm := "august human human affair mechanic night verb metal embark marine orient million"

	// 助记词转化为种子 - > 账户地址
	// 先推导路径，再获得钱包
	path := MustParseDerivationPath("m/44'/60'/0'/0/0")

	wallet, err := NewFromMnemonic(nm, "")

	if err != nil {
		log.Panic("failed to NewFromMnemonic:", err)
	}

	account, err := wallet.Derive(path, false)
	if err != nil {
		log.Panic("failed to Derive:", err)
	}

	fmt.Println(account.Address.Hex())
	//得到私钥
	pkey, err := wallet.derivePrivateKey(path)

	if err != nil {
		log.Panic("failed to derivePrivateKey:", err)
	}

	fmt.Println(*pkey)

	key := NewKeyFromECDSA(pkey)

	hdks := NewHDKeyStore("./data")

	err = hdks.StoreKey(hdks.JoinPath(account.Address.Hex()), key, "123")
	if err != nil {
		log.Panic("failed to StoreKey:", err)
	}

}

func test_sendTransaction() {
	cli, err := ethclient.Dial("HTTP://127.0.0.1:7545") //注意地址变化 8545
	if err != nil {
		log.Panic(err)
	}

	defer cli.Close()

	mnemonic := "august human human affair mechanic night verb metal embark marine orient million"
	wallet, err := NewFromMnemonic(mnemonic, "")
	if err != nil {
		log.Fatal(err)
	}

	path := MustParseDerivationPath("m/44'/60'/0'/0/1") //第2个账户地址
	account, err := wallet.Derive(path, true)
	if err != nil {
		log.Fatal(err)
	}

	nonce := uint64(2)
	value := big.NewInt(5000000000000000000)
	toAddress := common.HexToAddress("0x44f4CD617655104649C1b866D20D5EAE198deD38")
	gasLimit := uint64(21000)
	gasPrice := big.NewInt(21000000000)
	var data []byte

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
	signedTx, err := wallet.SignTx(account, tx, nil)
	if err != nil {
		log.Fatal("failed to signed tx:", err)
	}

	err = cli.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Panic("failed to send transaction:", err)
	}
	//toAddress := common.HexToAddress("0x44f4CD617655104649C1b866D20D5EAE198deD38")
	bl, err := cli.BalanceAt(context.Background(), toAddress, big.NewInt(1))
	fmt.Println(bl, bl.Uint64(), err)
}

func main() {
	//test_mnemonic()
	//test_ganache()
	test_transfer()
	//test_sendTransaction()
	//test_keystore()
}
