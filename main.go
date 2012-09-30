package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"log"
	"os"
	"runtime"
	"time"
)

const Duration1 time.Duration = time.Second
const Duration2 time.Duration = 3 * time.Second
const Duration3 time.Duration = time.Minute + time.Second

const DataSize = 48

type Data [DataSize]byte

func NewData() *Data {
	d := Data{}
	rand.Read(d[:])
	return &d
}

const ShaSize = sha256.Size

type Sha [ShaSize]byte

func (data *Data) Sha() *Sha {
	hash := sha256.New()
	hash.Write(data[:])
	sha := Sha{}
	hash.Sum(sha[:0])
	return &sha
}

func (sha *Sha) String() string {
	return hex.EncodeToString(sha[:])
}

const BufferSize = ShaSize + DataSize

type Buffer struct {
	*Sha
	*Data
}

func NewBuffer() Buffer {
	data := NewData()
	return Buffer{data.Sha(), data}
}

func main() {

	cpu := flag.Int("c", runtime.GOMAXPROCS(-1), "Number of goroutines")
	mem := flag.Int("m", 128*1024*1024, "Memory to use in bytes")
	flag.Parse()

	output := log.New(os.Stdout, "", log.LstdFlags)
	errput := log.New(os.Stderr, "", log.LstdFlags)

	if *cpu < 1 {
		*cpu = 1
	}
	errput.Println("cpu =", *cpu)

	buffers := make(chan Buffer, *cpu*4)
	for i := 0; i < *cpu; i++ {
		go func() {
			for {
				buffers <- NewBuffer()
			}
		}()
	}

	if *mem < BufferSize {
		*mem = BufferSize
	}
	errput.Println("mem =", *mem)

	errput.Println("size of sha  =", ShaSize)
	errput.Println("size of data =", DataSize)

	max, maxreached := *mem/BufferSize, false
	shas, comparisons := make(map[Sha]*Data, max), uint64(0)
	next := time.Now().Add(Duration1)
	for {
		if maxreached && time.Now().After(next) {
			errput.Println("total comparisons =", comparisons)
			next = next.Add(Duration3)
		}
		buffer := <-buffers
		sha := buffer.Sha
		if data, ok := shas[*sha]; ok {
			if bytes.Equal(buffer.Data[:], data[:]) {
				errput.Println("data collision", sha, data[:])
			} else {
				output.Println("sha collision", sha, buffer.Data[:], data[:])
			}
		} else if len(shas) < max {
			shas[*sha] = buffer.Data
			if !maxreached && (len(shas) == max || time.Now().After(next)) {
				errput.Printf("num/max buffers = %d/%d (%0.2f%%)", len(shas), max, float32(len(shas))*100/float32(max))
				if len(shas) < max {
					next = next.Add(Duration2)
				} else {
					maxreached, next = true, next.Add(Duration3)
				}
			}
		}
		comparisons++
	}
}
