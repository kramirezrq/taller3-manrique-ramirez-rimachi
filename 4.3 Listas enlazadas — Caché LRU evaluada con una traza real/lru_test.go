package main

import "testing"

// caso de error

// verificamos que crear una caché con capacidad cero o negativa retorne un error
func TestNuevaLRU_CapacidadInvalida(t *testing.T) {
	if _, err := NuevaLRU(0); err == nil {
		t.Error("se esperaba un error al crear una caché con capacidad 0")
	}

	if _, err := NuevaLRU(-5); err == nil {
		t.Error("se esperaba un error al crear una caché con capacidad negativa")
	}
}

// aso límite: caché vacía

// verificamos que get sobre una caché recién creada retorne (0, false)
func TestGet_CacheVacia(t *testing.T) {
	cache, err := NuevaLRU(2)
	if err != nil {
		t.Fatalf("no se esperaba error al crear la caché: %v", err)
	}

	if _, encontrado := cache.Get(1); encontrado {
		t.Error("no se esperaba encontrar nada en una caché vacía")
	}
}

// caso normal: agregar y leer

// verificamos que después de muchos put, los get correspondientes retornen los valores correctos
func TestPutGet_CasoNormal(t *testing.T) {
	cache, _ := NuevaLRU(2)

	cache.Put(1, 100) // IDpelicula y valor
	cache.Put(2, 200) // IDpelicula y valor

	if valor, encontrado := cache.Get(1); !encontrado || valor != 100 {
		t.Errorf("se esperaba (100, true), se obtuvo (%d, %v)", valor, encontrado)
	}

	if valor, encontrado := cache.Get(2); !encontrado || valor != 200 {
		t.Errorf("se esperaba (200, true), se obtuvo (%d, %v)", valor, encontrado)
	}
}

// caso normal: actualizar una clave existente

// verificamos que put sobre una clave que ya existe actualice su valor
func TestPut_ActualizaValorExistente(t *testing.T) {
	cache, _ := NuevaLRU(2)

	cache.Put(1, 100)
	cache.Put(1, 999) // actualizamos la misma clave

	if cache.Len() != 1 {
		t.Errorf("se esperaba 1 elemento, hay %d", cache.Len())
	}

	if valor, _ := cache.Get(1); valor != 999 {
		t.Errorf("se esperaba 999, se obtuvo %d", valor)
	}
}

// caso límite: borrar el menos recientemente usado

// verificamos que, al llenar la caché y agregar una clave nueva, se borre el lru
func TestPut_ExpulsaMenosRecienteUsado(t *testing.T) {
	cache, _ := NuevaLRU(2)

	cache.Put(1, 10) // estado: [1]
	cache.Put(2, 20) // estado: [2, 1]
	cache.Put(3, 30) // caché llena -> expulsa 1 -> estado: [3, 2]

	if _, encontrado := cache.Get(1); encontrado {
		t.Error("se esperaba que la clave 1 hubiera sido expulsada")
	}

	if _, encontrado := cache.Get(2); !encontrado {
		t.Error("se esperaba que la clave 2 siguiera en la caché")
	}

	if _, encontrado := cache.Get(3); !encontrado {
		t.Error("se esperaba que la clave 3 estuviera en la caché")
	}

	if cache.Len() != 2 {
		t.Errorf("se esperaban 2 elementos, hay %d", cache.Len())
	}
}

// caso normal: un acceso reciente evita la expulsión

// verificamos que un get sobre una clave la marque como lru
func TestPut_AccesoEvitaExpulsion(t *testing.T) {
	cache, _ := NuevaLRU(2)

	cache.Put(1, 10) // estado: [1]
	cache.Put(2, 20) // estado: [2, 1]
	cache.Get(1)     // accedemos a 1 -> ahora es el más reciente -> [1, 2]
	cache.Put(3, 30) // caché llena -> expulsa 2 (el menos reciente ahora) -> [3, 1]

	if _, encontrado := cache.Get(2); encontrado {
		t.Error("se esperaba que la clave 2 hubiera sido expulsada")
	}

	if _, encontrado := cache.Get(1); !encontrado {
		t.Error("se esperaba que la clave 1 siguiera en la caché (se accedió recientemente)")
	}
}

// caso límite: len y cap

// verificamos que cap retorne siempre la capacidad configurada, y len la cantidad de elementos
func TestLen_y_Cap(t *testing.T) {
	cache, _ := NuevaLRU(3)

	if cache.Cap() != 3 {
		t.Errorf("se esperaba capacidad 3, se obtuvo %d", cache.Cap())
	}

	if cache.Len() != 0 {
		t.Errorf("se esperaban 0 elementos, hay %d", cache.Len())
	}

	cache.Put(1, 10)
	cache.Put(2, 20)

	if cache.Len() != 2 {
		t.Errorf("se esperaban 2 elementos, hay %d", cache.Len())
	}

	// el tamaño nunca debe superar la capacidad
	cache.Put(3, 30)
	cache.Put(4, 40)

	if cache.Len() > cache.Cap() {
		t.Errorf("la caché excedió su capacidad: Len=%d, Cap=%d", cache.Len(), cache.Cap())
	}
}

// benchmarks

// medimos el costo de put de forma repetida
// como put es O(1), esperamos que el tiempo por operación se mantenga constante
func BenchmarkPut(b *testing.B) {
	cache, _ := NuevaLRU(1000)
	for i := 0; i < b.N; i++ {
		cache.Put(i, i)
	}
}

// medimos el costo de get de forma repetida sobre una caché ya llena
// como get es O(1), esperamos un tiempo por operación constante
func BenchmarkGet(b *testing.B) {
	cache, _ := NuevaLRU(1000)
	for i := 0; i < 1000; i++ {
		cache.Put(i, i) // llenamos la caché antes de medir
	}

	b.ResetTimer() // no contamos el tiempo de llenado, solo el de Get
	for i := 0; i < b.N; i++ {
		cache.Get(i % 1000) // recorremos las claves existentes en ciclo
	}
}
