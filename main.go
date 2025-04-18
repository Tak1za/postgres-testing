package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
	"github.com/jackc/pgx/v5"
)

type Point struct {
	Id         string           `faker:"uuid_hyphenated" json:"id"`
	Name       string           `faker:"name" json:"name"`
	ExternalId string           `faker:"name" json:"externalId"`
	TenantId   string           `faker:"tenantId" json:"tenantId"`
	Attributes []PointAttribute `json:"attributes"`
}

type PointJsonb struct {
	Id       string      `faker:"uuid_hyphenated" json:"id"`
	TenantId string      `faker:"tenantId" json:"tenantId"`
	Config   PointConfig `faker:"config" json:"config"`
}

type PointAttribute struct {
	Id    string `faker:"uuid_hyphenated" json:"id"`
	Name  string `faker:"name" json:"name"`
	Value string `faker:"name" json:"value"`
}

type PointConfig struct {
	Id         string           `faker:"uuid_hyphenated" json:"id"`
	Name       string           `faker:"name" json:"name"`
	ExternalId string           `faker:"name" json:"externalId"`
	Attributes []PointAttribute `faker:"attributes" json:"attributes"`
}

func main() {
	r := gin.Default()

	r.GET("/fetch/:tenantId", func(c *gin.Context) {
		tenantId := c.Param("tenantId")
		conn, err := pgx.Connect(context.Background(), "postgres://postgres:varunak47@localhost:5432/psr-testing")
		if err != nil {
			log.Printf("Unable to connect to database: %v\n", err)
			os.Exit(1)
		}
		defer conn.Close(context.Background())

		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		if err != nil {
			log.Printf("Unable to start transaction: %v\n", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		start := time.Now()
		points, err := fetchPoints(tx, tenantId)
		if err != nil {
			log.Printf("Unable to fetch points: %v\n", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		if err := tx.Commit(context.Background()); err != nil {
			log.Printf("Unable to commit transaction: %v\n", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		elapsed := time.Since(start)
		log.Printf("took %s", elapsed)

		c.JSON(http.StatusOK, gin.H{
			"data": points,
		})
	})

	r.GET("/fetchjsonb/:tenantId", func(c *gin.Context) {
		tenantId := c.Param("tenantId")
		conn, err := pgx.Connect(context.Background(), "postgres://postgres:varunak47@localhost:5432/psr-testing")
		if err != nil {
			log.Printf("Unable to connect to database: %v\n", err)
			os.Exit(1)
		}
		defer conn.Close(context.Background())

		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		if err != nil {
			log.Printf("Unable to start transaction: %v\n", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		start := time.Now()
		points, err := fetchPointsJsonb(tx, tenantId)
		if err != nil {
			log.Printf("Unable to fetch points: %v\n", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		if err := tx.Commit(context.Background()); err != nil {
			log.Printf("Unable to commit transaction: %v\n", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		elapsed := time.Since(start)
		log.Printf("took %s", elapsed)

		c.JSON(http.StatusOK, gin.H{
			"data": points,
		})
	})

	r.GET("/fetchjsonbwithattributes/:tenantId", func(c *gin.Context) {
		tenantId := c.Param("tenantId")
		conn, err := pgx.Connect(context.Background(), "postgres://postgres:varunak47@localhost:5432/psr-testing")
		if err != nil {
			log.Printf("Unable to connect to database: %v\n", err)
			os.Exit(1)
		}
		defer conn.Close(context.Background())

		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		if err != nil {
			log.Printf("Unable to start transaction: %v\n", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		start := time.Now()
		points, err := fetchPointsJsonbWithAttributes(tx, tenantId)
		if err != nil {
			log.Printf("Unable to fetch points: %v\n", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		if err := tx.Commit(context.Background()); err != nil {
			log.Printf("Unable to commit transaction: %v\n", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		elapsed := time.Since(start)
		log.Printf("took %s", elapsed)

		c.JSON(http.StatusOK, gin.H{
			"data": points,
		})
	})

	r.POST("/insert", func(c *gin.Context) {
		conn, err := pgx.Connect(context.Background(), "postgres://postgres:varunak47@localhost:5432/psr-testing")
		if err != nil {
			log.Printf("Unable to connect to database: %v\n", err)
			os.Exit(1)
		}
		defer conn.Close(context.Background())

		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		if err != nil {
			log.Printf("Unable to start transaction: %v\n", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		start := time.Now()

		tenantsToGenerate := 25
		pointsToGenerate := 100000
		points := make([]Point, 0)
		for i := 0; i < tenantsToGenerate; i++ {
			faker.AddProvider("tenantId", func(v reflect.Value) (interface{}, error) {
				return string(fmt.Sprintf("tenant%d", i+1)), nil
			})
			for j := 0; j < pointsToGenerate; j++ {
				p := generatePoint()
				points = append(points, p)
			}
			faker.RemoveProvider("tenantId")
		}

		if err := addPoints(tx, points); err != nil {
			log.Printf("Unable to insert point: %v\n", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		if err := tx.Commit(context.Background()); err != nil {
			log.Printf("Unable to commit transaction: %v\n", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		elapsed := time.Since(start)
		log.Printf("took %s", elapsed)

		c.Status(http.StatusOK)
	})

	r.POST("/insertjsonb", func(c *gin.Context) {
		conn, err := pgx.Connect(context.Background(), "postgres://postgres:varunak47@localhost:5432/psr-testing")
		if err != nil {
			log.Printf("Unable to connect to database: %v\n", err)
			os.Exit(1)
		}
		defer conn.Close(context.Background())

		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		if err != nil {
			log.Printf("Unable to start transaction: %v\n", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		start := time.Now()

		tenantsToGenerate := 25
		pointsToGenerate := 100000
		points := make([]PointJsonb, 0)
		faker.AddProvider("attributes", func(v reflect.Value) (interface{}, error) {
			attributes := make([]PointAttribute, 0)
			for i := 0; i < 5; i++ {
				attribute := PointAttribute{}
				_ = faker.FakeData(&attribute)
				attributes = append(attributes, attribute)
			}
			return attributes, nil
		})
		faker.AddProvider("config", func(v reflect.Value) (interface{}, error) {
			pc := PointConfig{}
			_ = faker.FakeData(&pc)
			return pc, nil
		})
		for i := 0; i < tenantsToGenerate; i++ {
			faker.AddProvider("tenantId", func(v reflect.Value) (interface{}, error) {
				return string(fmt.Sprintf("tenant%d", i+1)), nil
			})
			for j := 0; j < pointsToGenerate; j++ {
				p := generateJsonbPoint()
				points = append(points, p)
			}
			faker.RemoveProvider("tenantId")
		}

		if err := addPointsJsonb(tx, points); err != nil {
			log.Printf("Unable to insert point: %v\n", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		if err := tx.Commit(context.Background()); err != nil {
			log.Printf("Unable to commit transaction: %v\n", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		elapsed := time.Since(start)
		log.Printf("took %s", elapsed)

		c.Status(http.StatusOK)
	})

	r.Run()
}

func generatePoint() Point {
	p := Point{}
	_ = faker.FakeData(&p)
	return p
}

func generateJsonbPoint() PointJsonb {
	p := PointJsonb{}

	_ = faker.FakeData(&p)
	return p
}

func addPoints(tx pgx.Tx, point []Point) error {
	_, err := tx.CopyFrom(context.Background(), pgx.Identifier{"points"}, []string{"id", "name", "external_id", "tenant_id"}, pgx.CopyFromSlice(len(point), func(i int) ([]any, error) {
		return []any{point[i].Id, point[i].Name, point[i].ExternalId, point[i].TenantId}, nil
	}))
	if err != nil {
		return nil
	}

	return nil
}

func addPointsJsonb(tx pgx.Tx, point []PointJsonb) error {
	_, err := tx.CopyFrom(context.Background(), pgx.Identifier{"points_jsonb"}, []string{"id", "tenant_id", "config"}, pgx.CopyFromSlice(len(point), func(i int) ([]any, error) {
		return []any{point[i].Id, point[i].TenantId, point[i].Config}, nil
	}))
	if err != nil {
		return nil
	}

	return nil
}

func fetchPoints(tx pgx.Tx, tenantId string) ([]Point, error) {
	rows, err := tx.Query(context.Background(), "select id, name, external_id, tenant_id from points where tenant_id = $1 limit 10000", tenantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	points := make([]Point, 0)
	for rows.Next() {
		p := Point{}
		if err := rows.Scan(&p.Id, &p.Name, &p.ExternalId, &p.TenantId); err != nil {
			return nil, err
		}
		points = append(points, p)
	}

	return points, nil
}

func fetchPointsJsonb(tx pgx.Tx, tenantId string) ([]Point, error) {
	rows, err := tx.Query(context.Background(), "select id, config->>'name', config->>'externalId', tenant_id from points_jsonb where tenant_id = $1 limit 10000", tenantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	points := make([]Point, 0)
	for rows.Next() {
		p := Point{}
		if err := rows.Scan(&p.Id, &p.Name, &p.ExternalId, &p.TenantId); err != nil {
			return nil, err
		}
		points = append(points, p)
	}

	return points, nil
}

func fetchPointsJsonbWithAttributes(tx pgx.Tx, tenantId string) ([]Point, error) {
	rows, err := tx.Query(context.Background(), "select id, config->>'name', config->>'externalId', tenant_id, config->'attributes' from points_jsonb where tenant_id = $1 limit 10000", tenantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	points := make([]Point, 0)
	for rows.Next() {
		p := Point{}
		if err := rows.Scan(&p.Id, &p.Name, &p.ExternalId, &p.TenantId, &p.Attributes); err != nil {
			return nil, err
		}
		points = append(points, p)
	}

	return points, nil
}
