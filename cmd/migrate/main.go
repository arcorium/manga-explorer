package main

import (
  "log"
  "manga-explorer/cmd/migrate/command"
)

func main() {
  err := command.Execute()
  if err != nil {
    log.Fatalln(err)
  }
}
