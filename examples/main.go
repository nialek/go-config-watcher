package main

import (
	"github.com/nialek/go-config-watcher"
	"github.com/spf13/viper"
	"log"
	"time"
)

func main() {
	// create a viper instance
	v := viper.New()
	v.AddConfigPath("./examples")
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}

	// pass the viper instance in constructor, also decide if you want it to automatically reload and execute callbacks
	watcher, _ := configwatcher.New(v, true)
	number := watcher.Get("number", func(old interface{}, new interface{}) {
		log.Println("Number:", old, "->", new)
	})
	word := watcher.Get("word", func(old interface{}, new interface{}) {
		log.Println("Word:", old, "->", new)
	})

	log.Println("Initial:", "number =", number, "word =", word)

	// try modifying the config.yaml

	time.Sleep(time.Minute * 5)
}
