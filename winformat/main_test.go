package winformat

import (
	//"bytes"
	//"encoding/binary"
	. "github.com/r7kamura/gospel"
	"os"
	"testing"
)

func ReadBinaryFile(filename string) []byte {
	// ファイルを開く
	file, _ := os.Open(filename)
	// ファイルから1バイト読み出し
	b := make([]byte, 1400)
	file.Read(b)
	return b
}

func TestDescribe(t *testing.T) {
	Describe(t, "Parse", func() {
		Context("when SamplingSize is 4", func() {
			winFormat := Parse(ReadBinaryFile("./testdata/4.bin"))

			It("Sequence is 0 ", func() {
				Expect(winFormat.GetSequence()).To(Equal, 0)
			})

			It("SubSequence is  0", func() {
				Expect(winFormat.GetSubSequence()).To(Equal, 0)
			})

			It("SamplingSize is 4", func() {

				Expect(winFormat.GetSamplingSize()).To(Equal, 4)
			})
		})

		Context("when SamplingSize is 3", func() {
			winFormat := Parse(ReadBinaryFile("./testdata/3.bin"))
			It("SamplingSize is 3", func() {
				Expect(winFormat.GetSamplingSize()).To(Equal, 3)
			})
		})

		Context("when SamplingSize is 2", func() {
			winFormat := Parse(ReadBinaryFile("./testdata/2.bin"))
			It("SamplingSize is 2", func() {
				Expect(winFormat.GetSamplingSize()).To(Equal, 2)
			})
		})

		Context("when SamplingSize is 1", func() {
			winFormat := Parse(ReadBinaryFile("./testdata/1.bin"))
			It("SamplingSize is 1", func() {
				Expect(winFormat.GetSamplingSize()).To(Equal, 1)
			})
		})

		Context("when SamplingSize is 0", func() {
			winFormat := Parse(ReadBinaryFile("./testdata/0.bin"))
			It("SamplingSize is 0", func() {
				Expect(winFormat.GetSamplingSize()).To(Equal, 0)
			})
		})

	})

}
