FROM golang:latest

WORKDIR /app
COPY main.go main.go
# -s - Omit the symbol table and debug information.
# -w - Omit the DWARF symbol table.
RUN go build -ldflags="-s -w" -o main main.go
CMD /app/main
