package ipso_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/aliakseiz/ipso"
)

const benchRegistryFile = "share/registry.yaml"

func BenchmarkReg_Import(b *testing.B) {
	data := readBenchFile(b)

	b.ReportAllocs()
	b.SetBytes(int64(len(data)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reg := &ipso.Reg{}

		err := reg.Import(bytes.NewReader(data))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReg_Export(b *testing.B) {
	reg := newBenchReg(b)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := reg.Export(io.Discard)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReg_Sanitize(b *testing.B) {
	reg := newBenchReg(b)
	pristine := newBenchReg(b)
	sanitizer := ipso.DefaultSanitizer()

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		reg.Objects = pristine.GetObjects()

		b.StartTimer()

		reg.Sanitize(sanitizer)
	}
}

func BenchmarkReg_GetObjects(b *testing.B) {
	reg := newBenchReg(b)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = reg.GetObjects()
	}
}

func BenchmarkReg_FindObjectByIDAndVer(b *testing.B) {
	reg := newBenchReg(b)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := reg.FindObjectByIDAndVer(3303, "1.0")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReg_FindObjectsByID(b *testing.B) {
	reg := newBenchReg(b)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := reg.FindObjectsByID(3303)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReg_FindResourceByOIR(b *testing.B) {
	reg := newBenchReg(b)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := reg.FindResourceByOIR("3303/0/5700")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReg_FindResourcesByID(b *testing.B) {
	reg := newBenchReg(b)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := reg.FindResourcesByID(5700)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReg_FindResourcesByObjResIDs(b *testing.B) {
	reg := newBenchReg(b)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := reg.FindResourcesByObjResIDs(3303, 5700)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReg_FindObjectsByName(b *testing.B) {
	reg := newBenchReg(b)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := reg.FindObjectsByName("temperature")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReg_FindResourcesByName(b *testing.B) {
	reg := newBenchReg(b)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := reg.FindResourcesByName("Sensor Value")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReg_Compare(b *testing.B) {
	reg := newBenchReg(b)
	other := newBenchReg(b)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = reg.Compare(other)
	}
}

func readBenchFile(b *testing.B) []byte {
	b.Helper()

	data, err := os.ReadFile(benchRegistryFile)
	if err != nil {
		b.Fatal(err)
	}

	return data
}

func newBenchReg(b *testing.B) *ipso.Reg {
	b.Helper()

	reg := &ipso.Reg{}

	err := reg.Import(bytes.NewReader(readBenchFile(b)))
	if err != nil {
		b.Fatal(err)
	}

	return reg
}
