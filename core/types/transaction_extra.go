package types

import (
	"bufio"
	"encoding/base64"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
	"os"
)

type TxExtra struct {
	TxHash        common.Hash
	PreStateRoot  common.Hash
	PostStateRoot common.Hash

	PreState map[common.Address]*StateAccount

	PreStorage map[common.Address]map[common.Hash]common.Hash

	PreCode map[common.Address][]byte
}

func NewTxExtra(hash common.Hash) *TxExtra {
	return &TxExtra{
		TxHash:        hash,
		PreStateRoot:  common.Hash{},
		PostStateRoot: common.Hash{},
		PreState:      map[common.Address]*StateAccount{},
		PreStorage:    map[common.Address]map[common.Hash]common.Hash{},
		PreCode:       map[common.Address][]byte{},
	}
}

func (t *TxExtra) SetPreStateRoot(hash common.Hash) {
	t.PreStateRoot = hash
}

func (t *TxExtra) SetPostStateRoot(hash common.Hash) {
	t.PostStateRoot = hash
}

func (t *TxExtra) AddPreState(address common.Address, stateAccount *StateAccount) {
	if t.PreState == nil {
		t.PreState = map[common.Address]*StateAccount{}
	}
	t.PreState[address] = stateAccount
}

func (t *TxExtra) AddPreStorage(address common.Address, key, value common.Hash) {
	if t.PreStorage[address] == nil {
		t.PreStorage[address] = map[common.Hash]common.Hash{}
	}
	t.PreStorage[address][key] = value
}

func (t *TxExtra) AddPreCode(address common.Address, enc []byte) {
	if t.PreCode == nil {
		t.PreCode = map[common.Address][]byte{}
	}
	t.PreCode[address] = enc
}

// MinerExtra File

func InitFile(blockNumber *big.Int) {
	var f *os.File
	var err error
	left := new(big.Int).Quo(blockNumber, new(big.Int).SetUint64(1000))
	os.MkdirAll("./minerExtra/"+left.String(), 0755)
	fileName := "./minerExtra/" + left.String() + "/" + blockNumber.String() + ".txt"
	if CheckFileExist(fileName) { //文件存在
		os.Remove(fileName)
	}
	f, err = os.Create(fileName) //创建文件
	if err != nil {
		log.Info("创建文件出错", "file create fail", err)
		return
	}
	defer f.Close()

}

func InitTxFile(blockNumber *big.Int) {
	var f *os.File
	var err error
	left := new(big.Int).Quo(blockNumber, new(big.Int).SetUint64(1000))
	os.MkdirAll("./minerExtra/temp/"+left.String(), 0755)
	os.MkdirAll("./minerExtra/temp/C"+left.String(), 0755)
	fileName := "./minerExtra/temp/" + left.String() + "/" + blockNumber.String() + ".txt"
	fileName2 := "./minerExtra/temp/C" + left.String() + "/" + blockNumber.String() + ".txt"
	if CheckFileExist(fileName) { //文件存在
		os.Remove(fileName)
	}
	if CheckFileExist(fileName2) { //文件存在
		os.Remove(fileName2)
	}
	f, err = os.Create(fileName) //创建文件
	if err != nil {
		log.Info("InitTxFile出错", "file create fail", err)
		return
	}
	defer f.Close()
}

func ReNameTxFile(blockNumber *big.Int) {
	left := new(big.Int).Quo(blockNumber, new(big.Int).SetUint64(1000))
	oldName := "./minerExtra/temp/" + left.String() + "/" + blockNumber.String() + ".txt"
	newName := "./minerExtra/temp/C" + left.String() + "/" + blockNumber.String() + ".txt"
	err := os.Rename(oldName, newName)
	if err != nil {
		log.Info("ReNameTxFileErr", "error", err)
	}
}

func DelTxFile(blockNumber *big.Int) {
	left := new(big.Int).Quo(blockNumber, new(big.Int).SetUint64(1000))
	fileName := "./minerExtra/temp/" + left.String() + "/" + blockNumber.String() + ".txt"
	fileName2 := "./minerExtra/temp/C" + left.String() + "/" + blockNumber.String() + ".txt"
	if CheckFileExist(fileName) { //文件存在
		os.Remove(fileName)
	}
	if CheckFileExist(fileName2) { //文件存在
		os.Remove(fileName2)
	}
}

func WriteHash(blockNumber *big.Int, hash common.Hash) {
	WriteTxFile(blockNumber, hash.String())
}

func WriteStateKey(blockNumber *big.Int, address common.Address) {
	WriteTxFile(blockNumber, "0\t"+address.String())
}

func WriteStorageKey(blockNumber *big.Int, address common.Address, key common.Hash) {
	WriteTxFile(blockNumber, "1\t"+address.String()+"\t"+key.String())
}
func WriteTxCode(blockNumber *big.Int, address common.Address) {
	WriteTxFile(blockNumber, "2\t"+address.String())
}

