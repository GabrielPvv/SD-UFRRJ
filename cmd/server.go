package main

import (
    "fmt"
    "log"
    "net/rpc"
    "path/filepath"

    "myproject/internal/datastore"
)

type Server struct{}

func (s *Server) RealizarEmpréstimo(usuario, livro, dataEmpréstimo string, reply *string) error {
    err := datastore.RealizarEmpréstimo(usuario, livro, dataEmpréstimo)
    if err != nil {
        *reply = err.Error()
        return nil
    }
    *reply = "Sucesso"
    return nil
}

func (s *Server) DevolverLivro(usuario, livro string, reply *string) error {
    err := datastore.DevolverLivro(usuario, livro)
    if err != nil {
        *reply = err.Error()
        return nil
    }
    *reply = "Sucesso"
    return nil
}

func (s *Server) ConsultarLivro(livro string, reply *string) error {
    resultado, err := datastore.ConsultarLivro(livro)
    if err != nil {
        *reply = err.Error()
        return nil
    }
    *reply = resultado
    return nil
}

func (s *Server) ConsultarEmpréstimosUsuario(usuario string, reply *string) error {
    resultados, err := datastore.ConsultarEmpréstimosUsuario(usuario)
    if err != nil {
        *reply = err.Error()
        return nil
    }
    *reply = fmt.Sprintf("%v", resultados)
    return nil
}

func (s *Server) ConsultarEmpréstimosLivro(livro string, reply *string) error {
    resultados, err := datastore.ConsultarEmpréstimosLivro(livro)
    if err != nil {
        *reply = err.Error()
        return nil
    }
    *reply = fmt.Sprintf("%v", resultados)
    return nil
}

func main() {
    rpc.Register(new(Server))
    listener, err := net.Listen("tcp", ":1234")
    if err != nil {
        log.Fatal("Failed to listen on port 1234:", err)
    }

    fmt.Println("Server listening on :1234")
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Print("Failed to accept connection:", err)
            continue
        }
        go rpc.ServeConn(conn)
    }
}
