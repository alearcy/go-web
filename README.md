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

### POSTGRES CHEATSHEET

**create a db**

`CREATE DATABASE bookstore;`

**create user**

`CREATE USER bond WITH PASSWORD 'password';`

**grant privileges**

`GRANT ALL PRIVILEGES ON DATABASE bookstore to bond;`

**creare una tabella**

```sql
CREATE TABLE users
(
    id SERIAL NOT NULL,
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
	role INT NOT NULL DEFAULT 1,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);
```

**inserire il primo utente admin**

```sql
insert into users (name, surname, email, password, role) values ("alessandro", "arcidiaco", "arcidiaco.a@gmail.com", "password", 0)
```

### GO SQL COMMANDS

**Create Metodo 1**

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

**Create Metodo 2**

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
