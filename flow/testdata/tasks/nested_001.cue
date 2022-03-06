@flow()

nested: {
  @task(nest)

  tasks: {
    get: { FOO: _ } @task(os.Getenv)
    out: { text: get.FOO + "\n" } @task(os.Stdout)
    foo: get.FOO
  }
  out: { text: tasks.get.FOO + "\n" } @task(os.Stdout)

  foo: tasks.foo
}

out: { text: nested.foo + "\n" } @task(os.Stdout)
