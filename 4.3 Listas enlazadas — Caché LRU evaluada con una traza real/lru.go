package main

import "fmt"

// IDpelicula -> clave
// prev y next son punteros al nodo anterior y siguiente
type Nodo struct {
	IDpelicula, valor int
	prev, next        *Nodo
}

// listadoble -> lista completa
// head -> (más reciente) ... (menos reciente) -> tail
type listadoble struct {
	head, tail *Nodo
}

// nuevalistadoble creamos una lista vacía: conectamos los dos nodos ficticios entre sí
func nuevalistadoble() *listadoble {
	head := &Nodo{}                            // creamos el nodo ficticios de inicio
	tail := &Nodo{}                            // creamos el nodo ficticios de fin
	head.next = tail                           // conectamos el next de head con tail
	tail.prev = head                           // conectamos el prev de tail con head
	return &listadoble{head: head, tail: tail} // retornamos el puntero a la lista
}

// agregaralfrente agregamos el nodo n después de head
// marcándolo como el elemento más recientemente usado
func (l *listadoble) agregaralfrente(n *Nodo) {
	n.prev = l.head      // apuntamos el prev de n hacia head
	n.next = l.head.next // apuntamos el next de n hacia el nodo que seguía a head
	l.head.next.prev = n // apuntamos el prev de ese nodo hacia n
	l.head.next = n      // apuntamos el next de head hacia n
}

// borrar desconectamos el nodo n de la lista
func (l *listadoble) borrar(n *Nodo) {
	n.prev.next = n.next // conectamos el nodo anterior con el siguiente
	n.next.prev = n.prev // conectamos el nodo siguiente con el anterior
	n.prev = nil         // limpiamos el puntero prev de n
	n.next = nil         // limpiamos el puntero next de n
}

// llevaralfrente marcamos n como recién usado (solo con nodos que ya estaban en la lista)
// borrar asume que n.prev y n.next no son nil
func (l *listadoble) llevaralfrente(n *Nodo) {
	l.borrar(n)          // lo borramos de su posición actual
	l.agregaralfrente(n) // lo agregamos al frente
}

// nodofinal retornamos el nodo menos recientemente usado
func (l *listadoble) nodofinal() *Nodo {
	if l.tail.prev == l.head { // si el prev de tail apunta a head
		return nil // no hay elementos en la lista
	}
	return l.tail.prev // retornamos el nodo previo a tail
}

// listavacia indicamos si la lista no contiene nodos de datos
func (l *listadoble) listavacia() bool {
	return l.head.next == l.tail // si el next de head apunta a tail
}

// lru desarrollamos una caché de capacidad fija con algoritmo lru
type lru struct {
	capacidadmaxima int           // número máximo de películas guardadas
	mapa            map[int]*Nodo // llave es el ID de película, valor es el puntero al nodo
	lista           *listadoble   // lista doblemente enlazada con el orden de uso
}

// NuevaLRU recibimos la capacidad máxima y retornamos un puntero a la lru
func NuevaLRU(capacidadmaxima int) (*lru, error) {
	if capacidadmaxima <= 0 { // verificamos que la capacidad sea válida
		return nil, fmt.Errorf("la capacidad debe ser mayor que 0, recibido: %d", capacidadmaxima)
	}
	return &lru{
		capacidadmaxima: capacidadmaxima,
		mapa:            make(map[int]*Nodo), // creamos el mapa vacío
		lista:           nuevalistadoble(),   // creamos la lista vacía
	}, nil // sin errores
}

// Get buscamos IDpelicula en la caché
// si existe, la marcamos como recién usada y retornamos su valor con true
// si no existe, retornamos 0 y false.
func (l *lru) Get(IDpelicula int) (int, bool) {
	n, estaenmapa := l.mapa[IDpelicula] // buscamos IDpelicula en el mapa
	if !estaenmapa {
		return 0, false // la película no está en la caché
	}
	l.lista.llevaralfrente(n) // llevamos la película al frente
	return n.valor, true      // retornamos el valor encontrado
}

// Put insertamos o actualizamos IDpelicula en la caché
// si la caché está llena y la clave es nueva, expulsamos al elemento menos recientemente usado
func (l *lru) Put(IDpelicula, valor int) {
	// si la clave ya existía, solo actualizamos su valor
	if n, estaenmapa := l.mapa[IDpelicula]; estaenmapa {
		n.valor = valor           // actualizamos el valor
		l.lista.llevaralfrente(n) // la llevamos al frente
		return
	}

	// si la clave es nueva y la caché está llena, expulsamos al LRU
	if len(l.mapa) >= l.capacidadmaxima {
		nodoantiguo := l.lista.nodofinal() // obtenemos el nodo menos usado
		if nodoantiguo != nil {
			l.lista.borrar(nodoantiguo)            // sacamos el nodo antiguo de la lista
			delete(l.mapa, nodoantiguo.IDpelicula) // lo borramos del mapa
		}
	}

	// creamos el nodo nuevo y lo agregamos al frente (NO usamos
	// llevaralfrente aquí: el nodo recién creado no tiene prev/next)
	n := &Nodo{IDpelicula: IDpelicula, valor: valor}
	l.lista.agregaralfrente(n) // lo colocamos como el más reciente
	l.mapa[IDpelicula] = n     // lo registramos en el mapa
}

// len retornamos cuántas claves hay guardadas en la caché
func (l *lru) Len() int {
	return len(l.mapa) // contamos las claves del mapa
}

// cap retornamos la capacidad máxima configurada al crear la caché
func (l *lru) Cap() int {
	return l.capacidadmaxima
}
