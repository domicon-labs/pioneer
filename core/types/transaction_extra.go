package types

import (
	"bufio"
	"encoding/base64"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"os"
	"strings"
)

type AddressBytesList struct {
	List []AddressBytes
}
type AddressBytes struct {
	Address common.Address
	Bytes   []byte
}

type StateProofList struct {
	List []StateProof
}
type StateProof struct {
	Address common.Address
	Path    [][]byte
}

type AddressStorageList struct {
	List []AddressStorage
}
type AddressStorage struct {
	Address common.Address
	Storage StorageList
}
type StorageList struct {
	List []Storage
}
type Storage struct {
	Key common.Hash
	Val common.Hash
}
type AddressStorageProofList struct {
	List []AddressStorageProof
}
type AddressStorageProof struct {
	Address      common.Address
	StorageProof StorageProofList
}
type StorageProofList struct {
	List []StorageProof
}
type StorageProof struct {
	Hash common.Hash
	Path [][]byte
}
type ExtraData struct {
	TxHash        common.Hash
	PreStateRoot  common.Hash
	PostStateRoot common.Hash
	Codes         AddressBytesList

	PreStateData    AddressBytesList
	PreStateProof   StateProofList
	PreStorageData  AddressStorageList
	PreStorageProof AddressStorageProofList

	PostStateData    AddressBytesList
	PostStateProof   StateProofList
	PostStorageData  AddressStorageList
	PostStorageProof AddressStorageProofList
}

func NewExtraData() *ExtraData {
	return &ExtraData{
		TxHash:           common.Hash{},
		PreStateRoot:     common.Hash{},
		PostStateRoot:    common.Hash{},
		Codes:            AddressBytesList{},
		PreStateData:     AddressBytesList{},
		PreStateProof:    StateProofList{},
		PreStorageData:   AddressStorageList{},
		PreStorageProof:  AddressStorageProofList{},
		PostStateData:    AddressBytesList{},
		PostStateProof:   StateProofList{},
		PostStorageData:  AddressStorageList{},
		PostStorageProof: AddressStorageProofList{},
	}
}
func (e *ExtraData) SetTxHash(hash common.Hash) {
	e.TxHash = hash
}

func (e *ExtraData) SetPreStateRoot(hash common.Hash) {
	e.PreStateRoot = hash
}

func (e *ExtraData) GetPreStateRoot() common.Hash {
	return e.PreStateRoot
}

func (e *ExtraData) SetPostStateRoot(hash common.Hash) {
	e.PostStateRoot = hash
}

func (e *ExtraData) GetPostStateRoot() common.Hash {
	return e.PostStateRoot
}

func (e *ExtraData) AddCode(address common.Address, code []byte) {
	if _, exist := e.Codes.Get(address); exist {
		return
	}
	e.Codes.Add(AddressBytes{
		Address: address,
		Bytes:   code,
	})
}
func (e *ExtraData) Code(address common.Address) []byte {
	if obj, exist := e.Codes.Get(address); exist {
		return obj.Bytes
	}
	return nil
}

func (e *ExtraData) AddPreState(address common.Address, enc []byte) {
	if _, exist := e.PreStateData.Get(address); exist {
		return
	}
	e.PreStateData.Add(AddressBytes{
		Address: address,
		Bytes:   enc,
	})
}

func (e *ExtraData) AllPreState() AddressBytesList {
	return e.PreStateData
}

func (e *ExtraData) GetPreState(address common.Address) []byte {
	if obj, exist := e.PreStateData.Get(address); exist {
		return obj.Bytes
	}
	return nil
}

func (e *ExtraData) AddPreStateProof(address common.Address, path [][]byte) {
	if _, exist := e.PreStateProof.Get(address); exist {
		return
	}
	e.PreStateProof.Add(StateProof{
		Address: address,
		Path:    path,
	})
}

func (e *ExtraData) AllPreStateProof() StateProofList {
	return e.PreStateProof
}

func (e *ExtraData) GetPreStateProof(address common.Address) [][]byte {
	if obj, exist := e.PreStateProof.Get(address); exist {
		return obj.Path
	}
	return nil
}

func (e *ExtraData) AddPreStorage(address common.Address, hash common.Hash, val common.Hash) {
	//log.Info("AddPreStorage", "address", address, "hash", hash, "val", val)
	if obj, exist := e.PreStorageData.Get(address); !exist {
		temp := AddressStorage{
			Address: address,
			Storage: StorageList{},
		}
		temp.Storage.Add(Storage{
			Key: hash,
			Val: val,
		})
		e.PreStorageData.Add(temp)
	} else {
		if _, ex := obj.Storage.Get(hash); !ex {
			obj.Storage.Add(Storage{
				Key: hash,
				Val: val,
			})
			//log.Info("else", "storage", obj.Storage)
		}
	}
}

