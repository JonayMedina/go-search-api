package mongodb

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/JonayMedina/go-search-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	client *mongo.Client
	db     string
}

func NewUserRepository(client *mongo.Client, db string) *UserRepository {
	return &UserRepository{
		client: client,
		db:     db,
	}
}

func (repository *UserRepository) Create(ctx context.Context, user *models.User) error {
	log.Printf("Intentando crear usuario: %+v", user)
	collection := repository.client.Database(repository.db).Collection("users")

	user.ID = primitive.NewObjectID()
	log.Printf("ID generado: %s", user.ID.Hex())

	exists, err := repository.exists(ctx, user.Username, user.Email)
	if err != nil {
		log.Printf("Error verificando existencia: %v", err)
		return err
	}
	if exists {
		log.Println("Usuario ya existe")
		return errors.New("username or email already exists")
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if err := user.HashPassword(); err != nil {
		log.Printf("Error hasheando password: %v", err)
		return err
	}

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		log.Printf("Error insertando en MongoDB: %v", err)
		return err
	}
	log.Printf("Usuario creado con ID: %v", result.InsertedID)
	return nil
}

func (repository *UserRepository) exists(ctx context.Context, username, email string) (bool, error) {
	log.Printf("Verificando existencia de usuario - Username: %s, Email: %s", username, email)

	collection := repository.client.Database(repository.db).Collection("users")

	filter := bson.M{
		"$or": []bson.M{
			{"username": username},
			{"email": email},
		},
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("Error verificando existencia: %v", err)
		return false, err
	}

	exists := count > 0
	log.Printf("Usuario existe: %v", exists)
	return exists, nil
}

func (repository *UserRepository) GetUsers(ctx context.Context) ([]models.User, error) {
	log.Printf("Consultando usuarios en la base de datos %s", repository.db)
	collection := repository.client.Database(repository.db).Collection("users")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error en consulta Find: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		log.Printf("Error decodificando usuarios: %v", err)
		return nil, err
	}

	log.Printf("Usuarios recuperados exitosamente. Total: %d", len(users))
	return users, nil
}

func (repository *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	log.Printf("Consultando usuario por username: %s", username)
	collection := repository.client.Database(repository.db).Collection("users")
	filter := bson.M{"username": username}
	return repository.getUser(ctx, collection, filter)
}

func (repository *UserRepository) getUser(ctx context.Context, collection *mongo.Collection, filter bson.M) (*models.User, error) {
	var user models.User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		log.Printf("Error al obtener usuario: %v", err)
		return nil, err
	}
	return &user, nil
}
