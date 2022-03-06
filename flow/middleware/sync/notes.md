### sync.Pool

- how can this be a task?


```cue
foo: {
  @pool()

  load: {
    @task()
  }

  enrich: {
    @timeout()
    @pool() // middleware
    @task(api.Call)
  }

  send: {
    @task(msg.Send)
  }
}
```

