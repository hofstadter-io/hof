package main

hello: "world"

files: {
  @userfiles(content,trim=content)
  @userfiles(other,trim=other)
}