package materiaprima

import (
	"backend/banco"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type materiaprima struct {
	ID      uint32 `json:"id"`
	Nome    string `json:"nome"`
	Estoque uint32 `json:"estoque"`
}

func Create(w http.ResponseWriter, r *http.Request) {
	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		w.Write([]byte("Falha ao ler o corpo da requisição!"))
		return
	}

	var materiaprima materiaprima
	if erro = json.Unmarshal(corpoRequisicao, &materiaprima); erro != nil {
		w.Write([]byte("Erro ao converter a Matéria-prima para struct"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao converter conectar no banco de dados!"))
		return
	}
	defer db.Close()

	statement, erro := db.Prepare("insert into `metaria-prima` (nome, estoque) values (?, ?)")
	if erro != nil {
		w.Write([]byte("1- Erro ao criar o statement!"))
		return
	}
	defer statement.Close()

	insercao, erro := statement.Exec(materiaprima.Nome, materiaprima.Estoque)
	if erro != nil {
		w.Write([]byte("2- Erro ao executar o statement!"))
		return
	}

	idInserido, erro := insercao.LastInsertId()
	if erro != nil {
		w.Write([]byte("Erro ao obter o id inserido!"))
		return
	}

	// STATUS CODES

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Matéria-prima inserida com sucesso! Id: %d", idInserido)))

}

// Retorna todas as materias-primas salvos no banco de dados
func FindAll(w http.ResponseWriter, r *http.Request) {
	db, erro := banco.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar com o banco de dados!"))
		return
	}
	defer db.Close()

	linhas, erro := db.Query("select * from `metaria-prima`")
	if erro != nil {
		w.Write([]byte("Erro ao buscar as Matérias-primas"))
		return
	}
	defer linhas.Close()

	var materiasprimas []materiaprima
	for linhas.Next() {
		var materiaprima materiaprima

		if erro := linhas.Scan(&materiaprima.ID, &materiaprima.Nome, &materiaprima.Estoque); erro != nil {
			w.Write([]byte("Erro ao escanear a matéria-prima"))
			return
		}

		materiasprimas = append(materiasprimas, materiaprima)
	}

	w.WriteHeader(http.StatusOK)
	if erro := json.NewEncoder(w).Encode(materiasprimas); erro != nil {
		w.Write([]byte("Erro ao converter as matérias-primas para JSON"))
		return
	}
}

// Retorna uma matéria-prima específico salvo no banco de dados
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

	linha, erro := db.Query("select * from `metaria-prima` where id = ?", ID)
	if erro != nil {
		w.Write([]byte("Erro ao buscar o matéria-prima!"))
		return
	}
	defer linha.Close()

	var materiaprima materiaprima
	if linha.Next() {
		if erro := linha.Scan(&materiaprima.ID, &materiaprima.Nome, &materiaprima.Estoque); erro != nil {
			w.Write([]byte("Erro ao escanear a matéria-prima!"))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	if erro := json.NewEncoder(w).Encode(materiaprima); erro != nil {
		w.Write([]byte("Erro ao converter a matéria-prima para JSON!"))
		return
	}
}

// Atualizar os dados de um matéria-prima no banco de dados
func Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Erro ao converter o parâmetro para inteiro"))
	return
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

	var materiaprima materiaprima
	if erro := json.Unmarshal(corpoRequisicao, &materiaprima); erro != nil {
		w.Write([]byte("Erro ao converter a matéria-prima para struct"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar no banco de dados!"))
		return
	}
	defer db.Close()

	statement, erro := db.Prepare("update `metaria-prima` set nome = ?, estoque = ? where id = ?")
	if erro != nil {
		w.Write([]byte("Erro ao criar o statement!"))
		return
	}
	defer statement.Close()

	// Anderline abaixo serve para ignorar o retorno (Abaixo só pega se tiver erro)
	if _, erro := statement.Exec(materiaprima.Nome, materiaprima.Estoque, ID); erro != nil {
		w.Write([]byte("Erro ao atualizar o Matéria-prima!"))
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

// Deleta uma matéria-prima do banco de dados
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

	statement, erro := db.Prepare("delete from `metaria-prima` where id = ?")
	if erro != nil {
		w.Write([]byte("Erro ao criar o statement!"))
		return
	}
	defer statement.Close()

	if _, erro := statement.Exec(ID); erro != nil {
		w.Write([]byte("Erro ao deletar a Matéria-prima!"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
