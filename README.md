# Cachefy

Cachefy es una biblioteca de caché flexible para Go, compatible con cachés en memoria y persistentes. Ofrece varios backends, incluyendo cachés en memoria con sharding y capas de persistencia que utilizan SQLite o Postgres.

## Características

- Caché en memoria utilizando `sync.Map`, `RWMutex` y backends con sharding.
- Persistencia opcional con SQLite (almacenamiento de blobs binarios) o Postgres (almacenamiento en formato JSON).
- Configuración de tiempo de vida (TTL) para las entradas de la caché.
- Soporte para paginación para una recuperación eficiente desde almacenes persistentes.

## Instalación

```bash
go get github.com/ha1tch/cachefy
```

## Uso

### Caché Básico en Memoria

```go
config := CacheConfig{
    DefaultTTL: 5 * time.Minute,
    Backend:    "rwmutex",
}

cache, err := NewCache(config)
if err != nil {
    log.Fatalf("No se pudo inicializar la caché: %v", err)
}

cache.Set("key1", "value1")
value, _ := cache.Get("key1")
fmt.Println(value) // Salida: value1
```

### Caché en Memoria con Sharding

```go
config := CacheConfig{
    DefaultTTL:    5 * time.Minute,
    Backend:       "sharded",
    Shards:        4,
    ShardCapacity: 100,
}

cache, err := NewCache(config)
if err != nil {
    log.Fatalf("No se pudo inicializar la caché con sharding: %v", err)
}
```

### Caché Persistente (SQLite)

```go
config := CacheConfig{
    DefaultTTL:        5 * time.Minute,
    Backend:           "rwmutex",
    EnablePersistence: true,
    DatabaseType:      "sqlite",
    DatabaseDSN:       "cache.db",
}

cache, err := NewCache(config)
if err != nil {
    log.Fatalf("No se pudo inicializar la caché persistente: %v", err)
}
```

### Caché Persistente (Postgres)

```go
config := CacheConfig{
    DefaultTTL:        5 * time.Minute,
    Backend:           "rwmutex",
    EnablePersistence: true,
    DatabaseType:      "postgres",
    DatabaseDSN:       "user=your_user password=your_password dbname=your_db sslmode=disable",
}

cache, err := NewCache(config)
if err != nil {
    log.Fatalf("No se pudo inicializar la caché persistente: %v", err)
}
```

## Tests

Ejecutar los tests con:
```bash
go test ./... -v
```

