package main

import (
	"log"
	"os"
	"time"

	"github.com/faiface/beep"

	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type MusicControl struct {
	f beep.Format
	s beep.Streamer

	v *effects.Volume
}

func NewMusicControl() MusicControl {
	return MusicControl{}
}

func (m *MusicControl) Init() {
	f, err := os.Open("afterglow.mp3")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer), Paused: false}
	volume := &effects.Volume{
		Streamer: ctrl,
		Base:     2,
		Volume:   -4.0,
		Silent:   false,
	}

	m.f = format
	m.s = streamer
	m.v = volume

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
}

func (m *MusicControl) Play() {
	speaker.Play(m.v)
}

func (m *MusicControl) Close() {
	speaker.Close()
}
