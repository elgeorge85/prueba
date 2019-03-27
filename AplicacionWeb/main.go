package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
)

type cliente struct {
	id string
	nomb string
	ape  string
	dni  string
	ema  string
	val  bool
}

//BD
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "74848211"
	dbname   = "BD_apliWeb"
)

//contenido valido
var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}


//Funcion que conecta con la Base de Datos.
func conectaBD() *sql.DB{
	db, err := sql.Open("postgres", "postgres://postgres:74848211@127.0.0.1:5432/BD_apliWeb?sslmode=disable")

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return db
}
//Inserta un nuevo cliente en la BD
func crearCliente(c cliente) error{

	// id - nombre - apellidos - dni - email - validado
	//r.PostFormValue("nombre"), r.PostFormValue("apellidos"), r.PostFormValue("DNI"),r.PostFormValue("Email"), false

	sqlStatement := ` INSERT INTO client (nombre, apellidos, dni, email, validado)
VALUES ($1 , $2 , $3 , $4 , $5 )`

	db := conectaBD()
	defer db.Close()

	stmt , err := db.Prepare(sqlStatement)
	if err != nil {
		return err
	}
	defer stmt.Close()

	r, err := stmt.Exec(c.nomb, c.ape, c.dni, c.ema, false)
	if err != nil {
		return err
	}

	i, _ := r.RowsAffected()
	if i != 1 {
		return errors.New("Se esperaba insertar 1 fila.")
	}
	return nil
}

//Modifica un cliente existente en la BD.
func modificarCliente(c cliente) error{
	sqlStatement := ` update client
set validado = $1
where id = $2`

	db := conectaBD()
	defer db.Close()

	stmt , err := db.Prepare(sqlStatement)
	if err != nil {
		return err
	}
	defer stmt.Close()

	r, err := stmt.Exec(false , c.id)
	if err != nil {
		return err
	}

	i, _ := r.RowsAffected()
	if i != 1 {
		return errors.New("Se esperaba actualizar 1 fila.")
	}
	return nil
}

func borrarCliente(id int)  error {

	sqlStatement := ` delete from client
where id = $1`

	db := conectaBD()
	defer db.Close()

	stmt , err := db.Prepare(sqlStatement)
	if err != nil {
		return err
	}
	defer stmt.Close()

	r, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	i, _ := r.RowsAffected()
	if i != 1 {
		return errors.New("Se esperaba borrar 1 fila.")
	}
	return nil

}

//consulta todos los clientes de la BD.
func consultaCliente()( clientes []cliente, err error){

	sqlStatement := ` select *
from client`

	db := conectaBD()
	defer db.Close()

	rows, err := db.Query(sqlStatement)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next(){
		c := cliente{}
		var id, nombre, apell, dni, email string
		var valid bool
		err = rows.Scan(
			&id,
			&nombre,
			&apell,
			&dni,
			&email,
			&valid,
		)
		if err != nil {
			return
		}

		c.id = id
		c.nomb = nombre
		c.ape = apell
		c.dni = dni
		c.ema = email
		c.val = valid

		clientes = append(clientes, c)
	}

	return clientes, nil
}

//Manenjadores
func indice(w http.ResponseWriter, r *http.Request) {

	//var formular datos

	err := tpl.ExecuteTemplate(w, "index.gohtml", nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	// recibir datos del html
	err = r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	c := cliente{

		nomb: r.PostFormValue("nombre"),
		ape: r.PostFormValue("apellidos"),
		dni: r.PostFormValue("DNI"),
		ema: r.PostFormValue("Email"),
		val: false,
}

	err = crearCliente(c)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Guardado con exito.")
	fmt.Println(c)

}


func validacionC(w http.ResponseWriter, r *http.Request) {

	clientes := []cliente{}

	clientes, err := consultaCliente()
	if err != nil{
		log.Fatal()
	}

	fmt.Println(clientes)
	err = tpl.ExecuteTemplate(w,"validacionAdmin.gohtml", clientes)
	if err != nil{
		fmt.Println("fallo execute validacion", err)
	}


}


func clien(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "cliente.gohtml", nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)

		c := cliente{}
		resultado := []cliente{}
		c.id = r.URL.Query().Get("id")
		db := conectaBD()
		row, err := db.Query("select * from client where id=?", c.id)
		if err != nil{
			log.Fatal()
		}

		for row.Next() {
			var id, nombre, apellidos, dni, email string
			var validado bool
			err = row.Scan(&id,&nombre, &apellidos, &dni, &email, &validado)
			if err != nil {
				panic(err.Error())
			}
			c.nomb = nombre
			c.ape = apellidos
			c.dni = dni
			c.ema = email
			c.val = validado

			resultado = append(resultado, c)
		}

		err = tpl.ExecuteTemplate(w, "cliente.gohtml", resultado)
		if err != nil{
			fmt.Println("fallo execute cliente")
		}
		defer db.Close()



	}


}

func main() {


	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", indice)
	http.HandleFunc("/validacionAdmin", validacionC)
	http.HandleFunc("/cliente", clien)

	http.ListenAndServe(":8080", nil)

}
