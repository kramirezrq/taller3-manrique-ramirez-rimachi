package pilas

import (
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ── Stack tests ──────────────────────────────────────────────

func TestStack_PushPop_Normal(t *testing.T) {
	var s Stack[int]
	s.Push(10)
	s.Push(20)
	s.Push(30)
	for _, want := range []int{30, 20, 10} {
		got, ok := s.Pop()
		if !ok {
			t.Fatalf("Pop() ok=false, quería %d", want)
		}
		if got != want {
			t.Errorf("Pop() = %d, quería %d", got, want)
		}
	}
}

func TestStack_Peek_NoModifica(t *testing.T) {
	var s Stack[string]
	s.Push("a")
	s.Push("b")
	v, ok := s.Peek()
	if !ok || v != "b" {
		t.Fatalf("Peek() = %q, %v; quería \"b\", true", v, ok)
	}
	if s.Len() != 2 {
		t.Errorf("Len() = %d tras Peek; quería 2", s.Len())
	}
}

func TestStack_IsEmpty(t *testing.T) {
	var s Stack[float64]
	if !s.IsEmpty() {
		t.Error("pila nueva debe estar vacía")
	}
	s.Push(3.14)
	if s.IsEmpty() {
		t.Error("con un elemento no debe estar vacía")
	}
	s.Pop()
	if !s.IsEmpty() {
		t.Error("tras Pop del único elemento debe estar vacía")
	}
}

func TestStack_Pop_Vacia(t *testing.T) {
	var s Stack[int]
	v, ok := s.Pop()
	if ok || v != 0 {
		t.Errorf("Pop() en vacía = %d, %v; quería 0, false", v, ok)
	}
}

func TestStack_Peek_Vacia(t *testing.T) {
	var s Stack[int]
	_, ok := s.Peek()
	if ok {
		t.Error("Peek() en vacía debe devolver false")
	}
}

func TestStack_UnSoloElemento(t *testing.T) {
	var s Stack[int]
	s.Push(42)
	top, _ := s.Peek()
	if top != 42 {
		t.Errorf("Peek() = %d; quería 42", top)
	}
	v, ok := s.Pop()
	if !ok || v != 42 {
		t.Errorf("Pop() = %d, %v; quería 42, true", v, ok)
	}
	if !s.IsEmpty() {
		t.Error("debe quedar vacía")
	}
}

// ── CalcularStockSpan tests ──────────────────────────────────

func TestCalcularStockSpan_Normal(t *testing.T) {
	precios := []float64{100, 80, 60, 70, 60, 75, 85}
	want := []int{1, 1, 1, 2, 1, 4, 6}
	got := CalcularStockSpan(precios)
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("span[%d] = %d; quería %d", i, got[i], want[i])
		}
	}
}

func TestCalcularStockSpan_PreciosCrecientes(t *testing.T) {
	precios := []float64{1, 2, 3, 4, 5}
	got := CalcularStockSpan(precios)
	for i, v := range got {
		if v != i+1 {
			t.Errorf("span[%d] = %d; quería %d", i, v, i+1)
		}
	}
}

func TestCalcularStockSpan_PreciosDecrecientes(t *testing.T) {
	precios := []float64{5, 4, 3, 2, 1}
	got := CalcularStockSpan(precios)
	for i, v := range got {
		if v != 1 {
			t.Errorf("span[%d] = %d; quería 1", i, v)
		}
	}
}

func TestCalcularStockSpan_PreciosIguales(t *testing.T) {
	precios := []float64{5, 5, 5, 5}
	want := []int{1, 2, 3, 4}
	got := CalcularStockSpan(precios)
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("span[%d] = %d; quería %d", i, got[i], want[i])
		}
	}
}

func TestCalcularStockSpan_UnElemento(t *testing.T) {
	got := CalcularStockSpan([]float64{42.0})
	if len(got) != 1 || got[0] != 1 {
		t.Errorf("resultado = %v; quería [1]", got)
	}
}

func TestCalcularStockSpan_Vacio(t *testing.T) {
	got := CalcularStockSpan([]float64{})
	if len(got) != 0 {
		t.Errorf("resultado = %v; quería []", got)
	}
}

// ── LeerPrecios tests ────────────────────────────────────────

func csvTemporal(t *testing.T, contenido string) string {
	t.Helper()
	dir := t.TempDir()
	ruta := filepath.Join(dir, "datos.txt")
	os.WriteFile(ruta, []byte(contenido), 0o644)
	return ruta
}

func TestLeerPrecios_Normal(t *testing.T) {
	csv := "Date,Open,High,Low,Close,Volume,OpenInt\n" +
		"2000-01-03,4.50,4.75,4.40,4.68,123456,0\n" +
		"2000-01-04,4.68,4.90,4.60,4.85,234567,0\n"
	precios, err := LeerPrecios(csvTemporal(t, csv))
	if err != nil || len(precios) != 2 || precios[0] != 4.68 {
		t.Errorf("error inesperado: %v, precios=%v", err, precios)
	}
}

