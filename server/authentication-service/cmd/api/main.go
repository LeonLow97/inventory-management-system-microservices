package main

import (
	"log"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"

	"github.com/LeonLow97/internal/adapters/outbound"
	"github.com/LeonLow97/internal/core/services"
	"github.com/LeonLow97/internal/pkg/config"
	"github.com/LeonLow97/internal/pkg/consul"
	"github.com/LeonLow97/internal/pkg/db"
	"github.com/LeonLow97/internal/pkg/grpcserver"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed to load config with error:", err)
	}

	conn, err := db.ConnectToDB(*cfg)
	if err != nil {
		log.Fatalln("Failed to connect to database with error:", err)
	}

	DevTestData(conn)

	// initialize grpc authentication and user services
	repo := outbound.NewRepository(conn)
	service := services.NewService(repo, *cfg)

	app := grpcserver.Application{
		Config:  *cfg,
		Service: service,
	}

	go app.InitiateGRPCServer()

	// register authentication microservice with service discovery consul
	serviceDiscovery := consul.NewConsul(*cfg)
	if err := serviceDiscovery.RegisterService(); err != nil {
		log.Fatalf("failed to register authentication microservice with error: %v\n", err)
	}

	select {}
}

func DevTestData(conn *sqlx.DB) {
	conn.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id              BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
		email           TEXT                     NOT NULL,
		hashed_password TEXT                     NOT NULL,
		first_name      TEXT                     NOT NULL,
		last_name       TEXT                     NOT NULL,
		active          BOOLEAN                  NOT NULL DEFAULT TRUE,
		updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), 
		created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);

	CREATE UNIQUE INDEX idx_users_email ON users(email);

	CREATE TABLE IF NOT EXISTS admin_users (
		user_id         BIGINT                   NOT NULL REFERENCES users(id),
		active          BOOLEAN                  NOT NULL DEFAULT TRUE,     
		updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), 
		created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);

	CREATE UNIQUE INDEX idx_admin_users_user_id ON admin_users(user_id);

	INSERT INTO users (first_name, last_name, hashed_password, email)
	VALUES
		('Jie Wei', 'Low', '$2a$10$OULOXURo57bo5keyNXGQxefqMyEM67JIscqLVKWgd/S.siCqNAHdC', 'lowjiewei@email.com'),
		('Leon', 'Low', '$2a$10$OULOXURo57bo5keyNXGQxefqMyEM67JIscqLVKWgd/S.siCqNAHdC', 'leonlow@email.com'),
	    ('John', 'Doe', '$2a$10$OULOXURo57bo5keyNXGQxefqMyEM67JIscqLVKWgd/S.siCqNAHdC', 'john.doe@example.com'),
		('Jane', 'Smith', '$2a$10$OULOXURo57bo5keyNXGQxefqMyEM67JIscqLVKWgd/S.siCqNAHdC', 'jane.smith@example.com'),
		('Michael', 'Johnson', '$2a$10$3pW8P9U8TeWZCu66QdI.hBGfxh28f/Mq4VJH4HJWmV.qUp3UGewu2', 'michael.johnson@example.com'),
		('Emily', 'Williams', '$2a$10$R1bAS2TYZHz1lgdVndc7hr2lJhbY7kbXyy8KvS6RtUowQke/f5b9a', 'emily.williams@example.com'),
		('David', 'Brown', '$2a$10$WuJ9Pm4WuvlB8Xg6gt5p9pmK2MO5Imdjngg1iz8k9IXvP9Zklj9cG', 'david.brown@example.com'),
		('Sarah', 'Jones', '$2a$10$REfNoY2y.0p1hIgeV/JHfe7z2X32BRh7XIH8JOPjrKHTe67dZkQ.O', 'sarah.jones@example.com'),
		('James', 'Miller', '$2a$10$9bMLnJ0D4qE0LsAvrKpC5KQkO2aqsrmYSZ3QpeHNOA4R0UqCfq1yq', 'james.miller@example.com'),
		('Amanda', 'Davis', '$2a$10$yV0OSbSc5bA2GpZ2DRPbfqCBH6P2dVkq2K6umrktV4Yk5HOy67a8C', 'amanda.davis@example.com'),
		('William', 'Garcia', '$2a$10$Wq9uvL.TuZCGPbLdbCdrs.GTzVDBbGxt1yLlOYKTHfrmjF8HfED7O', 'william.garcia@example.com'),
		('Elizabeth', 'Martinez', '$2a$10$AWTjeLRHE5RgFcQ3KK5eK0t8pgnYYyALwOWJH4NYxx2I9LKOe5Ayy', 'elizabeth.martinez@example.com'),
		('Robert', 'Hernandez', '$2a$10$P2VsEvOss32OFfnZTcUG1CJY79swQ8rHovFgHLdXs1J6BWsC4bzm6', 'robert.hernandez@example.com'),
		('Jessica', 'Lopez', '$2a$10$K3wIq1xH2pTEY7fzybhYPFz0ZRYBfgQLwzmtrfOdZTt7ctZ0yPjo.', 'jessica.lopez@example.com'),
		('Daniel', 'Gonzalez', '$2a$10$KPfQtvBExRz5miVdzJ4fNj4OSXn9m7q47sRjtoP0.1mbA6X8N2oa2', 'daniel.gonzalez@example.com'),
		('Olivia', 'Wilson', '$2a$10$UwABT6Q9XTKmrxGGIltTPJ2szX2Ln5gdRkvL4fU0sujbTZ.OBhhzC', 'olivia.wilson@example.com'),
		('Ethan', 'Anderson', '$2a$10$G4Wdg3Xik5npRll.LsoV2h3E3g6LtthClHzB1twPgyv.M93khd5ju', 'ethan.anderson@example.com'),
		('Avery', 'Thomas', '$2a$10$IRZ0KN8AHTShL/Z8s8gck8PCX4nLmf2I52Ak9wJxEqN63D0Jpr6rV', 'avery.thomas@example.com'),
		('Mason', 'Taylor', '$2a$10$akWV3kB4cLrFQjcQmU2J5KeFqJ9vQJvZI1OH1uSzXjx2v1zTrReV2', 'mason.taylor@example.com'),
		('Chloe', 'Moore', '$2a$10$BpzBd27DUCqgDfaEfoQvsWkzYXjxUHT9uhgvERbr2F2VguD5d11QG', 'chloe.moore@example.com'),
		('Lucas', 'Jackson', '$2a$10$hm5nZTc25B6fgYM2bRu6MbqUzDdsuXj9XDPKAI1go8wI6PRAakczm', 'lucas.jackson@example.com'),
		('Sophia', 'White', '$2a$10$qaOaLTmlS8FEVhSrtgZ.xBz06hKm6YNfcu9RY9v2iV9XMzuTtFGmQ', 'sophia.white@example.com'),
		('Mia', 'Harris', '$2a$10$9VlaoXovdWV7my6brdlF75nN4fWGT7KmT3DpmXckc9IqUywN7J90C', 'mia.harris@example.com'),
		('Jackson', 'Martin', '$2a$10$J0fpd8IPrfvxKOEK59qg..vqWjXt1NEklgFZqEq6hyP56VKtJeeCC', 'jackson.martin@example.com'),
		('Liam', 'Thompson', '$2a$10$z6brsGH0xu6Te9k9rwrn8n9hMbA7RHLx5wuy1zn3ZpJl1kdyZqfGe', 'liam.thompson@example.com'),
		('Evelyn', 'Garcia', '$2a$10$XBVsIxfw5P6yHgB0At4sRrrV0XieM5m3SYE71skFvQTqfKZmXhHgG', 'evelyn.garcia@example.com'),
		('Henry', 'Martinez', '$2a$10$0B8WcBBsq9TpwBZ9VGtq8R3Hk2tNaxnmg9X2Q7l5jXVVtZo92K6SG', 'henry.martinez@example.com'),
		('Aiden', 'Roberts', '$2a$10$ix69Yc7YPK4gBdQ/9B2IE5TUnn0v9BeGB7UzHpvDrlMFDOg1zEnme', 'aiden.roberts@example.com'),
		('Charlotte', 'King', '$2a$10$8uZz8LfHLdUP.fQ0qdeK5pUG8Bd/gujjQzxQBOB5B7F9QFjxuIt5y', 'charlotte.king@example.com'),
		('Benjamin', 'Scott', '$2a$10$kRHbFlbd4X7mc0g9UJ1gkVt0OD8nmRmP4osYXJ/q4wHfzUmixZqMe', 'benjamin.scott@example.com'),
		('Grace', 'Young', '$2a$10$1wYsSLSfzqS0jV74O2apmWlwctdQFCzOHfgcMvqSKmnBglOYqka2C', 'grace.young@example.com'),
		('Elijah', 'Allen', '$2a$10$ktJp6t9ZPSkImWn6pGpW6tmvttdHe7aY2dce9DQowQUpwPOsV1Mqa', 'elijah.allen@example.com');
	
	INSERT INTO admin_users (user_id) VALUES (1);
	`)

	log.Println("Inserted dev data!")
}
