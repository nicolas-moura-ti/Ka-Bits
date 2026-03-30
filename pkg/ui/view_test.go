package ui

import (
	"testing"
)

func BenchmarkRenderDataRain(b *testing.B) {
	dataRain := make([]string, 15)
	for i := range dataRain {
		dataRain[i] = " "
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		renderDataRain(dataRain)
	}
}
