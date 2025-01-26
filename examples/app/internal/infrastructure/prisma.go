package prisma

import (
	"log"

	"github.com/AbM7797/prisma-query-parser/examples/app/prisma/db"
)

func InitPrismaClient() *db.PrismaClient {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		log.Fatalf("Failed to connect to Prisma: %v", err)
	}
	return client
}
