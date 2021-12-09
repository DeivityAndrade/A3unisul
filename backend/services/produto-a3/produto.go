package produto

import (
	"backend/banco"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type produto struct {
	ID    uint32 `json:"id"`
	Nome  string `json:"nome"`
	Valor float32 `json:"valor"`
}

func Create(w http.ResponseWriter, r *http.Request) {
	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		w.Write([]byte("Falha ao ler o corpo da requisição!"))
		return
	}

	var produto produto
	if erro = json.Unmarshal(corpoRequisicao, &produto); erro != nil {
		w.Write([]byte("Erro ao converter o produto para struct"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao converter conectar no banco de dados!"))
		return
	}
	defer db.Close()

	statement, erro := db.Prepare("insert into produto (nome, valor) values (?, ?)")
	if erro != nil {
		w.Write([]byte("Erro ao criar o statement!"))
		return
	}
	defer statement.Close()

	insercao, erro := statement.Exec(produto.Nome, produto.Valor)
	if erro != nil {
		w.Write([]byte("Erro ao executar o statement!"))
		return
	}

	idInserido, erro := insercao.LastInsertId()
	if erro != nil {
		w.Write([]byte("Erro ao obter o id inserido!"))
		return
	}

	// STATUS CODES

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Produto inserido com sucesso! Id: %d", idInserido)))

}

// Retorna todos os produtos salvos no banco de dados
func FindAll(w http.ResponseWriter, r *http.Request) {
	db, erro := banco.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar com o banco de dados!"))
		return
	}
	defer db.Close()

	linhas, erro := db.Query("select * from produto")
	if erro != nil {
		w.Write([]byte("Erro ao buscar os produtos"))
		return
	}
	defer linhas.Close()

	var produtos []produto
	for linhas.Next() {
		var produto produto

		if erro := linhas.Scan(&produto.ID, &produto.Nome, &produto.Valor); erro != nil {
			w.Write([]byte("Erro ao escanear o produto"))
			return
		}

		produtos = append(produtos, produto)
	}

	w.WriteHeader(http.StatusOK)
	if erro := json.NewEncoder(w).Encode(produtos); erro != nil {
		w.Write([]byte("Erro ao converter os produtos para JSON"))
		return
	}
}

// Retorna um produto específico salvo no banco de dados
func Find(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	ID, erro := strconv.ParseUint(parametros["id"], 10, 32)
	if erro != nil {
		w.Write([]byte("Erro ao converter o parâmetro para inteiro"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar com o banco de dados!"))
		return
	}
	defer db.Close()

	linha, erro := db.Query("select * from produtos where id = ?", ID)
	if erro != nil {
		w.Write([]byte("Erro ao buscar o produto!"))
		return
	}
	defer linha.Close()

	var produto produto
	if linha.Next() {
		if erro := linha.Scan(&produto.ID, &produto.Nome, &produto.Valor); erro != nil {
			w.Write([]byte("Erro ao escanear o produto!"))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	if erro := json.NewEncoder(w).Encode(produto); erro != nil {
		w.Write([]byte("Erro ao converter o produto para JSON!"))
		return
	}
}

// Atualizar os dados de um produto no banco de dados
func Update(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	ID, erro := strconv.ParseUint(parametros["id"], 10, 32)
	if erro != nil {
		w.Write([]byte("Erro ao converter o parâmetro para inteiro"))
		return
	}

	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		w.Write([]byte("Erro ao ler o corpo da requisição!"))
		return
	}

	var produto produto
	if erro := json.Unmarshal(corpoRequisicao, &produto); erro != nil {
		w.Write([]byte("Erro ao converter o produto para struct"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar no banco de dados!"))
		return
	}
	defer db.Close()

	statement, erro := db.Prepare("update produto set nome = ?, valor = ? where id = ?")
	if erro != nil {
		w.Write([]byte("Erro ao criar o statement!"))
		return
	}
	defer statement.Close()

	// Anderline abaixo serve para ignorar o retorno (Abaixo só pega se tiver erro)
	if _, erro := statement.Exec(produto.Nome, produto.Valor, ID); erro != nil {
		w.Write([]byte("Erro ao atualizar o produto!"))
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

// Deleta um produto do banco de dados
func Delete(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	ID, erro := strconv.ParseUint(parametros["id"], 10, 32)
	if erro != nil {
		w.Write([]byte("Erro ao converter o parâmetro para inteiro"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar no banco de dados!"))
		return
	}
	defer db.Close()

	statement, erro := db.Prepare("delete from produto where id = ?")
	if erro != nil {
		w.Write([]byte("Erro ao criar o statement!"))
		return
	}
	defer statement.Close()

	if _, erro := statement.Exec(ID); erro != nil {
		w.Write([]byte("Erro ao deletar o produto!"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
