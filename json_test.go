package main

import (
	"fmt"
	"log"
	"testing"
)

const nMax = 51

type mahasiswa struct {
	NIM   string
	Nama  string
	Nilai int
}
type arrayMahasiswa [nMax]mahasiswa

func nilaiPertama(T arrayMahasiswa, N int, NIM string) int {
	for i := 0; i < N; i++ {
		if T[i].NIM == NIM {
			return T[i].Nilai
		}
	}
	return -1
}

func TestJson(t *testing.T) {
	var data arrayMahasiswa
	var n int

	_, err := fmt.Scan(&n)
	if err != nil {
		t.Fatalf("Failed to scan the number of students: %v", err)
	}

	if n > nMax {
		log.Fatalf("Number of students exceeds maximum allowed: %d", nMax)
	}

	for i := 0; i < n; i++ {
		fmt.Scan(&data[i].NIM, &data[i].Nama, &data[i].Nilai)
	}

	nilai := nilaiPertama(data, n, "113")

	fmt.Println(nilai)

}
