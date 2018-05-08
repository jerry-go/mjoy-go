package accountTest

import (
	"testing"
	"fmt"
	"mjoy.io/accounts/keystore"
	"mjoy.io/core/transaction"
	"math/big"
)

//test Account Read
func TestAccountRead(t *testing.T){
	fmt.Println("WellCome to CreateAccount")
	//init myKeyStore
	myKeyStore := keystore.NewKeyStore("./keystore",1 << 18,1)
	if myKeyStore == nil{
		panic("mykeystore== nil")
	}
	//read accounts and print
	acExists := myKeyStore.Accounts()
	for _,ac := range acExists{
		fmt.Printf("Address:%x,   Url:%s\n",ac.Address,ac.URL)
	}
}

//test Account Create
func TestAccountCreate(t *testing.T){
	fmt.Println("WellCome to CreateAccount")
	//init myKeyStore
	myKeyStore := keystore.NewKeyStore("./keystore",1 << 18,1)
	if myKeyStore == nil{
		panic("mykeystore== nil")
	}
	//read accounts and print
	acExists := myKeyStore.Accounts()
	for _,ac := range acExists{
		fmt.Printf("Address:%x,   Url:%s\n",ac.Address,ac.URL)
	}
	//create account
	ac,err := myKeyStore.NewAccount("123")
	if err != nil{
		panic(err)
	}
	fmt.Printf("Print NewAccount Address:%x,   Url:%s\n",ac.Address,ac.URL)
	//read accounts again and print
	acExists = acExists[:0]
	acExists = myKeyStore.Accounts()
	for _,ac := range acExists{
		fmt.Printf("After Address:%x,   Url:%s\n",ac.Address,ac.URL)
	}
}

//test Wallets
func TestWalletsShow(t *testing.T){
	fmt.Println("Wellcome to Wallets show.....")
	//init myKeyStore
	myKeyStore := keystore.NewKeyStore("./keystore",1 << 18,1)
	if myKeyStore == nil{
		panic("mykeystore== nil")
	}

	//check all the wallets
	aw := myKeyStore.Wallets()
	//show all wallets
	for _,w := range aw{
		fmt.Printf("Wallet URL:%s\n" , w.URL())
	}
}

//unlock All account
func TestUnlockAllAccount(t *testing.T){
	fmt.Println("Wellcome to Wallets show.....")
	//init myKeyStore
	myKeyStore := keystore.NewKeyStore("./keystore",1 << 18,1)
	if myKeyStore == nil{
		panic("mykeystore== nil")
	}


	//read accounts and print

	acExists := myKeyStore.Accounts()
	for _,ac := range acExists{
		fmt.Printf("After Address:%x,   Url:%s\n",ac.Address,ac.URL)
		//unlock
		err := myKeyStore.Unlock(ac,"123")
		if err != nil{
			fmt.Println("unlock Wrong....",err)
		}
	}
	//print accounts have been unlocked
	aw := myKeyStore.Wallets()
	//show all wallets
	for _,w := range aw{
		s,_ := w.Status()
		fmt.Printf("wallet: Url:%s,  Status:%s\n" , w.URL(),s)
	}

}

//test transaction signing
func TestTransactionSign(t *testing.T){
	fmt.Println("Wellcome to Wallets show.....")
	myKeyStore := keystore.NewKeyStore("./keystore",1 << 18,1)
	if myKeyStore == nil{
		panic("mykeystore== nil")
	}


	//read accounts and print

	acExists := myKeyStore.Accounts()
	for _,ac := range acExists{
		fmt.Printf("After Address:%x,   Url:%s\n",ac.Address,ac.URL)
		//unlock
		err := myKeyStore.Unlock(ac,"123")
		if err != nil{
			fmt.Println("unlock Wrong....",err)
		}
	}
	//print all accounts have been unlocked
	aw := myKeyStore.Wallets()
	//show all wallets
	for _,w := range aw{
		s,_ := w.Status()
		fmt.Printf("wallet: Url:%s,  Status:%s\n" , w.URL(),s)
	}

	fmt.Println("AccountExists len:" , len(acExists))
	if len(acExists) < 2 {
		t.Skip("Test Transaction Signing Should have 2 or more accounts,Please Run test : TestAccountCreate ")
	}
	amount:=big.NewInt(int64(20))

	res:=big.NewInt(int64(10))
	//create transaction,ac[0]--->ac[1]
	newTx:=transaction.NewTransaction(1,acExists[1].Address,amount,10,res,nil)

	newTx.PrintVSR()
	//use walet[0] to sign
	tx,err:=aw[0].SignTxWithPassphrase(acExists[0],"123",newTx,big.NewInt(int64(1)))
	if err != nil{
		fmt.Println("err=",err)
		panic("signTxWithPassphrase failed")
	}
	tx.PrintVSR()
	_ = tx
	fmt.Println("over.....")
}
