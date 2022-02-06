package debug

test1: {
  @flow()
  t: { FOO: _ } @task(os.Getenv)
  o: { text: "(ct): the cow goes \(t.FOO)\n" } @task(os.Stdout)
}
