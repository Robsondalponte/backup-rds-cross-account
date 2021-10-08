package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
)

var (
	ManyInstances = errors.New("Foram encontrados mais de um snapshots. Delete manualmente para que o fluxo possa seguir")
	NoneInstance  = errors.New("Nenhum snapshot recente foi encontrado")
	NonDelete     = errors.New("Algo aconteceu e não foi possível deletar o snapshot antigo")
)

func modifySnapshot(client *rds.RDS, name string) {
	input := &rds.ModifyDBSnapshotAttributeInput{
		AttributeName:        aws.String("restore"),
		DBSnapshotIdentifier: aws.String(name),
		ValuesToAdd: []*string{
			//substitua o valor pelo key ID da sua conta de cofre
			aws.String("9999999999"),
		},
	}

	fmt.Println("Compartilhando...")
	_, err := client.ModifyDBSnapshotAttribute(input)
	if err != nil {
		sendMessegeError("Erro ao modificar a propriedade do snpashot para permitir compartilhamento. Retorno aws: %s\n", err)
	}
}
func share(client *rds.RDS, snapName string) {
	dayBeforeYesterday := time.Now().Add(-48 * time.Hour)
	result, err := client.DescribeDBSnapshots(nil)
	if err != nil {
		sendMessegeError("Não habilitado para listar snapshots, %s\n", err)
	}
	for _, s := range result.DBSnapshots {
		if strings.EqualFold(aws.StringValue(s.SnapshotType), "manual") {
			if strings.EqualFold(aws.StringValue(s.DBSnapshotIdentifier), snapName) {
				fmt.Println(snapName)
				modifySnapshot(client, snapName)
			} else if dayBeforeYesterday.After(aws.TimeValue(s.SnapshotCreateTime)) {
				deleteOldSnapshot(client, aws.StringValue(s.DBSnapshotIdentifier))
			}
		}
	}
}