func WriteTxFile(blockNumber *big.Int, text string) {
	left := new(big.Int).Quo(blockNumber, new(big.Int).SetUint64(1000))
	filePath := "./minerExtra/temp/" + left.String() + "/" + blockNumber.String() + ".txt"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		//log.Info("文件打开失败", "文件打开失败", err, "string", text)
	}
	//及时关闭file句柄
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	write.WriteString(text + "\n")
	//Flush将缓存的文件真正写入到文件中
	write.Flush()
}

func ReadTxFile(blockNumber *big.Int, hash common.Hash) []string {
	var lines []string
	left := new(big.Int).Quo(blockNumber, new(big.Int).SetUint64(1000))
	myfile, err := os.Open("./minerExtra/temp/C" + left.String() + "/" + blockNumber.String() + ".txt") //open the file
	if err != nil {
		log.Info("ReadTxFile", "Error opening file:", err)
		return lines
	}
	defer myfile.Close()
	scanner := bufio.NewScanner(myfile)
	start := false
	for scanner.Scan() {
		if start {
			if scanner.Text() == hash.String() {
				break
			}
			lines = append(lines, scanner.Text())
		}
		if scanner.Text() == hash.String() {
			start = true
		}
	}
	return lines
}

func WriteTxHash(blockNumber *big.Int, hash common.Hash) {
	WriteFile(blockNumber, "txHash\t"+hash.String())
}

func WritePreRoot(blockNumber *big.Int, hash common.Hash) {
	WriteFile(blockNumber, "preStateRoot\t"+hash.String())
}

func WritePostRoot(blockNumber *big.Int, hash common.Hash) {
	WriteFile(blockNumber, "postStateRoot\t"+hash.String())
}

func WritePreState(blockNumber *big.Int, address common.Address, enc []byte) {
	WriteFile(blockNumber, "preState\t"+address.String()+"\t"+base64.StdEncoding.EncodeToString(enc))
}

func WritePostState(blockNumber *big.Int, address common.Address, enc []byte) {
	WriteFile(blockNumber, "postState\t"+address.String()+"\t"+base64.StdEncoding.EncodeToString(enc))
}

func WritePreStateProof(blockNumber *big.Int, address common.Address, path [][]byte) {
	text := ""
	for _, bytes := range path {
		text += "\t" + base64.StdEncoding.EncodeToString(bytes)
	}

	WriteFile(blockNumber, "preStateProof\t"+address.String()+text)
}

func WritePostStateProof(blockNumber *big.Int, address common.Address, path [][]byte) {
	text := ""
	for _, bytes := range path {
		text += "\t" + base64.StdEncoding.EncodeToString(bytes)
	}

	WriteFile(blockNumber, "postStateProof\t"+address.String()+text)
}

func WritePreStorage(blockNumber *big.Int, address common.Address, hash, val common.Hash) {
	WriteFile(blockNumber, "preStorage\t"+address.String()+"\t"+hash.String()+"\t"+val.String())
}

func WritePostStorage(blockNumber *big.Int, address common.Address, hash, val common.Hash) {
	WriteFile(blockNumber, "postStorage\t"+address.String()+"\t"+hash.String()+"\t"+val.String())
}

func WritePreStorageProof(blockNumber *big.Int, address common.Address, hash common.Hash, path [][]byte) {
	text := ""
	for _, bytes := range path {
		text += "\t" + base64.StdEncoding.EncodeToString(bytes)
	}

	WriteFile(blockNumber, "preStorageProof\t"+address.String()+"\t"+hash.String()+text)
}

func WritePostStorageProof(blockNumber *big.Int, address common.Address, hash common.Hash, path [][]byte) {
	text := ""
	for _, bytes := range path {
		text += "\t" + base64.StdEncoding.EncodeToString(bytes)
	}

	WriteFile(blockNumber, "postStorageProof\t"+address.String()+"\t"+hash.String()+text)
}

func WritePreCode(blockNumber *big.Int, address common.Address, code []byte) {
	WriteFile(blockNumber, "preCode\t"+address.String()+"\t"+base64.StdEncoding.EncodeToString(code))
}

func WritePostCode(blockNumber *big.Int, address common.Address, code []byte) {
	WriteFile(blockNumber, "postCode\t"+address.String()+"\t"+base64.StdEncoding.EncodeToString(code))
}

func WriteFile(blockNumber *big.Int, text string) {
	left := new(big.Int).Quo(blockNumber, new(big.Int).SetUint64(1000))
	filePath := "./minerExtra/" + left.String() + "/" + blockNumber.String() + ".txt"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		//log.Info("文件打开失败", "文件打开失败", err, "string", text)
	}
	//及时关闭file句柄
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	write.WriteString(text + "\n")
	//Flush将缓存的文件真正写入到文件中
	write.Flush()
}

func CheckFileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
