# Backup cross account na AWS para banco de dados RDS 

<img src="https://miro.medium.com/max/1400/1*CdjOgfolLt_GNJYBzI-1QQ.jpeg" width="300"/>

Projeto criado em golang. Veja a documentação [aqui](https://golang.org/doc/)

---


## Detalhes sobre o projeto

Com esse projeto é possível automatizar o processo de backup da instancia RDS para outra conta de destino (uma conta cofre).



O script roda em um container docker. A imagem desse container é composta pela imagem para golang e features para funções de scheduler. Veja detalhes da imagem [aqui](https://hub.docker.com/r/rodrigodiez/golang-cron/), ou acesse o repositorio no github [aqui](https://github.com/rodrigodiez/golang-cron).


**Obs.: Atente para o fato da imagem original não ter um UTC definido. Ou seja, existe uma diferença de 3h para o horario de Brasilia (por exemplo: quando aqui for 3h, na imagem será 6h)**

---


## Como rodar o projeto localmente

* Primeiro instale o golang utilizando a documentação oficial.
* Faça Download do projeto: `git clone https://github.com/Robsondalponte/backup-rds-cross-account.git`
* Acesse este diretorio: `cd backup-rds-cross-account/src`
* Criei o arquivo de rastreamento de modulos (go.mod): `go mod init infra.retta.com.br/go`
* Instale as dependencias do projeto: `go install -v`
* Veja se o arquivo binario foi criado. Se isso aconteceu, rode o projeto: `go run .`

**Obs.: Caso já exista um snapshot (manual) na conta cofre ou na conta principal, ocorrera um erro. Para testar, limpe os snapshots manuais no painel do RDS**

---


## Dependências

Estão sendo utilizadas algumas dependências externas para este projeto:

#### aws-sdk-go

Utilizamos o SDK da aws para Golang
##### [aws](https://github.com/aws/aws-sdk-go/tree/main/aws)
##### [awserr](https://github.com/aws/aws-sdk-go/tree/main/aws/awserr)
##### [rds](https://github.com/aws/aws-sdk-go/tree/main/service/rds)
##### [session](https://github.com/aws/aws-sdk-go/tree/main/aws/session)
##### [credentials](https://github.com/aws/aws-sdk-go/tree/main/aws/credentials)


---


## Credenciais

Foram utilizadas duas credenciais IAM (uma para cada conta). Com as seguintes regras (formato json):

`
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": [
                "rds:DescribeDBSnapshots",
                "rds:CreateDBSnapshot",
                "rds:CopyDBSnapshot",
                "rds:DeleteDBSnapshot",
                "rds:DescribeDBSnapshotAttributes"
            ],
            "Resource": "*",
            "Condition": {
                "IpAddress": {
                    "aws:SourceIp": "<ip addres>"
                }
            }
        }
    ]
}
`

---

<br>
<br>
<br>
<br>
<img src="https://miro.medium.com/max/1400/1*CdjOgfolLt_GNJYBzI-1QQ.jpeg" width="150"/> <img src="https://cdn.pixabay.com/photo/2016/12/21/17/11/signe-1923369_960_720.png" width="20"/> <img src="https://upload.wikimedia.org/wikipedia/commons/thumb/9/93/Amazon_Web_Services_Logo.svg/512px-Amazon_Web_Services_Logo.svg.png" width="150"/>

