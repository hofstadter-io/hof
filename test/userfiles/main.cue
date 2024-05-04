package main

hello: "world"

files: {
  @userfiles(content,trim=content)
}