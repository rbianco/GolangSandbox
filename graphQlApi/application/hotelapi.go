package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	hotelsearch "github.com/rbianco/GolangSandbox/graphQlApi/infrastructure/search"

	"github.com/graphql-go/graphql"
)

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

						var hotelFake = hotelsearch.Hotel(idQuery)
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

						var hotelList = hotelsearch.Hotel(idQuery)
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

	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})

	fmt.Println("Now server is running on port 8080")
	fmt.Println("Test with Get      : curl -g 'http://localhost:8080/graphql?query={user(id:\"1\"){name}}'")
	http.ListenAndServe(":8080", nil)
}

//Helper function to import json from file to map
func importJSONDataFromFile(fileName string, result interface{}) (isOK bool) {
	isOK = true
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Print("Error:", err)
		isOK = false
	}
	err = json.Unmarshal(content, result)
	if err != nil {
		isOK = false
		fmt.Print("Error:", err)
	}
	return
}
