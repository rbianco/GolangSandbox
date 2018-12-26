package main

import (
	"log"

	"github.com/spf13/viper"
)

type Configuration struct {
	SearchServiceEndPoint string
}

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	var configuration Configuration

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error al leer del archivo config, %s", err)
	}
	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("Error al deserializar la configuración, %s", err)
	}

	log.Printf("La ruta del servicio de búsqueda es %s", configuration.SearchServiceEndPoint)
}
