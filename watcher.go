package configwatcher

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type (
	callbackFunc func(old interface{}, new interface{})
	callbackMap  map[string]callbackFunc
	valueMap     map[string]interface{}
)

type Watcher struct {
	v         *viper.Viper
	callbacks callbackMap
	values    valueMap
}

func New(v *viper.Viper) (*Watcher, error) {
	w := &Watcher{
		v:         v,
		values:    make(valueMap),
		callbacks: make(callbackMap),
	}

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		w.traverseChanges()
	})

	return w, nil
}

func (w *Watcher) traverseChanges() {
	for key, value := range w.values {
		go func(key string, value interface{}) {
			valueFromConfig := w.v.Get(key)
			if valueFromConfig != value {
				w.values[key] = valueFromConfig
				w.callbacks[key](value, valueFromConfig)
			}
		}(key, value)
	}
}

func (w *Watcher) Get(key string, callback callbackFunc) interface{} {
	value := w.v.Get(key)
	w.values[key] = value
	w.callbacks[key] = callback
	return value
}