func (e *ExtraData) GetPreStorage(address common.Address, hash common.Hash) common.Hash {
	if obj, exist := e.PreStorageData.Get(address); exist {
		if storage, ex := obj.Storage.Get(hash); ex {
			return storage
		}
	}
	return common.Hash{}
}

func (e *ExtraData) AddPreStorageProof(address common.Address, hash common.Hash, path [][]byte) {
	if obj, exist := e.PreStorageProof.Get(address); !exist {
		temp := AddressStorageProof{
			Address:      address,
			StorageProof: StorageProofList{},
		}
		temp.StorageProof.Add(StorageProof{
			Hash: hash,
			Path: path,
		})
		e.PreStorageProof.Add(temp)
	} else {
		if _, ex := obj.StorageProof.Get(hash); !ex {
			obj.StorageProof.Add(StorageProof{
				Hash: hash,
				Path: path,
			})
		}
	}
}

func (e *ExtraData) GetPreStorageProof(address common.Address, hash common.Hash) [][]byte {
	if obj, exist := e.PreStorageProof.Get(address); exist {
		if proof, ex := obj.StorageProof.Get(hash); ex {
			return proof.Path
		}
	}
	return nil
}
func (e *ExtraData) AddPostState(address common.Address, enc []byte) {
	if obj, exist := e.PostStateData.Get(address); exist {
		obj.Bytes = enc
		return
	}
	e.PostStateData.Add(AddressBytes{
		Address: address,
		Bytes:   enc,
	})
}

func (e *ExtraData) AllPostState() AddressBytesList {
	return e.PostStateData
}

func (e *ExtraData) GetPostState(address common.Address) []byte {
	if obj, exist := e.PostStateData.Get(address); exist {
		return obj.Bytes
	}
	return nil
}

func (e *ExtraData) AddPostStateProof(address common.Address, path [][]byte) {
	if obj, exist := e.PostStateProof.Get(address); exist {
		obj.Path = path
		return
	}
	e.PostStateProof.Add(StateProof{
		Address: address,
		Path:    path,
	})
}

func (e *ExtraData) AllPostStateProof() StateProofList {
	return e.PostStateProof
}

func (e *ExtraData) GetPostStateProof(address common.Address) [][]byte {
	if obj, exist := e.PostStateProof.Get(address); exist {
		return obj.Path
	}
	return nil
}

func (e *ExtraData) AddPostStorage(address common.Address, hash common.Hash, val common.Hash) {
	if obj, exist := e.PostStorageData.Get(address); !exist {
		temp := AddressStorage{
			Address: address,
			Storage: StorageList{},
		}
		temp.Storage.Add(Storage{
			Key: hash,
			Val: val,
		})
		e.PostStorageData.Add(temp)
	} else {
		if _, ex := obj.Storage.Get(hash); !ex {
			obj.Storage.Add(Storage{
				Key: hash,
				Val: val,
			})
		}
	}
}

func (e *ExtraData) GetPostStorage(address common.Address, hash common.Hash) common.Hash {
	if obj, exist := e.PostStorageData.Get(address); exist {
		if storage, ex := obj.Storage.Get(hash); ex {
			return storage
		}
	}
	return common.Hash{}
}

func (e *ExtraData) AddPostStorageProof(address common.Address, hash common.Hash, path [][]byte) {
	if obj, exist := e.PostStorageProof.Get(address); !exist {
		temp := AddressStorageProof{
			Address:      address,
			StorageProof: StorageProofList{},
		}
		temp.StorageProof.Add(StorageProof{
			Hash: hash,
			Path: path,
		})
		e.PostStorageProof.Add(temp)
	} else {
		if _, ex := obj.StorageProof.Get(hash); !ex {
			obj.StorageProof.Add(StorageProof{
				Hash: hash,
				Path: path,
			})
		}
	}
}

func (e *ExtraData) GetPostStorageProof(address common.Address, hash common.Hash) [][]byte {
	if obj, exist := e.PostStorageProof.Get(address); exist {
		if proof, ex := obj.StorageProof.Get(hash); ex {
			return proof.Path
		}
	}
	return nil
}

func (l *StorageList) Add(obj Storage) {
	l.List = append(l.List, obj)
}

func (l *StorageList) Get(hash common.Hash) (common.Hash, bool) {
	for _, obj := range l.List {
		if obj.Key == hash {
			return obj.Val, true
		}
	}
	return common.Hash{}, false
}

func (l *AddressBytesList) Add(obj AddressBytes) {
	l.List = append(l.List, obj)
}

