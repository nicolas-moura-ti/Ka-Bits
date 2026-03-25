package ui

import (
	"strings"
	"testing"
)

// old implementation
func renderDataRainOld(dataRain []string) string {
	rain := ""
	for _, char := range dataRain {
		rain += char + "\n"
	}
	return StyleDataRain.Render(rain)
}

// new implementation
func renderDataRainNew(dataRain []string) string {
	var builder strings.Builder
	// We know the length roughly
	builder.Grow(len(dataRain) * 2)
	for _, char := range dataRain {
		builder.WriteString(char)
		builder.WriteByte('\n')
	}
	return StyleDataRain.Render(builder.String())
}

func BenchmarkRenderDataRainOld(b *testing.B) {
	dataRain := make([]string, 15)
	for i := range dataRain {
		dataRain[i] = " "
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		renderDataRainOld(dataRain)
	}
}

func BenchmarkRenderDataRainNew(b *testing.B) {
	dataRain := make([]string, 15)
	for i := range dataRain {
		dataRain[i] = " "
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		renderDataRainNew(dataRain)
	}
}
