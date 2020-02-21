package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
	"time"
)

//Creamos la estructura del bloque
type Block struct {
	//Ingresamos la fecha
	Timestamp int64
	//Ingresamos los datos
	Data []byte
	//Obtenemos el hash del bloque anterior
	PrevBlockHash []byte
	//Declaramos el hash del bloque
	Hash []byte
	//Dificultad del bloque
	Nonce int
}

//Creamos la funcion de nuevo bloque
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	//Hacemos hash a todos los datos del bloque
	block.Hash = hash[:]
	block.Nonce = nonce
	//Retornamos el bloque
	return block
}

//Declaramos la estructura de la cadena de bloques
type Blockchain struct {
	//La declaramos como un arreglo con un puntero
	blocks []*Block
}

//Creamos la funcion de crear bloque
func (bc *Blockchain) AddBlock(data string) {
	//Revisamos el valor del bloque anterior
	prevBlock := bc.blocks[len(bc.blocks)-1]
	//Enviamos los datos y mandamos a llamar la funcion de crear bloque
	newBlock := NewBlock(data, prevBlock.Hash)
	//Agreamos el bloque al arreglo
	bc.blocks = append(bc.blocks, newBlock)
}

//Funcion para agregar el bloque genesis
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

//Aqui comenzamos con la creacion de la cadena de bloques
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

//Agregamos la dificultad para la creacion del hash
const targetBits = 24

//Declaramos el limite
var (
	maxNonce = math.MaxInt64
)

//Hacemos una validacion del trabajo
type ProofOfWork struct {
	//Apuntamos hacia el bloque
	block *Block
	//Declaramos el arranque
	target *big.Int
}

//Hacemos una nueva validacion del trabajo
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}

	return pow
}

//Funcion para convertir un entero en formato hexadecimal
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

//Funcion de la validacion del trabajo
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

//Algoritmo de prueba de trabajo
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Minando el contenido del bloque \"%s\"\n", pow.block.Data)
	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}
func main() {
	bc := NewBlockchain()

	bc.AddBlock("Juan Pablo de Jesus Figueroa Jaramillo")
	bc.AddBlock("Karen Valeria Ramirez Perez")

	for _, block := range bc.blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
	}
}
