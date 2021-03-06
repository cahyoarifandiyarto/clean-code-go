package repository

import (
	"golang-clean-architecture/config"
	"golang-clean-architecture/entity"
	"golang-clean-architecture/exception"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type productRepositoryImpl struct {
	Collection *mongo.Collection
}

func NewProductRepository(database *mongo.Database) ProductRepository {
	return &productRepositoryImpl{
		Collection: database.Collection("products"),
	}
}

func (repository *productRepositoryImpl) Insert(product entity.Product) {
	ctx, cancel := config.NewMongoContext()
	defer cancel()

	_, err := repository.Collection.InsertOne(ctx, bson.M{
		"_id":      product.Id,
		"name":     product.Name,
		"price":    product.Price,
		"quantity": product.Quantity,
	})
	exception.PanicIfNeeded(err)
}

func (repository *productRepositoryImpl) FindAll() (products []entity.Product) {
	ctx, cancel := config.NewMongoContext()
	defer cancel()

	cursor, err := repository.Collection.Find(ctx, bson.M{})
	exception.PanicIfNeeded(err)

	var documents []bson.M
	err = cursor.All(ctx, &documents)
	exception.PanicIfNeeded(err)

	for _, document := range documents {
		products = append(products, entity.Product{
			Id:       document["_id"].(string),
			Name:     document["name"].(string),
			Price:    document["price"].(int),
			Quantity: document["quantity"].(int),
		})
	}

	return products
}

func (repository *productRepositoryImpl) DeleteAll() {
	ctx, cancel := config.NewMongoContext()
	defer cancel()

	_, err := repository.Collection.DeleteMany(ctx, bson.M{})
	exception.PanicIfNeeded(err)
}
