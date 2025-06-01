package main

import (
	"fmt"
	"log"
	"net/rpc"
)

func realizarEmpréstimo(usuario, livro, dataEmpréstimo string) {
	client, err := rpc.Dial("tcp", ":1234")
	if err != nil {
		log.Fatal("Failed to connect to server:", err)
	}

	var reply string
	err = client.Call("Server.RealizarEmpréstimo", struct{ Usuario, Livro, DataEmpréstimo string }{Usuario: usuario, Livro: livro, DataEmpréstimo: dataEmpréstimo}, &reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply)
}

func devolverLivro(usuario, livro string) {
	client, err := rpc.Dial("tcp", ":1234")
	if err != nil {
		log.Fatal("Failed to connect to server:", err)
	}

	var reply string
	err = client.Call("Server.DevolverLivro", struct{ Usuario, Livro string }{Usuario: usuario, Livro: livro}, &reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply)
}

func consultarLivro(livro string) {
	client, err := rpc.Dial("tcp", ":1234")
	if err != nil {
		log.Fatal("Failed to connect to server:", err)
	}

	var reply string
	err = client.Call("Server.ConsultarLivro", livro, &reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply)
}

func consultarEmpréstimosUsuario(usuario string) {
	client, err := rpc.Dial("tcp", ":1234")
	if err != nil {
		log.Fatal("Failed to connect to server:", err)
	}

	var reply string
	err = client.Call("Server.ConsultarEmpréstimosUsuario", usuario, &reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply)
}

func consultarEmpréstimosLivro(livro string) {
	client, err := rpc.Dial("tcp", ":1234")
	if err != nil {
		log.Fatal("Failed to connect to server:", err)
	}

	var reply string
	err = client.Call("Server.ConsultarEmpréstimosLivro", livro, &reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply)
}

func main() {
	// Exemplos de chamada das funções RPC
	realizarEmpréstimo("João", "Golang Programming", "2023-10-01")
	devolverLivro("João", "Golang Programming")
	consultarLivro("Golang Programming")
	consultarEmpréstimosUsuario("João")
	consultarEmpréstimosLivro("Golang Programming")
}
