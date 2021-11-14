package main

import (
	"github.com/nialek/go-config-watcher"
	"github.com/spf13/viper"
	"log"
	"time"
)

func main() {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}

	watcher, _ := configwatcher.New(v)
	number := watcher.Get("number", func(old interface{}, new interface{}) {
		log.Println("Number:", old, "->", new)
	})
	word := watcher.Get("word", func(old interface{}, new interface{}) {
		log.Println("Word:", old, "->", new)
	})

	log.Println("Initial:", "number =", number, "word =", word)

	time.Sleep(time.Minute * 5)
}
