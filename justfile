
make:
    mkdir -p bin
    go build -o ./bin/cloaktray .

install:
    cp ./bin/cloaktray /usr/local/bin