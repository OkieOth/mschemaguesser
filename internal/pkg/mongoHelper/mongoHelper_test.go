package mongoHelper

import (
	"context"
	"testing"
)

var conStr = "mongodb://{MONGO_USER}:{MONGO_PASSWORD}@{MONGO_HOST}:{MONGO_PORT}/admin"

func TestListDbs_IT(t *testing.T) {
	client, err := Connect(conStr)
	defer CloseConnection(client)

	if err != nil {
		t.Errorf("Failed to get client: %v", err)
		return
	}
	dbs, err := ListDatabases(client)

	if err != nil {
		t.Errorf("Failed to list databases: %v", err)
		return
	}

	if len(dbs) != 4 {
		t.Errorf("Retrieved wrong number of dbs: %v", dbs)
	}
}

func TestConnect_IT(t *testing.T) {
	client, err := Connect(conStr)

	defer func() {
		if client == nil {
			return
		}
		client.Disconnect(context.Background())
	}()

	if err != nil {
		t.Errorf("Failed to connect to db: %v", err)
	}
}
