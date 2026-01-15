package examples

import "strings"

// inputs supplied via tags
inputs: {
  model: string @tag(model)
  prompt: string @tag(prompt)
  msg: string @tag(msg)
}

vertex_chat: {
  @flow() // define a flow

  steps: {

    // task: get auth from external command
    gcp: {
      @task(os.Exec)
      cmd: ["gcloud", "auth", "print-access-token"]
      stdout: string
      key:    strings.TrimSpace(stdout)
    }

    // task: api call via reusable task
    call: _gemini & {
      apikey: gcp.key

      model: inputs.model
      prompt: inputs.prompt
      msg: inputs.msg

      resp: body: _
    }

    // task: print text to std output
    out: {
      @task(os.Stdout)
      text: call.final.text
    }

  }
}

// reusable task
_gemini: {
  @task(api.Call)

  apikey: string
  model: string | *"gemini-1.5-pro-001:generateContent"
  prompt: string | *"You are an assistant who is very concise when responding."
  msg: string

  req: {
    host: "https://us-central1-aiplatform.googleapis.com"
    path: "/v1/projects/hof-io--develop/locations/us-central1/publishers/google/models/\(model)"
    headers: {
      "Content-Type": "application/json"
      Authorization:  "Bearer \(apikey)"
    }
    data: {
      systemInstruction: {
        role: "MODEL"
        parts: [{
          text: prompt
        }]
      }

      contents: [{
        role: "USER"
        parts: [{
          text: msg
        }]
      }]
    }
    method: "POST"
  }

  resp: {
    body: {}
  }
  // @print(resp.body)

  // task-local ETL
  final: {
    cand: resp.body.candidates[0]
    text: cand.content.parts[0].text
  }
}