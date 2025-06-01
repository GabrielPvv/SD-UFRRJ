package datastore

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

var db *sql.DB

type Livro struct {
	ID         int
	Titulo     string
	Disponivel bool
}

type Emprestimo struct {
	ID             int
	Usuario        string
	LivroID        int
	DataEmprestimo string
	DataDevolucao  string
}

func InitDB() error {
	var err error
	db, err = sql.Open("sqlite", "./biblioteca.db")
	if err != nil {
		return err
	}

	createTables := `
    CREATE TABLE IF NOT EXISTS livros (
        id INTEGER PRIMARY KEY,
        titulo TEXT NOT NULL,
        disponivel BOOLEAN NOT NULL
    );

    CREATE TABLE IF NOT EXISTS emprestimos (
        id INTEGER PRIMARY KEY,
        usuario TEXT NOT NULL,
        livro_id INTEGER NOT NULL,
        data_emprestimo TEXT,
        data_devolucao TEXT,
        FOREIGN KEY(livro_id) REFERENCES livros(id)
    );`
	_, err = db.Exec(createTables)
	if err != nil {
		return err
	}

	_, err = db.Exec(`INSERT OR IGNORE INTO livros (id, titulo, disponivel) VALUES (1, 'Golang Programming', 1)`)
	return err
}

func DB() *sql.DB {
	return db
}

func ObterLivro(id int) (Livro, error) {
	var l Livro
	err := db.QueryRow("SELECT id, titulo, disponivel FROM livros WHERE id = ?", id).
		Scan(&l.ID, &l.Titulo, &l.Disponivel)
	if err != nil {
		return Livro{}, fmt.Errorf("livro não encontrado")
	}
	return l, nil
}

func RealizarEmprestimo(usuario string, livroID int, dataEmprestimo string) error {
	livro, err := ObterLivro(livroID)
	if err != nil {
		return err
	}
	if !livro.Disponivel {
		return fmt.Errorf("livro indisponível")
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO emprestimos (usuario, livro_id, data_emprestimo) VALUES (?, ?, ?)",
		usuario, livroID, dataEmprestimo)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("UPDATE livros SET disponivel = 0 WHERE id = ?", livroID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func DevolverLivro(usuario string, livroID int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	var emprestimoID int
	row := tx.QueryRow("SELECT id FROM emprestimos WHERE usuario = ? AND livro_id = ?", usuario, livroID)
	err = row.Scan(&emprestimoID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("empréstimo não encontrado")
	}

	_, err = tx.Exec("DELETE FROM emprestimos WHERE id = ?", emprestimoID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("UPDATE livros SET disponivel = 1 WHERE id = ?", livroID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func ConsultarLivro(livroID int) (string, error) {
	livro, err := ObterLivro(livroID)
	if err != nil {
		return "", err
	}
	if livro.Disponivel {
		return "Disponível", nil
	}

	var usuario, data string
	err = db.QueryRow("SELECT usuario, data_emprestimo FROM emprestimos WHERE livro_id = ?", livroID).
		Scan(&usuario, &data)
	if err != nil {
		return "", fmt.Errorf("erro ao consultar empréstimo")
	}

	return fmt.Sprintf("Emprestado para %s desde %s", usuario, data), nil
}

func ConsultarEmprestimosUsuario(usuario string) ([]string, error) {
	rows, err := db.Query("SELECT livro_id, data_emprestimo, data_devolucao FROM emprestimos WHERE usuario = ?", usuario)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var livroID int
		var dataEmprestimo, dataDevolucao string
		rows.Scan(&livroID, &dataEmprestimo, &dataDevolucao)

		status := "no prazo"
		if dataDevolucao != "" {
			dt, _ := time.Parse("2006-01-02", dataDevolucao)
			if time.Now().After(dt) {
				status = "atrasado"
			}
		}

		livro, _ := ObterLivro(livroID)
		result = append(result, fmt.Sprintf("Título: %s, Data: %s, Status: %s", livro.Titulo, dataEmprestimo, status))
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("nenhum empréstimo encontrado para o usuário")
	}

	return result, nil
}

func ConsultarEmprestimosLivro(livroID int) ([]string, error) {
	rows, err := db.Query("SELECT usuario, data_emprestimo FROM emprestimos WHERE livro_id = ?", livroID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var usuario, data string
		rows.Scan(&usuario, &data)
		result = append(result, fmt.Sprintf("Usuário: %s, Data Empréstimo: %s", usuario, data))
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("nenhum histórico de empréstimos encontrado para o livro")
	}

	return result, nil
}
