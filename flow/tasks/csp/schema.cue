// CSP (Communicating Sequential Processes)
// Named channels (mailboxes) and send/recv messages
// Similar to Go / Erlang concurrency
// Can be accessed from different flows/tasks
package csp

// Chan is a named mailbox
Chan: {
  @task(csp.Chan)
  $task: "csp.Chan"

  // the name of the channel
  mailbox: string

  // how many messages can queue
  // defaults to a blocking send / recv
  buf: int | *0
}

// Send a message to a mailbox
Send: {
  @task(csp.Send)
  $task(csp.Send)

  // the name of the channel
  mailbox: string

  // an optional key
  key?: string

  // the message value
  val: _
}

// Recv is a coroutine which runs indefinitely
Recv: {
  @task(csp.Recv)
  $task(csp.Recv)

  // the name of the channel
  mailbox: string

  // the name of a second mailbox to quit on any message received
  quitMailbox?: string

  // the handler for a message, run as a flow per message
  handler: {
    // filled when a message is received
    msg: _

    // do whatever you like in here
  }
}
