package main

import (
	"bufio"
	"fmt"
	uc "github.com/unicorn-engine/unicorn/bindings/go/unicorn"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var fn = "../mipigo/test/test.bin"
var dat, _ = ioutil.ReadFile(fn)

var f, _ = os.OpenFile("unicode", os.O_CREATE|os.O_RDWR, 0666)

func ensure(err error) {
	if err != nil {
		panic(err)
	}
}
func TestMIPSCode(t *testing.T) {
	defer f.Close()
	mu, err := uc.NewUnicorn(uc.ARCH_MIPS, uc.MODE_32|uc.MODE_BIG_ENDIAN)
	ensure(err)
	addTestHook(mu)
	ensure(mu.MemMap(0, 0x80000000))
	mu.MemWrite(0, dat)

	ensure(mu.Start(0, 0x5ead0004))
}

var buffer = bufio.NewWriter(f)

func addTestHook(mu uc.Unicorn) {
	mu.HookAdd(uc.HOOK_INTR, func(mu uc.Unicorn, intno uint32) {
		if intno != 17 {
			log.Fatal("invalid interrupt ", intno)
		}
		syscall_no, _ := mu.RegRead(uc.MIPS_REG_V0)
		if syscall_no == 4246 {
			// exit group
			log.Fatal("exist group")
		}
	}, 0, 0)

	mu.HookAdd(uc.HOOK_CODE, func(mu uc.Unicorn, addr uint64, size uint32) {
		pc, _ := mu.RegRead(uc.MIPS_REG_PC)
		insn, _ := mu.MemRead(pc, 4)
		b1 := bytetoBit(insn[0])
		op := b1[:6]
		var _b [4][8]uint8
		for i, b := range insn {
			_b[i] = bytetoBit(b)
		}
		buffer.WriteString(fmt.Sprintf("Syscall: op: %v insn:%v\n", op, _b))
		buffer.Flush()
	}, 0, 0x80000000)
}

func write() {
	defer f.Close()
	buffer := bufio.NewWriter(f)
	for i := 0; i < len(dat); i += 4 {

		b1 := bytetoBit(dat[i])
		op := b1[:6]
		buffer.WriteString(fmt.Sprintf("%6x: op:%v -:%b\n", i, op, dat[i:i+4]))
	}
	buffer.Flush()
}

func bytetoBit(b byte) (ret [8]uint8) {
	for i := 7; i >= 0; i-- {
		ret[i] = b & 1
		b >>= 1
	}
	return ret
}