func (l *AddressBytesList) Get(address common.Address) (*AddressBytes, bool) {
	for _, obj := range l.List {
		if obj.Address == address {
			return &obj, true
		}
	}
	return nil, false
}
func (s *StorageProofList) Add(obj StorageProof) {
	s.List = append(s.List, obj)
}

func (s *StorageProofList) Get(hash common.Hash) (*StorageProof, bool) {
	for _, obj := range s.List {
		if obj.Hash == hash {
			return &obj, true
		}
	}
	return nil, false
}

func (s *AddressStorageProofList) Add(obj AddressStorageProof) {
	s.List = append(s.List, obj)
}

func (s *AddressStorageProofList) Get(address common.Address) (*AddressStorageProof, bool) {
	for _, obj := range s.List {
		if obj.Address == address {
			return &obj, true
		}
	}
	return nil, false
}

func (s *StateProofList) Add(obj StateProof) {
	s.List = append(s.List, obj)
}

func (s *StateProofList) Get(address common.Address) (*StateProof, bool) {
	for _, obj := range s.List {
		if obj.Address == address {
			return &obj, true
		}
	}
	return nil, false
}

func (s *AddressStorageList) Add(obj AddressStorage) {
	s.List = append(s.List, obj)
}
func (s *AddressStorageList) Get(address common.Address) (*AddressStorage, bool) {
	for _, obj := range s.List {
		if obj.Address == address {
			return &obj, true
		}
	}
	return nil, false
}

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

func ReadFile(blockNumber *big.Int, db ethdb.KeyValueStore) {
	myfile, err := os.Open("./minerExtra/" + blockNumber.String() + ".txt") //open the file
	if err != nil {
		log.Info("ReadFile", "Error opening file:", err)
		return
	}
	defer myfile.Close()
	extraData := NewExtraData()
	scanner := bufio.NewScanner(myfile) //scan the contents of a file and print line by line
	for scanner.Scan() {
		line := scanner.Text()
		lArr := strings.Split(line, "\t")
		stype := lArr[0]
		switch stype {
		case "txHash":
			extraData.SetTxHash(common.HexToHash(lArr[1]))
		case "preStateRoot":
			extraData.SetPreStateRoot(common.HexToHash(lArr[1]))
		case "postStateRoot":
			extraData.SetPostStateRoot(common.HexToHash(lArr[1]))
			data, err0 := rlp.EncodeToBytes(extraData)
			if err0 != nil {
				log.Error("Failed to encode txExtra", "err", err0)
			}
			db.Put(extraData.TxHash.Bytes(), data)
			//log.Info("最终数据", "最终数据", extraData)
			extraData = NewExtraData()
		case "code":
			if c, err := base64.StdEncoding.DecodeString(lArr[2]); err == nil {
				extraData.AddCode(common.HexToAddress(lArr[1]), c)
			}
		case "preState":
			if c, err := base64.StdEncoding.DecodeString(lArr[2]); err == nil {
				extraData.AddPreState(common.HexToAddress(lArr[1]), c)
			}
		case "postState":
			if c, err := base64.StdEncoding.DecodeString(lArr[2]); err == nil {
				extraData.AddPostState(common.HexToAddress(lArr[1]), c)
			}
		case "preStateProof":
			var path [][]byte
			for _, s := range lArr[2:] {
				if c, err := base64.StdEncoding.DecodeString(s); err == nil {
					path = append(path, c)
				}
			}
			extraData.AddPreStateProof(common.HexToAddress(lArr[1]), path)
		case "postStateProof":
			var path [][]byte
			for _, s := range lArr[2:] {
				if c, err := base64.StdEncoding.DecodeString(s); err == nil {
					path = append(path, c)
				}
			}
			extraData.AddPostStateProof(common.HexToAddress(lArr[1]), path)
		case "preStorage":
			//log.Info("preStorage", "address", lArr[1], "hash", lArr[2], "val", lArr[3])
			extraData.AddPreStorage(common.HexToAddress(lArr[1]), common.HexToHash(lArr[2]), common.HexToHash(lArr[3]))
			//log.Info("preStorage", "extra.PreStorage", extraData.PreStorageData.List)
		case "postStorage":
			//log.Info("postStorage", "address", lArr[1], "hash", lArr[2], "val", lArr[3])
			extraData.AddPostStorage(common.HexToAddress(lArr[1]), common.HexToHash(lArr[2]), common.HexToHash(lArr[3]))
		case "preStorageProof":
			var path [][]byte
			for _, s := range lArr[3:] {
				if c, err := base64.StdEncoding.DecodeString(s); err == nil {
					path = append(path, c)
				}
			}
			extraData.AddPreStorageProof(common.HexToAddress(lArr[1]), common.HexToHash(lArr[2]), path)
		case "postStorageProof":
			var path [][]byte
			for _, s := range lArr[3:] {
				if c, err := base64.StdEncoding.DecodeString(s); err == nil {
					path = append(path, c)
				}
			}
			extraData.AddPostStorageProof(common.HexToAddress(lArr[1]), common.HexToHash(lArr[2]), path)
		default:
		}
	}

	if err := scanner.Err(); err != nil {
		log.Info("end", "Error reading from file:", err) //print error if scanning is not done properly
	}
}

func CheckFileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

type ExtraProof struct {
	List map[common.Hash][]byte
}

func (e *ExtraProof) Has(key []byte) (bool, error) {
	if e.List[common.BytesToHash(key)] == nil {
		return false, nil
	}
	return true, nil
}

func (e *ExtraProof) Get(key []byte) ([]byte, error) {
	if data := e.List[common.BytesToHash(key)]; data != nil {
		return data, nil
	}
	return nil, nil
}

func NewExtraProof(b [][]byte) *ExtraProof {
	proof := map[common.Hash][]byte{}
	for _, bytes := range b {
		temp := crypto.Keccak256(bytes)
		proof[common.BytesToHash(temp)] = bytes
	}
	return &ExtraProof{proof}
}

// toHexSlice creates a slice of hex-strings based on []byte.
func ToHexSlice(b [][]byte) []string {
	r := make([]string, len(b))
	for i := range b {
		r[i] = hexutil.Encode(b[i])
	}
	return r
}

type TxExtra struct {
	TxHash        common.Hash
	PreStateRoot  common.Hash
	PostStateRoot common.Hash

	PreState  map[common.Address]StateAccount
	PostState map[common.Address][]byte

	PreStorage  map[common.Address]map[common.Hash]common.Hash
	PostStorage map[common.Address]map[common.Hash]common.Hash

	PreStateProof  map[common.Address][][]byte
	PostStateProof map[common.Address][][]byte

	PreStorageProof  map[common.Address]map[common.Hash][][]byte
	PostStorageProof map[common.Address]map[common.Hash][][]byte
}

func NewTxExtra(hash common.Hash) *TxExtra {
	return &TxExtra{
		TxHash:           hash,
		PreStateRoot:     common.Hash{},
		PostStateRoot:    common.Hash{},
		PreState:         map[common.Address]StateAccount{},
		PostState:        map[common.Address][]byte{},
		PreStorage:       map[common.Address]map[common.Hash]common.Hash{},
		PostStorage:      map[common.Address]map[common.Hash]common.Hash{},
		PreStateProof:    map[common.Address][][]byte{},
		PostStateProof:   map[common.Address][][]byte{},
		PreStorageProof:  map[common.Address]map[common.Hash][][]byte{},
		PostStorageProof: map[common.Address]map[common.Hash][][]byte{},
	}
}

func (t *TxExtra) AddPreState(address common.Address, stateAccount StateAccount) {
	if t.PreState == nil {
		t.PreState = map[common.Address]StateAccount{}
	}
	t.PreState[address] = stateAccount
}

func (t *TxExtra) AddPostState(address common.Address, enc []byte) {
	if t.PostState == nil {
		t.PostState = map[common.Address][]byte{}
	}
	t.PostState[address] = enc
}

func (t TxExtra) AddPreStorage(address common.Address, key, value common.Hash) {
	if t.PreStorage[address] == nil {
		t.PreStorage[address] = map[common.Hash]common.Hash{}
	}
	t.PreStorage[address][key] = value
}

func (t TxExtra) AddPostStorage(address common.Address, key, value common.Hash) {
	if t.PostStorage[address] == nil {
		t.PostStorage[address] = map[common.Hash]common.Hash{}
	}
	t.PostStorage[address][key] = value
}

func (t *TxExtra) AddPreStateProof(address common.Address, proof [][]byte) {
	t.PreStateProof[address] = proof
}

func (t *TxExtra) AddPostStateProof(address common.Address, proof [][]byte) {
	t.PostStateProof[address] = proof
}

func (t TxExtra) AddPreStorageProof(address common.Address, key common.Hash, proof [][]byte) {
	if t.PreStorageProof[address] == nil {
		t.PreStorageProof[address] = map[common.Hash][][]byte{}
	}
	t.PreStorageProof[address][key] = proof
}

func (t TxExtra) AddPostStorageProof(address common.Address, key common.Hash, proof [][]byte) {
	if t.PostStorageProof[address] == nil {
		t.PostStorageProof[address] = map[common.Hash][][]byte{}
	}
	t.PostStorageProof[address][key] = proof
}
