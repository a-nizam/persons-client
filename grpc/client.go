package grpcclient

import (
	"context"
	"fmt"
	"time"

	pb "github.com/a-nizam/persons-client/gen"
	"github.com/a-nizam/persons-client/models"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	nullValue = 0
)

type Client struct {
	grpcclient pb.PersonsClient
}

func New() *Client {
	conn, err := grpc.NewClient(":40404", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("Failed to connect grpc server")
	}
	client := Client{grpcclient: pb.NewPersonsClient(conn)}
	return &client
}

func (c *Client) AddPerson(ctx context.Context, person *models.Person) (int64, error) {
	personID, err := c.grpcclient.AddPerson(ctx, &pb.Person{
		ID:        person.ID,
		Name:      person.Name,
		Birthdate: person.Birthdate.Format("2006-01-02"),
	})
	if err != nil {
		fmt.Printf("Failed to add person: %v", err)
		return nullValue, err
	}
	return personID.Value, nil
}

func (c *Client) GetPerson(ctx context.Context, id int64) (*models.Person, error) {
	person, err := c.grpcclient.GetPerson(ctx, &pb.PersonID{Value: id})
	if err != nil {
		fmt.Printf("Failed to get person: %v", err)
		return nil, err
	}
	birthdate, err := time.Parse("2006-01-02", person.Birthdate)
	if err != nil {
		fmt.Printf("Failed to get person: %v", err)
		return nil, err
	}
	return &models.Person{
		ID:        person.ID,
		Name:      person.Name,
		Birthdate: birthdate,
	}, nil
}

func (c *Client) EditPerson(ctx context.Context, person *models.Person) error {
	_, err := c.grpcclient.EditPerson(ctx, &pb.Person{
		ID:        person.ID,
		Name:      person.Name,
		Birthdate: person.Birthdate.Format("2006-01-02"),
	})
	if err != nil {
		fmt.Printf("Failed to edit person: %v", err)
		return err
	}
	return nil
}

func (c *Client) RemovePerson(ctx context.Context, id int64) error {
	_, err := c.grpcclient.RemovePerson(ctx, &pb.PersonID{Value: id})
	if err != nil {
		fmt.Printf("Failed to remove person: %v", err)
		return err
	}
	return nil
}

func (c *Client) GetList(ctx context.Context) (grpc.ServerStreamingClient[pb.Person], error) {
	personList, err := c.grpcclient.GetList(ctx, &pb.Empty{})
	if err != nil {
		fmt.Printf("Failed to get person list: %v", err)
		return nil, err
	}
	return personList, nil
}
