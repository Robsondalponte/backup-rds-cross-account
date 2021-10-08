package main

import (
	"fmt"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

func main() {
	const instance = "<nome-da-instancia>"
	currentTime := time.Now()
	time := currentTime.Format("2006-01-02")
	reg, _ := regexp.Compile("[^0-9]+")
	stringTime := reg.ReplaceAllString(time, "")
	snapName := string(instance + "V" + stringTime)

	fmt.Printf("\niniciando sessão na conta principal\n")
	sess, err := session.NewSession(&aws.Config{
		//Em <regiao>, adicione a região onde se encontra o seu banco de dados. Ex.: sa-east-1;
		Region: aws.String("<regiao>")},
	)
	if err != nil {
		sendMessegeError("Problema ao iniciar sessão %v\n", err)
	}

	client := rds.New(sess)
	fmt.Printf("\nCriando snapshot: %v\n", snapName)
	_, err = client.CreateDBSnapshot(&rds.CreateDBSnapshotInput{
		DBInstanceIdentifier: aws.String(instance),
		DBSnapshotIdentifier: aws.String(snapName),
	})
	if err != nil {
		sendMessegeError("Não habilitado a criar snapshots da instancia. Retorno aws: %v\n", err)
	}

	fmt.Printf("Aguardando o snapshot ser criado...\n")

	err = client.WaitUntilDBSnapshotAvailable(&rds.DescribeDBSnapshotsInput{
		DBInstanceIdentifier: aws.String(instance),
	})
	if err != nil {
		sendMessegeError("Ocorreu um erro enquanto aguardava a criação do snapshot da instancia. Retorno aws: %v\n", err)
	}

	fmt.Printf("Snapshot %q criado com sucesso\n", instance)
	fmt.Printf("Iniciando processo de compartilhamento\n")
	share(client, snapName)

	fmt.Printf("Persistindo os dados na conta cofre\n")
	persist(snapName, instance)
}
