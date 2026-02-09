package gql

import (
	"log"
	"time"

	"pg-management-system/internal/database"
	"pg-management-system/internal/models"

	"github.com/graphql-go/graphql"
)

// Define Types
// Define Types
var roomType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Room",
	Fields: graphql.Fields{
		"id":          &graphql.Field{Type: graphql.Int},
		"room_number": &graphql.Field{Type: graphql.String},
		"capacity":    &graphql.Field{Type: graphql.Int},
		"occupancy":   &graphql.Field{Type: graphql.Int},
		"price":       &graphql.Field{Type: graphql.Float},
	},
})

var guestType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Guest",
	Fields: graphql.Fields{
		"id":        &graphql.Field{Type: graphql.Int},
		"name":      &graphql.Field{Type: graphql.String},
		"email":     &graphql.Field{Type: graphql.String},
		"phone":     &graphql.Field{Type: graphql.String},
		"room_id":   &graphql.Field{Type: graphql.Int},
		"join_date": &graphql.Field{Type: graphql.String}, // Simplified as string for simplicity
	},
})

var paymentType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Payment",
	Fields: graphql.Fields{
		"id":             &graphql.Field{Type: graphql.Int},
		"guest_id":       &graphql.Field{Type: graphql.Int},
		"amount":         &graphql.Field{Type: graphql.Float},
		"payment_date":   &graphql.Field{Type: graphql.String},
		"payment_method": &graphql.Field{Type: graphql.String},
	},
})

// Define Root Query
var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"rooms": &graphql.Field{
			Type: graphql.NewList(roomType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return database.GetAllRooms()
			},
		},
		"room": &graphql.Field{
			Type: roomType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.Int},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, _ := p.Args["id"].(int)
				return database.GetRoomByID(id)
			},
		},
		"guests": &graphql.Field{
			Type: graphql.NewList(guestType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return database.GetAllGuests()
			},
		},
		"guest": &graphql.Field{
			Type: guestType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.Int},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, _ := p.Args["id"].(int)
				return database.GetGuestByID(id)
			},
		},
		"allPayments": &graphql.Field{
			Type: graphql.NewList(paymentType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return database.GetAllPayments()
			},
		},
		"payment": &graphql.Field{
			Type: paymentType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, _ := p.Args["id"].(int)
				return database.GetPaymentByID(id)
			},
		},
		"payments": &graphql.Field{
			Type: graphql.NewList(paymentType),
			Args: graphql.FieldConfigArgument{
				"guest_id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				guestID, _ := p.Args["guest_id"].(int)
				return database.GetPaymentsByGuestID(guestID)
			},
		},
	},
})

// Define Mutation
var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		// Room Mutations
		"createRoom": &graphql.Field{
			Type: roomType,
			Args: graphql.FieldConfigArgument{
				"room_number": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"capacity":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"price":       &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Float)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				room := models.Room{
					RoomNumber: p.Args["room_number"].(string),
					Capacity:   p.Args["capacity"].(int),
					Price:      p.Args["price"].(float64),
				}
				err := database.CreateRoom(&room)
				if err != nil {
					return nil, err
				}
				return room, nil
			},
		},
		"updateRoom": &graphql.Field{
			Type: roomType,
			Args: graphql.FieldConfigArgument{
				"id":          &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"room_number": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"capacity":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"occupancy":   &graphql.ArgumentConfig{Type: graphql.Int}, // Optional
				"price":       &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Float)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Args["id"].(int)
				room := models.Room{
					RoomNumber: p.Args["room_number"].(string),
					Capacity:   p.Args["capacity"].(int),
					Price:      p.Args["price"].(float64),
				}
				if val, ok := p.Args["occupancy"].(int); ok {
					room.Occupancy = val
				}
				err := database.UpdateRoom(id, &room)
				if err != nil {
					return nil, err
				}
				room.ID = id
				return room, nil
			},
		},
		"deleteRoom": &graphql.Field{
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Args["id"].(int)
				err := database.DeleteRoom(id)
				if err != nil {
					return false, err
				}
				return true, nil
			},
		},

		// Guest Mutations
		"createGuest": &graphql.Field{
			Type: guestType,
			Args: graphql.FieldConfigArgument{
				"name":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"email":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"phone":   &graphql.ArgumentConfig{Type: graphql.String},
				"room_id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				guest := models.Guest{
					Name:     p.Args["name"].(string),
					Email:    p.Args["email"].(string),
					Phone:    p.Args["phone"].(string),
					RoomID:   p.Args["room_id"].(int),
					JoinDate: time.Now(),
				}
				err := database.CreateGuest(&guest)
				if err != nil {
					return nil, err
				}
				return guest, nil
			},
		},
		"updateGuest": &graphql.Field{
			Type: guestType,
			Args: graphql.FieldConfigArgument{
				"id":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"name":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"email":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"phone":   &graphql.ArgumentConfig{Type: graphql.String},
				"room_id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Args["id"].(int)
				guest := models.Guest{
					Name:   p.Args["name"].(string),
					Email:  p.Args["email"].(string),
					Phone:  p.Args["phone"].(string),
					RoomID: p.Args["room_id"].(int),
				}
				err := database.UpdateGuest(id, &guest)
				if err != nil {
					return nil, err
				}
				guest.ID = id
				return guest, nil
			},
		},
		"deleteGuest": &graphql.Field{
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Args["id"].(int)
				err := database.DeleteGuest(id)
				if err != nil {
					return false, err
				}
				return true, nil
			},
		},

		// Payment Mutation
		"createPayment": &graphql.Field{
			Type: paymentType,
			Args: graphql.FieldConfigArgument{
				"guest_id":       &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"amount":         &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Float)},
				"payment_method": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				payment := models.Payment{
					GuestID:       p.Args["guest_id"].(int),
					Amount:        p.Args["amount"].(float64),
					PaymentMethod: p.Args["payment_method"].(string),
					PaymentDate:   time.Now(),
				}
				err := database.CreatePayment(&payment)
				if err != nil {
					return nil, err
				}
				return payment, nil
			},
		},
		"updatePayment": &graphql.Field{
			Type: paymentType,
			Args: graphql.FieldConfigArgument{
				"id":             &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"guest_id":       &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
				"amount":         &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Float)},
				"payment_method": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Args["id"].(int)
				payment := models.Payment{
					ID:            id,
					GuestID:       p.Args["guest_id"].(int),
					Amount:        p.Args["amount"].(float64),
					PaymentMethod: p.Args["payment_method"].(string),
					PaymentDate:   time.Now(), // Or fetch current and keep it, but Repository uses this.
				}
				err := database.UpdatePayment(&payment)
				if err != nil {
					return nil, err
				}
				return payment, nil
			},
		},
		"deletePayment": &graphql.Field{
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Args["id"].(int)
				err := database.DeletePayment(id)
				if err != nil {
					return false, err
				}
				return true, nil
			},
		},
	},
})

// Schema
var Schema graphql.Schema

func init() {
	var err error
	Schema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})
	if err != nil {
		log.Fatalf("Failed to create GraphQL schema: %v", err)
	}
}
