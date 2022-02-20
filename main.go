package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

func AudioPipe(audio []float64) beep.Streamer {
	idx := 0
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		var i int
		for i = range samples {
			if i+idx >= len(audio) {
				return i, false
			}
			samples[i][0], samples[i][1] = audio[i+idx], audio[i+idx]
		}

		idx += i
		return len(samples), true
	})
}

type AudioBuffer struct {
	Buffer     []float64 `json:"buf"`
	SampleRate int       `json:"sampleRate"`
	BufferSize int       `json:"bufferSize"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var aud AudioBuffer
		err := json.NewDecoder(r.Body).Decode(&aud)
		if err != nil {
			log.Panic(err)
		}

		go Play(aud)
	})

	http.ListenAndServe(":1926", nil)
}

func Play(audio AudioBuffer) {
	err := speaker.Init(beep.SampleRate(audio.SampleRate), audio.BufferSize)
	if err != nil {
		log.Panic(err)
	}

	done := make(chan bool)
	speaker.Play(beep.Seq(AudioPipe(audio.Buffer), beep.Callback(func() {
		fmt.Println("Done!")
		done <- true
	})))

	<-done
}

func Speaker(audio []float64) {
	speaker.Init(48000, 8192)

	done := make(chan bool)
	speaker.Play(beep.Seq(AudioPipe(audio), beep.Callback(func() {
		fmt.Println("Done!")
		done <- true
	})))

	<-done
}

