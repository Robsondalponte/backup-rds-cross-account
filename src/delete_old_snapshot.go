package main

import (
	"fmt"

	
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/rds"
)

func deleteOldSnapshot (client *rds.RDS, name string){
	input := &rds.DeleteDBSnapshotInput{
		DBSnapshotIdentifier: aws.String(name),
	}
	result, err := client.DeleteDBSnapshot(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rds.ErrCodeInvalidDBSnapshotStateFault:
				sendMessegeError(rds.ErrCodeInvalidDBSnapshotStateFault, aerr.Error())
			case rds.ErrCodeDBSnapshotNotFoundFault:
				sendMessegeError(rds.ErrCodeDBSnapshotNotFoundFault, aerr.Error())
			default:
				sendMessegeError("Erro inesperado. Retorno aws: ", aerr.Error())
			}
		} else {
			sendMessegeError("Erro inesperado: ", err.Error())
		}
		return
	}
	if aws.StringValue(result.DBSnapshot.Status) == "deleted"{
		fmt.Printf("Snapshot %s deletado com sucesso", aws.StringValue(result.DBSnapshot.DBSnapshotIdentifier))
	} else {
		sendMessegeError("Novo erro identificado: ", NonDelete)
	}
}
