Game of Life
------------
This project implements Conway's Game of Life in Golang.
The performance can be improved by specifying the number of goroutines
to use. The board is divided into a grid with the number of cells equal to the number
of goroutines. As the number of goroutines P has to be a perfect square (1, 4, 9 ...), sqrt(P) is supplied as input argument.

Example
------------
The following command loads in a predefined world in setup/pulsar.txt and divides to work over 4(=2^2) goroutines. It computes
the state of the world over 150 iterations and introduces a delay between each update of 100ms.

`go run cmd/gol/main.go -filename=setup/pulsar.txt -sqrtP=2 -delay=100
`

The following command creates a random world of 50 by 50 with a cell density of 0.3 and divides to work over 4(=2^2) goroutines. 
It computes the state of the world over 300 iterations and introduces a delay between each update of 30ms.

`go run cmd/gol/main.go -rows=50 -cols=50 -sqrtP=2 -nits=300
`

The following command creates a board of 4096 by 4096 with a cell density of 0.5 and divides to work over 256(=16^2) goroutines. 
It computes the state of the world over 1000 iterations and does not show the board on the terminal.

`go run cmd/gol/main.go -rows=4096 -cols=4096 -sqrtP=16 -nits=1000 -dLevel=0.5 -plot=false
`


