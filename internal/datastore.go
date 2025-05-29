package datastore

type Livro struct {
    ID          int
    Título      string
    Disponível  bool
}

type Empréstimo struct {
    ID             int
    Usuário        string
    Livro         string
    DataEmpréstimo string
    DataDevolução  string
}

var livros = map[int]Livro{}
var emprestimos = map[int]Empréstimo{}

func init() {
    // Adicione alguns dados iniciais para teste
    livros[1] = Livro{ID: 1, Título: "Golang Programming", Disponível: true}
}

func obterLivro(livroID int) (Livro, error) {
    livro, exists := livros[livroID]
    if !exists {
        return Livro{}, fmt.Errorf("livro não encontrado")
    }
    return livro, nil
}

func realizarEmpréstimo(usuario, livro, dataEmpréstimo string) error {
    _, err := obterLivro(livro)
    if err != nil {
        return err
    }

    // Verificar se o livro está disponível
    if !livros[livro].Disponível {
        return fmt.Errorf("livro indisponível")
    }

    // Registrar o empréstimo
    id := len(emprestimos) + 1
    emprestimos[id] = Empréstimo{ID: id, Usuário: usuario, Livro: livro, DataEmpréstimo: dataEmpréstimo}
    livros[livro].Disponível = false
    return nil
}

func devolverLivro(usuario, livro string) error {
    _, err := obterLivro(livro)
    if err != nil {
        return err
    }

    // Verificar se o livro está emprestado ao usuário
    for _, e := range emprestimos {
        if e.Usuário == usuario && e.Livro == livro {
            delete(emprestimos, e.ID)
            livros[livro].Disponível = true
            return nil
        }
    }

    return fmt.Errorf("empréstimo não encontrado")
}

func consultarLivro(livro string) (string, error) {
    _, err := obterLivro(livro)
    if err != nil {
        return "", err
    }

    for _, e := range emprestimos {
        if e.Livro == livro {
            return fmt.Sprintf("Emprestado para %s desde %s", e.Usuário, e.DataEmpréstimo), nil
        }
    }

    return "Disponível", nil
}

func consultarEmpréstimosUsuario(usuario string) ([]string, error) {
    var result []string
    for _, e := range emprestimos {
        if e.Usuário == usuario {
            status := "no prazo"
            if time.Now().After(time.Date(e.DataDevolução[:], 0, 0, 0, 0, 0, 0, time.UTC)) {
                status = "atrasado"
            }
            result = append(result, fmt.Sprintf("Título: %s, Data: %s, Status: %s", e.Livro, e.DataEmpréstimo, status))
        }
    }

    if len(result) == 0 {
        return []string{}, fmt.Errorf("nenhum empréstimo encontrado para o usuário")
    }

    return result, nil
}

func consultarEmpréstimosLivro(livro string) ([]string, error) {
    var result []string
    for _, e := range emprestimos {
        if e.Livro == livro {
            result = append(result, fmt.Sprintf("Usuário: %s, Data Empréstimo: %s", e.Usuário, e.DataEmpréstimo))
        }
    }

    if len(result) == 0 {
        return []string{}, fmt.Errorf("nenhum histórico de empréstimos encontrado para o livro")
    }

    return result, nil
}
