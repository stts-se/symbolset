# symbolset

Symbolset is a repository for handling phonetic symbol sets and mappers/converters between different symbol sets and languages. Written in `go`.

[![GoDoc](https://godoc.org/github.com/stts-se/symbolset?status.svg)](https://godoc.org/github.com/stts-se/symbolset)
[![Go Report Card](https://goreportcard.com/badge/github.com/stts-se/symbolset)](https://goreportcard.com/report/github.com/stts-se/symbolset) [![Build Status](https://travis-ci.org/stts-se/symbolset.svg?branch=master)](https://travis-ci.org/stts-se/symbolset)

## I. Server installation

1. Set up `go`

     Download: https://golang.org/dl/ (1.13 or higher)   
     Installation instructions: https://golang.org/doc/install             


2. Clone the source code

   `$ git clone https://github.com/stts-se/symbolset.git`  
   `$ cd symbolset`   
   
3. Test (optional)

   `symbolset$ go test ./...`


4. Pre-compile server (for faster execution times).

    `symbolset$ cd server`   
    `server$ go build .`


## II. Quick start: Start the server with demo set of symbol sets

`server$ ./server -ss_files demo_files`


## III. Setup with Wikispeech symbolsets

1. Clone Wikispeech lexdata (this might take a couple of minutes)

   `$ git clone https://github.com/stts-se/wikispeech-lexdata.git`


2. Setup 

    `server$ bash setup.sh wikispeech-lexdata ss_files`


3. Start server

    `server$ ./server -ss_files ss_files`
