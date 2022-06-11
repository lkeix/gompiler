link:
	as -o main.o main.s runtime.s && \
	ld -o main.out main.o

run:
	go run main.go > main.s && \
	as -o main.o main.s && \
	ld -o a.out main.o && \
	mv a.out main.out && \
	./main.out

test: main.out
			./test.sh

assemble:
	go run main.go > main.s && \
	as -o main.o main.s

build:
	go run main.go > main.s && \
	as -o main.o main.s && \
	ld -o a.out main.o && \
	mv a.out main.out

clean:
	rm -rf *.s *.o *.out