package pkg

import (
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/gustavlsouz/documents-service/internal/common"
	"github.com/gustavlsouz/documents-service/internal/common/persistence"
	documentControllers "github.com/gustavlsouz/documents-service/internal/document/controllers"
	documentOperations "github.com/gustavlsouz/documents-service/internal/document/operations"
	"github.com/gustavlsouz/documents-service/internal/middlewares"
	statusControllers "github.com/gustavlsouz/documents-service/internal/status/controllers"
	"github.com/joho/godotenv"
)

func WrapWithGlobalMiddlewares(handler http.Handler) http.Handler {

	return middlewares.NewCorsMiddlewareDecorator(
		middlewares.NewCounterMiddlewareDecorator(handler),
	)
}

func Start(startedSuccessfully chan<- bool, envPath string, migrationPath string) {
	log.Println("Iniciando serviço")
	runtime.GOMAXPROCS(runtime.NumCPU())

	err := godotenv.Load(envPath)
	if err != nil {
		log.Panicf("Erro ao carregar configuração de variáveis de ambiente: %v\n", err)
	}

	err = persistence.GetPersistenceInstance().Start(migrationPath)

	if err != nil {
		log.Panic("Não foi possível iniciar conexão com o banco de dados: ", err)
	}

	mux := http.NewServeMux()

	documentController := documentControllers.NewDocumentController(
		documentOperations.NewDocumentReaderCreator(),
		documentOperations.NewDocumentInserterCreator(),
		documentOperations.NewDocumentRemoverCreator(),
		documentOperations.NewDocumentUpdaterCreator(),
	)

	documentHandler := common.NewHttpHandlerBuilder().
		Post(documentController.Create).
		Get(documentController.Read).
		Put(documentController.Update).
		Delete(documentController.Delete).
		Build()

	mux.Handle("/api/document", documentHandler)

	mux.Handle("/status", common.NewHttpHandlerBuilder().
		Get(statusControllers.NewStatusController().GetStatus).
		Build())

	go func() {
		if err != nil {
			startedSuccessfully <- false
			return
		}
		startedSuccessfully <- true
	}()

	err = http.ListenAndServe(":8080", WrapWithGlobalMiddlewares(mux))
	if err != nil {
		fmt.Printf("Erro ao iniciar o servidor: %v\n", err)
	}
}
