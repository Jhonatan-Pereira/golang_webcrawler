package db

import (
	"context"
)

func Insert(colletion string, data interface{}) error {
	client, ctx := getConnection()
	defer client.Disconnect(ctx)

	c := client.Database("crawler").Collection(colletion)

	_, err := c.InsertOne(context.Background(), data)

	return err
}
