package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rbianco/GolangSandbox/graphQlApi/application/config"
	"github.com/rbianco/GolangSandbox/graphQlApi/infrastructure/hotelsearch"
	"github.com/spf13/viper"

	"github.com/graphql-go/graphql"
)

var configuration config.Configuration

var hotelType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Hotel",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"uri": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"hotel": &graphql.Field{
				Type: hotelType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					idQuery, isOK := p.Args["id"].(string)

					if isOK {

						var hotelFake = hotelsearch.Hotel(configuration.SearchServiceEndPoint, idQuery)
						return hotelFake[0], nil
					}
					return nil, nil
				},
			},
			"hotelList": &graphql.Field{
				Type:        graphql.NewList(hotelType),
				Description: "Hotel search service",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					idQuery, isOK := p.Args["id"].(string)

					if isOK {

						var hotelList = hotelsearch.Hotel(configuration.SearchServiceEndPoint, idQuery)
						return hotelList, nil
					}
					return nil, nil
				},
			},
		},
	})

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: queryType,
	},
)

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

func main() {

	initConfig()
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})

	fmt.Println("Now server is running on port 8080")
	fmt.Println("Test with Get      : curl -g 'http://localhost:8080/graphql?query={hotelList(id:\"55688\"){id,name,title,uri}}'")
	http.ListenAndServe(":8080", nil)
}

func initConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error al leer del archivo config, %s", err)
	}
	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("Error al deserializar la configuración, %s", err)
	}
}
