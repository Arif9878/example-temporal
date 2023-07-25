package activity

import (
	"context"

	"github.com/Arif9878/example-temporal/messages"
	"github.com/Arif9878/example-temporal/repository"
)

func SetCart(ctx context.Context, cart *messages.Cart) (err error) {
	c := &repository.Cart{}

	err = c.SetCart(ctx, cart)
	if err != nil {
		return err
	}

	return err
}

func GetCart(ctx context.Context) (data *[]messages.Product, err error) {
	c := &repository.Cart{}

	data, err = c.GetCart(ctx)
	if err != nil {
		return nil, err
	}

	return data, err
}
