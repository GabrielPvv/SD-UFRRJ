package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"

	"sd-ufrrj.local/internal/datastore"
)

type Server struct{}

// Auxiliar para buscar o ID do livro pelo título
func getLivroID(titulo string) (int, error) {
	row := datastore.DB().QueryRow("SELECT id FROM livros WHERE titulo = ?", titulo)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("livro não encontrado")
	}

	return id, nil
}

func (s *Server) RealizarEmpréstimo(args struct {
	Usuario        string
	Livro          string
	DataEmpréstimo string
}, reply *string) error {
	livroID, err := getLivroID(args.Livro)
	if err != nil {
		*reply = err.Error()
		return nil
	}

	err = datastore.RealizarEmprestimo(args.Usuario, livroID, args.DataEmpréstimo)
	if err != nil {
		*reply = err.Error()
		return nil
	}

	*reply = "Sucesso"
	return nil
}

func (s *Server) DevolverLivro(args struct {
	Usuario string
	Livro   string
}, reply *string) error {
	livroID, err := getLivroID(args.Livro)
	if err != nil {
		*reply = err.Error()
		return nil
	}

	err = datastore.DevolverLivro(args.Usuario, livroID)
	if err != nil {
		*reply = err.Error()
		return nil
	}

	*reply = "Sucesso"
	return nil
}

func (s *Server) ConsultarLivro(titulo string, reply *string) error {
	livroID, err := getLivroID(titulo)
	if err != nil {
		*reply = err.Error()
		return nil
	}

	resultado, err := datastore.ConsultarLivro(livroID)
	if err != nil {
		*reply = err.Error()
		return nil
	}

	*reply = resultado
	return nil
}

func (s *Server) ConsultarEmpréstimosUsuario(usuario string, reply *string) error {
	resultados, err := datastore.ConsultarEmprestimosUsuario(usuario)
	if err != nil {
		*reply = err.Error()
		return nil
	}

	*reply = fmt.Sprintf("%v", resultados)
	return nil
}

func (s *Server) ConsultarEmpréstimosLivro(titulo string, reply *string) error {
	livroID, err := getLivroID(titulo)
	if err != nil {
		*reply = err.Error()
		return nil
	}

	resultados, err := datastore.ConsultarEmprestimosLivro(livroID)
	if err != nil {
		*reply = err.Error()
		return nil
	}

	*reply = fmt.Sprintf("%v", resultados)
	return nil
}

func main() {
	// Inicializa o banco de dados SQLite
	if err := datastore.InitDB(); err != nil {
		log.Fatalf("Erro ao inicializar banco: %v", err)
	}

	// Registra o servidor RPC
	rpc.Register(new(Server))
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("Erro ao escutar na porta 1234:", err)
	}

	fmt.Println("Servidor escutando na porta 1234")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print("Erro ao aceitar conexão:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
