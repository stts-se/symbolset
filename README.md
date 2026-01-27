# symbolset

Symbolset is a repository for handling phonetic symbol sets and mappers/converters between different symbol sets and languages. Written in `go`.

[![GoDoc](https://godoc.org/github.com/stts-se/symbolset?status.svg)](https://godoc.org/github.com/stts-se/symbolset)
[![Go Report Card](https://goreportcard.com/badge/github.com/stts-se/symbolset)](https://goreportcard.com/report/github.com/stts-se/symbolset) [![Build Status](https://travis-ci.com/stts-se/symbolset.svg?branch=master)](https://app.travis-ci.com/stts-se/symbolset)

## I. Server installation

1. Set up `go`

     Download: https://golang.org/dl/ (1.24 or higher)   
     Installation instructions: https://golang.org/doc/install             


2. Clone the source code

   ``` sh
   $ git clone https://github.com/stts-se/symbolset.git
   $ git clone https://github.com/stts-se/wikispeech-lexdata.git # optional
   $ cd symbolset/
   ```   
   
4. Test (optional)

   ```sh
   symbolset$ go test ./...
   ```


6. Pre-compile server (for faster execution times).

   ``` sh
   symbolset$ cd server/
   server$ go build .
    ```


## II. Setup with Wikispeech symbolsets (optional)

1. Setup 

    ``` sh
   server$ bash setup.sh ../../wikispeech-lexdata/ ss_files
   ```

3. Start server

    ``` sh
   server$ ./server -ss_files ss_files/
   ```

## IIb. Alternative: Start the server without Wikispeech symbol sets (using demo symbol sets)

    ```
    server$ ./server -ss_files demo_files/
    ```


---

_This work was supported by the Swedish Post and Telecom Authority (PTS) through the grant "Wikispeech – en användargenererad talsyntes på Wikipedia" (2016–2017), and the Swedish Inheritance Fund (Allmänna arvsfonden) through the grant "Wikispeech – Talsyntes och taldatainsamlare." (2024–2027)._
