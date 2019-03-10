package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
)

type datos struct {
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
	password = ""
	dbname   = "ApliClienteAdmin"
)

//contenido valido
var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func indice(w http.ResponseWriter, r *http.Request) {

	var formular datos

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

	//Cambiar struct por un INSERT en la BD
	/*formular.nomb = r.PostFormValue("nombre")
	formular.ape = r.PostFormValue("apellidos")
	formular.dni = r.PostFormValue("DNI")
	formular.ema = r.PostFormValue("Email")*/

	fmt.Println(formular)

	//fmt.Println("method: ", r.Method)

}

func validacionC(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "validacionAdmin.gohtml", nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
		//DEsde la vista del cliente, presionar boton validar y modificar la BD de cliente a validado=true


		//Crear a mano varios clientes y pasarlos a la pagina ValidarCliente para que los cargue

		registro1 := datos{"pedro", "Picapiedra", "11111111A", "PPica@piedra.usa", false}
		//registro2 := datos{"Pablo", "Marmol", "22222222B", "PabloM@marmol.usa", false}
		//registro3 := datos{"tronco", "movil", "99999999E", "troncomovil@troncos.usa", true}

		//tpl.New("rellenrGrid")

		tpl, err = tpl.Parse(`
				<tr>
				<td>{{.nomb}}</td>
                <td>{{.ape}}</td>
                <td>{{.dni}}</td>
                <td>{{.ema}}</td>
				{{if .val == false}}
                <td>No</td>
				{{else}} 
				<td>Si</td>
				{{end}}
                <td><a href="cliente" class="btn btn-info">Ver</a></td>
				</tr>
            `)
	if err != nil {
		log.Fatal("Internal server error", http.StatusInternalServerError)
	}

		err = tpl.Execute(w,  &registro1)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

}

func clien(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "cliente.gohtml", nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)

		//generar template rellenando los datos de un cliente

	}


}

func main() {

	/*Conexión BD

	db, err := sql.Open("postgres", "postgres://postgres:postgres@127.0.0.1:5432/BD_apliWeb?sslmode=disable")

	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	*/ //////////

	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", indice)
	http.HandleFunc("/validacionAdmin", validacionC)
	http.HandleFunc("/cliente", clien)

	http.ListenAndServe(":8080", nil)

}
