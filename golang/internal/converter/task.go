package converter

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"imersaofc/internal/rabbitmq"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

type VideoConverter struct {
	db *sql.DB
	rabbitClient *rabbitmq.RabbitClient 
}

func NewVideoConverter(rabbitClient *rabbitmq.RabbitClient , db *sql.DB) *VideoConverter {
	return &VideoConverter{
		db: db,
	}
}

// {"video_id": 1, "path": "media/uploads/1"}
type VideoTask struct{
	VideoID int `json:"video_id"` // Toda vez que tiver um json e quiser converter ele para o formato VideoTask, toda vez que ele encontrar video_id ele está se referindo ao video_id do json
	Path string `json:"path` // Aqui referindo ao path do json
}

// Converter VideoTask para json
func (vc *VideoConverter) Handle(d amqp.Delivery, conversionExch, confirmationKey, confirmationQueue string) { // Esse * é um ponteiro, qualquer valor que eu auterar aqui utilizando o vc. vai ser refletido em qualquer lugar do meu código)
	var task VideoTask  
	err := json.Unmarshal(d.Body, &task) // Pegar a mensagem que vai ser o json de input e vai converter ele no formado de VideoTask
	if err != nil {
		vc.logError(task, "failed to unmarshal task", err)
		return
	}

	if IsProcessed(vc.db, task.VideoID) {
		slog.Warn("Video already processed", slog.Int("video_id", task.VideoID))
		d.Ack(false) // Se a mensagem já foi processada eu removo ela da fila
		return
	}

	err = MarkProcessed(vc.db, task.VideoID)
	if err != nil {
		vc.logError(task, "failed to mark video as processed", err)
		return
	}
	d.Ack(false) // Avisar que já foi processado e pode jogar a mensagem para fora da fila 
	slog.Info("Video marked as processed", slog.Int("video_id", task.VideoID))

	confirmationMessage := []byte(fmt.Sprintf(`{"video_id": %d, "path":"%s"}`, task.VideoID, task.Path))
	err = vc.rabbitClient.PublishMessage(conversionExch, confirmationKey, confirmationQueue, confirmationMessage)

	// Pegar essa task que a gente tem e mandar processar
	err = vc.processVideo(&task)
	if err != nil {
		vc.logError(task, "failed to process video", err)
		return
	}
}

// Processar o vídeo
func (vc *VideoConverter) processVideo(task *VideoTask) error { // O * siginifica que qualquer lugar que eu mudar o VideoTask ele muda no programa inteiro
	mergedFile := filepath.Join(task.Path, "merged.mp4") // Caminho aonde eu vou gerar o meu arquivo mergeado
	mpegDashPath := filepath.Join(task.Path, "mpeg-dash") // Qual é a pasta que vou converter esses arquivos

	slog.Info("Merging chunks", slog.String("path", task.Path))
	err := vc.mergeChunks(task.Path, mergedFile)
	if err != nil {
		vc.logError(*task, "failed to merge chunks", err)
		return fmt.Errorf("failed to merge chunks")
	}

	// Criar a pasta para converter os arquivos
	slog.Info("Creating mpeg-dash dir", slog.String("path", task.Path))
	err = os.MkdirAll(mpegDashPath, os.ModePerm) // Esse ModePerm serve para eu ter permissão dentro do diretório
	 if err != nil {
		vc.logError(*task, "failed to create mpeg-dash directory", err)
		return err
	 }

	 slog.Info("Converting video to mpeg-dash", slog.String("path", task.Path))
	 ffmpegCmd := exec.Command(
		"ffmpeg", "-i", mergedFile,
		"-f", "dash",
		filepath.Join(mpegDashPath, "output.mpd"),
	)

	output, err :=  ffmpegCmd.CombinedOutput()
	if err != nil {
		vc.logError(*task, "failed to convert video to mpeg-dash, output: " +string(output), err)
		return err
	}
	slog.Info("Video converted to mpeg-dash", slog.String("path", mpegDashPath))

	// Remover o arquivo merged.mp4
	slog.Info("Removing merged file", slog.String("path", mergedFile))
	err = os.Remove(mergedFile)
	if err != nil {
		vc.logError(*task, "failed to remove merged file", err)
		return err
	}

	return nil

}

func (vc *VideoConverter) logError(task VideoTask, message string, err error) {
	errorData := map[string]any { // Eu uso o any porque a chave pode ser de qualquer tipo
		"video_id": task.VideoID,
		"error": 	message,
		"datails": 	err.Error(),
		"time": 	time.Now(),
	}
	serializedError, err := json.Marshal(errorData)
	slog.Error("Processing error", slog.String("error_datails", string(serializedError)))

	RegisterError(vc.db, errorData, err)

}

// Ordenar o caminho do arquivo somente com o nome do arquivo
func (vc *VideoConverter) extractNumber(fileName string) int {
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
func (vc *VideoConverter) mergeChunks(inputDir, outputFile string) error {
	chunks, err := filepath.Glob(filepath.Join(inputDir, "*.chunk")) // Buscar todos os arquivos *.chunk e pegar o caminho completo

	// Tratamento de error
	if err != nil {
		return fmt.Errorf("failed to find chunks: %v", err)
	}

	// Colocar regra de ordenação dos chunks
	sort.Slice(chunks, func(i, j int) bool {
		return vc.extractNumber(chunks[i]) < vc.extractNumber(chunks[j])
	})

	// Criar um arquivo de saida
	output, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	// Depois que toda a função mergeChunks rodar, fecha o arquivo
	defer output.Close()

	// Abrir o chunks (Arquivos)
	for _, chunk := range chunks {
		input, err := os.Open(chunk)
		if err != nil {
			return fmt.Errorf("failed to open chunk: %v", err)
		}

		_, err = output.ReadFrom(input) // Para jogar algo nesse output é preciso ler o input que acabou-se de abrir na linha 53
		if err != nil {
			return fmt.Errorf("failed to write chunk %s to merged file: %v", chunk, err)
		}
		input.Close() // Fechar o chunk (Arquivo)		
	}

	return nil

}