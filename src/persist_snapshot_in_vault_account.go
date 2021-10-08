package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

func persist(snapName string, instance string) {
	dayBeforeYesterday := time.Now().Add(-48 * time.Hour)
	identifierSnapshotShared := "arn:aws:rds:<REGION>:<ACCOUNT ID>:snapshot:" + snapName
	copySnapshotName := snapName + "copy"

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("<REGION>"),
		//Credentials: credentials.NewSharedCredentials("", "vault-account"),
		Credentials: credentials.NewStaticCredentials("<YOUR-KEY-ID>", "<YOUR-SECERT-KEY>", "")},
	)
	if err != nil {
		sendMessegeError("Problema com as configurações de acesso a conta cofre:\t%v\n", err)
	}

	client := rds.New(sess)

	_, err = client.CopyDBSnapshot(&rds.CopyDBSnapshotInput{
		SourceDBSnapshotIdentifier: aws.String(identifierSnapshotShared),
		TargetDBSnapshotIdentifier: aws.String(copySnapshotName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rds.ErrCodeDBSnapshotAlreadyExistsFault:
				sendMessegeError(rds.ErrCodeDBSnapshotAlreadyExistsFault, aerr.Error())
			case rds.ErrCodeDBSnapshotNotFoundFault:
				sendMessegeError(rds.ErrCodeDBSnapshotNotFoundFault, aerr.Error())
			case rds.ErrCodeInvalidDBSnapshotStateFault:
				sendMessegeError(rds.ErrCodeInvalidDBSnapshotStateFault, aerr.Error())
			case rds.ErrCodeSnapshotQuotaExceededFault:
				sendMessegeError(rds.ErrCodeSnapshotQuotaExceededFault, aerr.Error())
			case rds.ErrCodeKMSKeyNotAccessibleFault:
				sendMessegeError(rds.ErrCodeKMSKeyNotAccessibleFault, aerr.Error())
			case rds.ErrCodeCustomAvailabilityZoneNotFoundFault:
				sendMessegeError(rds.ErrCodeCustomAvailabilityZoneNotFoundFault, aerr.Error())
			default:
				sendMessegeError("Erro inesperado. Retorno aws: ", aerr.Error())
			}
		} else {
			sendMessegeError("Erro inesperado: ", err.Error())
		}
		return
	}
	fmt.Printf("Aguardando pela copia do snapshot...")
	err = client.WaitUntilDBSnapshotAvailable(&rds.DescribeDBSnapshotsInput{
		DBInstanceIdentifier: aws.String(instance),
	})
	if err != nil {
		sendMessegeError("Ocorreu um erro enquanto aguardava a criação do snapshot.", err)
	}
	fmt.Printf("Snapshot foi copiado com sucesso para conta de destino\n")

	result, err := client.DescribeDBSnapshots(nil)
	if err != nil {
		sendMessegeError("Não habilitado para listar snapshots, %s\n", err)
	}
	fmt.Printf("Deletando snapshots com mais de 48h\n")
	for _, s := range result.DBSnapshots {
		if aws.StringValue(s.SnapshotType) == "manual" && dayBeforeYesterday.After(aws.TimeValue(s.SnapshotCreateTime)) {
			deleteOldSnapshot(client, aws.StringValue(s.DBSnapshotIdentifier))
		}
	}
}
