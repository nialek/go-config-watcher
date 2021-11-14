package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"time"
)

type (
	valueMap   map[string]interface{}
	channelMap map[string]chan interface{}
)

type Watcher struct {
	v        *viper.Viper
	channels map[string]chan interface{}
	values   valueMap
}

func New(v *viper.Viper) (*Watcher, error) {
	w := &Watcher{
		v:        v,
		values:   make(valueMap),
		channels: make(channelMap),
	}

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		w.traverseChanges()
	})

	return w, nil
}

func (w *Watcher) traverseChanges() {
	for key, value := range w.values {
		valueFromConfig := w.v.Get(key)
		if valueFromConfig != value {
			w.values[key] = valueFromConfig
			w.channels[key] <- valueFromConfig
		}
	}
}

func (w *Watcher) GetChannel(key string) <-chan interface{} {
	value := w.v.Get(key)
	w.values[key] = value

	ch := make(chan interface{}, 1)
	w.channels[key] = ch
	ch <- value

	return ch
}

func main() {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}

	watcher, _ := New(v)
	number := watcher.GetChannel("number")
	word := watcher.GetChannel("word")

	go func() {
		for {
			select {
			case n := <-number:
				log.Println("Number:", n)
			case w := <-word:
				log.Println("Word:", w)
			}
		}
	}()

	time.Sleep(time.Minute * 5)
}
