package main

type Read interface {
	read(p []byte) (int, error)
}

type Write interface {
	write(p []byte) (int, error)
}

// composed interface
type ReadAndWrite interface {
	Read
	Write
}
