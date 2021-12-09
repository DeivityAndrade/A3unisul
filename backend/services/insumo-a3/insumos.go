package insumos

import (
	"backend/banco"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type insumo struct {
	IdProduto      uint32 `json:"id_produto"`
	IdMateriaPrima uint32 `json:"id_materia_prima"`
	Estoque        uint32 `json:"quantidade"`
}

func Create(w http.ResponseWriter, r *http.Request) {
	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		w.Write([]byte("Falha ao ler o corpo da requisição!"))
		return
	}

	var insumo insumo
	if erro = json.Unmarshal(corpoRequisicao, &insumo); erro != nil {
		w.Write([]byte("Erro ao converter o insumo para struct"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao converter conectar no banco de dados!"))
		return
	}
	defer db.Close()

	statement, erro := db.Prepare("insert into insumo (id_produto, id_materia_prima, quantidade) values (?, ?, ?)")
	if erro != nil {
		w.Write([]byte("Erro ao criar o statement!"))
		return
	}
	defer statement.Close()

	insercao, erro := statement.Exec(insumo.IdProduto, insumo.IdMateriaPrima, insumo.Estoque)
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
	w.Write([]byte(fmt.Sprintf("insumo inserido com sucesso! Id: %d", idInserido)))

}

// Retorna todos os insumos salvos no banco de dados
func FindAll(w http.ResponseWriter, r *http.Request) {
	db, erro := banco.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar com o banco de dados!"))
		return
	}
	defer db.Close()

	linhas, erro := db.Query("select * from insumo")
	if erro != nil {
		w.Write([]byte("Erro ao buscar os insumos"))
		return
	}
	defer linhas.Close()

	var insumos []insumo
	for linhas.Next() {
		var insumo insumo

		if erro := linhas.Scan(&insumo.IdProduto, &insumo.IdMateriaPrima, &insumo.Estoque); erro != nil {
			w.Write([]byte("Erro ao escanear o insumo"))
			return
		}

		insumos = append(insumos, insumo)
	}

	w.WriteHeader(http.StatusOK)
	if erro := json.NewEncoder(w).Encode(insumos); erro != nil {
		w.Write([]byte("Erro ao converter os insumos para JSON"))
		return
	}
}

// Retorna um insumo específico salvo no banco de dados
func Find(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	IdProduto, erro := strconv.ParseUint(parametros["id_produto"], 10, 32)
	IdMateriaPrima, erro := strconv.ParseUint(parametros["id_materia_prima"], 10, 32)
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

	linha, erro := db.Query("select * from insumos where id_produto = ? and id_materia_prima = ?", IdProduto, IdMateriaPrima)
	if erro != nil {
		w.Write([]byte("Erro ao buscar o insumo!"))
		return
	}
	defer linha.Close()

	var insumo insumo
	if linha.Next() {
		if erro := linha.Scan(&insumo.IdProduto, &insumo.IdMateriaPrima, &insumo.Estoque); erro != nil {
			w.Write([]byte("Erro ao escanear o insumo!"))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	if erro := json.NewEncoder(w).Encode(insumo); erro != nil {
		w.Write([]byte("Erro ao converter o insumo para JSON!"))
		return
	}
}

// Atualizar os dados de um insumo no banco de dados
func Update(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	IdProduto, erro := strconv.ParseUint(parametros["id_produto"], 10, 32)
	IdMateriaPrima, erro := strconv.ParseUint(parametros["id_materia_prima"], 10, 32)
	if erro != nil {
		w.Write([]byte("Erro ao converter o parâmetro para inteiro"))
		return
	}

	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		w.Write([]byte("Erro ao ler o corpo da requisição!"))
		return
	}

	var insumo insumo
	if erro := json.Unmarshal(corpoRequisicao, &insumo); erro != nil {
		w.Write([]byte("Erro ao converter o insumo para struct"))
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		w.Write([]byte("Erro ao conectar no banco de dados!"))
		return
	}
	defer db.Close()

	statement, erro := db.Prepare("update insumo set quantidade = ? where id_produto = ?, id_materia_prima = ? ")
	if erro != nil {
		w.Write([]byte("Erro ao criar o statement!"))
		return
	}
	defer statement.Close()

	// Anderline abaixo serve para ignorar o retorno (Abaixo só pega se tiver erro)
	if _, erro := statement.Exec(insumo.Estoque, IdProduto, IdMateriaPrima); erro != nil {
		w.Write([]byte("Erro ao atualizar o insumo!"))
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

// Deleta um insumo do banco de dados
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

	statement, erro := db.Prepare("delete from insumo where id = ?")
	if erro != nil {
		w.Write([]byte("Erro ao criar o statement!"))
		return
	}
	defer statement.Close()

	if _, erro := statement.Exec(ID); erro != nil {
		w.Write([]byte("Erro ao deletar o insumo!"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
