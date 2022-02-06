package do

import (
  "encoding/json"
  "hof.io/repro/load"
)

vars: {
  user: string | *"dr_verm" @tag(user)
}

meta: {
  secrets: {
    tLoad: load.load
    token: tLoad.say
  }

  req: {
    url: "https://postman-echo.com/get?cow=\(secrets.token)"
    method: "GET"
  }
}

command: {
  do: {
    cfg: meta
    get: {
      req: cfg.req & { 
        $id: "tool/http.Do"
      }
      resp: req.response
      Out: json.Indent(resp.body, "", " ")
    }
    print: {
      $id: "tool/cli.Print"
      text: "\(get.Out) @\(vars.user)"
    }
  }
}
