# Appunti GO

Con `go run main.go > index.html` genero un file html eseguibile.

## SQL

### Leggere dal DB

```go
func amigos(w http.ResponseWriter, req *http.Request) {
    // per fare select uso query che mi ritorna le ROWS
	rows, err := db.Query(`SELECT aName FROM amigos;`)
	check(err)
	defer rows.Close()

	// preparo le variabili su cui assegnare i valori delle rows
	var s, name string
	s = "RETRIEVED RECORDS:\n"

	// ciclo fintanto che next() ritorna true per ogni row che trova e assegno con Scan il valore alla variabile
	for rows.Next() {
		err = rows.Scan(&name)
		check(err)
		s += name + "\n"
	}
	fmt.Fprintln(w, s)
}
```

### Scrivere sul DB

**Metodo 1**

```go
func create(w http.ResponseWriter, req *http.Request) {
    // uso PREPARE per preparare lo statment che userò poi successivamente
	stmt, err := db.Prepare(`CREATE TABLE customer (name VARCHAR(20));`)
	check(err)
	defer stmt.Close()

    // eseguo la query
	r, err := stmt.Exec()
	check(err)

    // opzionale se voglio leggere quante rows sono state toccate
	n, err := r.RowsAffected()
	check(err)

	fmt.Fprintln(w, "CREATED TABLE customer", n)
}
```

**Metodo 2**

```go
func create(w http.ResponseWriter, req *http.Request) {
    // uso direttamente EXEC
	res, err := db.Exec(`CREATE TABLE customer (name VARCHAR(20));`)
	    panic(err)
    }

    // opzionale se voglio leggere quante rows sono state toccate
    numDeleted, err := res.RowsAffected()
    if err != nil {
        panic(err)
    }
    print(numDeleted)
}
```
