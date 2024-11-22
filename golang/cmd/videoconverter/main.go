package main

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)

func main() {
	println("Hello, Wolrd!")
}

// Ordenar o caminho do arquivo somente com o nome do arquivo
func extractNumber(fileName string) int {
	re := regexp.MustCompile(`\d+`) // Pegar apenas digitos
	numStr := re.FindString(filepath.Base(fileName)) // string

	// Tratamento de erro
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return -1
	}
	return num
}

// Função de merge (criar um novo arquivo a partir dos chunks)
func mergeChunks(inputDir, outPutFile string) error {
	chunks, err := filepath.Glob(filepath.Join(inputDir, "*.chunk")) // Buscar todos os arquivos *.chunk e pegar o caminho completo

	// Tratamento de error
	if err != nil {
		return fmt.Errorf("failed to find chunks: %v", err)
	}

	// Colocar regra de ordenação dos chunks
	sort.Slice(chunks, func(i, j int) bool {
		return extractNumber(chunks[i]) < extractNumber(chunks[j])
	})

}