func TestLeerPrecios_ArchivoNoExiste(t *testing.T) {
	_, err := LeerPrecios("/no/existe.txt")
	if err == nil {
		t.Error("se esperaba error")
	}
}

func TestLeerPrecios_SinColumnaClose(t *testing.T) {
	csv := "Date,Open,High,Low,Volume\n2000-01-03,4.50,4.75,4.40,123456\n"
	_, err := LeerPrecios(csvTemporal(t, csv))
	if err == nil {
		t.Error("se esperaba error por columna Close ausente")
	}
}

func TestLeerPrecios_SinFilasValidas(t *testing.T) {
	_, err := LeerPrecios(csvTemporal(t, "Date,Open,High,Low,Close,Volume,OpenInt\n"))
	if err == nil {
		t.Error("se esperaba error")
	}
}

func TestLeerPrecios_EncabezadoVacio(t *testing.T) {
	_, err := LeerPrecios(csvTemporal(t, ""))
	if err == nil {
		t.Error("se esperaba error")
	}
}

func TestLeerPrecios_FilaConValorNoNumerico(t *testing.T) {
	csv := "Date,Open,High,Low,Close,Volume,OpenInt\n" +
		"2000-01-03,4.50,4.75,4.40,N/A,123456,0\n" +
		"2000-01-04,4.68,4.90,4.60,4.85,234567,0\n"
	precios, err := LeerPrecios(csvTemporal(t, csv))
	if err != nil || len(precios) != 1 || precios[0] != 4.85 {
		t.Errorf("error=%v precios=%v; quería [4.85]", err, precios)
	}
}

// ── LeerRegistros tests ──────────────────────────────────────

func TestLeerRegistros_Normal(t *testing.T) {
	data := "Date,Open,High,Low,Close,Volume,OpenInt\n" +
		"2000-01-03,4.50,4.75,4.40,4.68,123456,0\n"
	regs, err := LeerRegistros(csvTemporal(t, data))
	if err != nil || len(regs) != 1 || regs[0].Fecha != "2000-01-03" {
		t.Errorf("error=%v regs=%v", err, regs)
	}
}

func TestLeerRegistros_SinClose(t *testing.T) {
	_, err := LeerRegistros(csvTemporal(t, "Date,Open\n2000-01-03,4.50\n"))
	if err == nil {
		t.Error("se esperaba error")
	}
}

func TestLeerRegistros_ArchivoNoExiste(t *testing.T) {
	_, err := LeerRegistros("/no/existe.txt")
	if err == nil {
		t.Error("se esperaba error")
	}
}

func TestLeerRegistros_SinFilas(t *testing.T) {
	_, err := LeerRegistros(csvTemporal(t, "Date,Open,High,Low,Close,Volume,OpenInt\n"))
	if err == nil {
		t.Error("se esperaba error")
	}
}

// ── DiaMayorSpan tests ───────────────────────────────────────

func TestDiaMayorSpan_Normal(t *testing.T) {
	if got := DiaMayorSpan([]int{1, 1, 1, 2, 1, 4, 6}); got != 6 {
		t.Errorf("DiaMayorSpan = %d; quería 6", got)
	}
}

func TestDiaMayorSpan_Vacio(t *testing.T) {
	if got := DiaMayorSpan(nil); got != -1 {
		t.Errorf("DiaMayorSpan(nil) = %d; quería -1", got)
	}
}

func TestDiaMayorSpan_Empate(t *testing.T) {
	if got := DiaMayorSpan([]int{3, 3, 3}); got != 0 {
		t.Errorf("empate = %d; quería 0", got)
	}
}

// ── Benchmarks ───────────────────────────────────────────────

func generarPrecios(n int) []float64 {
	r := rand.New(rand.NewSource(42))
	p := make([]float64, n)
	for i := range p {
		p[i] = float64(r.Intn(1000) + 1)
	}
	return p
}

func BenchmarkCalcularStockSpan_1K(b *testing.B) {
	precios := generarPrecios(1_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CalcularStockSpan(precios)
	}
}

func BenchmarkCalcularStockSpan_10K(b *testing.B) {
	precios := generarPrecios(10_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CalcularStockSpan(precios)
	}
}

func BenchmarkCalcularStockSpan_100K(b *testing.B) {
	precios := generarPrecios(100_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CalcularStockSpan(precios)
	}
}

func BenchmarkCalcularStockSpan_PeorCaso_100K(b *testing.B) {
	precios := make([]float64, 100_000)
	for i := range precios {
		precios[i] = float64(i + 1)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CalcularStockSpan(precios)
	}
}

func BenchmarkLeerPrecios_10K(b *testing.B) {
	var sb strings.Builder
	sb.WriteString("Date,Open,High,Low,Close,Volume,OpenInt\n")
	for i := 0; i < 10_000; i++ {
		sb.WriteString("2000-01-01,100.0,105.0,99.0,102.5,1000000,0\n")
	}
	ruta := filepath.Join(b.TempDir(), "bench.txt")
	os.WriteFile(ruta, []byte(sb.String()), 0o644)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LeerPrecios(ruta)
	}
}
