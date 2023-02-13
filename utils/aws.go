package utils

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func GetPathStoreParameters(path string) string {
	// Inicializa a sessão do AWS SDK
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Cria uma nova instância do serviço SSM
	ssmSvc := ssm.New(sess)

	// Busca o parâmetro específico
	param, err := ssmSvc.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(path),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		fmt.Println("Erro ao buscar o parâmetro:", err)
		return ""
	}

	// Imprime o valor do parâmetro
	fmt.Println("Parametro ", path, " encontrado...")
	var result = *param.Parameter.Value
	return result
}
