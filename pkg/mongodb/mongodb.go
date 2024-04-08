package mongodb

import (
    "context"
    "errors"
    "fmt"
    "wait4it/pkg/model"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
)

func (m *MongoDbConnection) buildConnectionString() string {
    // MongoDB connection string format
    return fmt.Sprintf("mongodb://%s:%s@%s:%d", m.Username, m.Password, m.Host, m.Port)
}

func (m *MongoDbConnection) BuildContext(cx model.CheckContext) {
    m.Port = cx.Port
    m.Host = cx.Host
    m.Username = cx.Username
    m.Password = cx.Password
}

func (m *MongoDbConnection) Validate() error {
    if m.Host == "" {
        return errors.New("host can't be empty")
    }

    if m.Username != "" && m.Password == "" {
        return errors.New("password cannot be empty if a username is provided")
    }

    if m.Port < 1 || m.Port > 65535 {
        return errors.New("invalid port range for MongoDB")
    }

    return nil
}

func (m *MongoDbConnection) Check(ctx context.Context) (bool, bool, error) {
    client, err := mongo.NewClient(options.Client().ApplyURI(m.buildConnectionString()))
    if err != nil {
        return false, true, err
    }

    err = client.Connect(ctx)
    if err != nil {
        return false, true, err
    }

    err = client.Ping(ctx, readpref.Primary())
    if err != nil {
        return false, false, err
    }

    return true, true, nil
}
