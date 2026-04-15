package main

import "sync"

func or(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	out := make(chan interface{})
	var once sync.Once
	
	for _, ch := range channels{
		go func(ch <-chan interface{}) {
			select
		} () 
	}
	
	return out
}
