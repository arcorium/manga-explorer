package mail

import (
  "strings"
  "text/template"
)

type Mail struct {
  Recipients  []string
  Subject     string
  BodyType    BodyType
  Body        string
  EmbedFiles  []string
  AttachFiles []string
}

func NewHTML(filename string, link string, path ...string) (*Mail, error) {
  paths := "template/"
  if len(path) == 1 {
    paths = path[0]
  }
  tmpl, err := template.New(filename).ParseFiles(paths + filename)
  if err != nil {
    return nil, err
  }

  var builder strings.Builder
  err = tmpl.Execute(&builder, struct {
    Link string
  }{Link: link})

  if err != nil {
    return nil, err
  }

  return &Mail{
    BodyType: BodyTypeHTML,
    Body:     builder.String(),
  }, nil
}

type BodyType string

func (b BodyType) String() string {
  return string(b)
}

const (
  BodyTypeHTML  BodyType = "text/html"
  BodyTypePlain          = "text/plain"
)
