package mongoHelper

import (
	"context"
	"fmt"
	"log"
	"okieoth/schemaguesser/internal/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(conStr string) (*mongo.Client, error) {
	conStr2Use := utils.ReplaceWithEnvVar(conStr, "MONGO_USER", "admin")
	conStr2Use = utils.ReplaceWithEnvVar(conStr2Use, "MONGO_PASSWORD", "secretpassword")

	return mongo.Connect(context.Background(), options.Client().ApplyURI(conStr2Use))
}

func Dummy() {
	fmt.Println("'Dummy' is called")
}

func ListDatabases(conStr string) ([]string, error) {
	var ret []string
	client, err := Connect(conStr)
	if err != nil {
		return ret, err
	}

	defer func() {
		if client == nil {
			return
		}
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	filter := bson.M{}
	cursor, err := client.ListDatabases(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
		return ret, err
	}

	for _, db := range cursor.Databases {
		ret = append(ret, db.Name)
	}
	return ret, nil
}
