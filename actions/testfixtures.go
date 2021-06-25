package actions

// trainer
// Copyright 2021 The Home Depot
// Licensed under the Apache License v2.0
// See LICENSE for further details.

import "github.com/MakeNowJust/heredoc"

var success = heredoc.Doc(`{
  "result": {
    "allyourbase.example.com": {
      "ci": "allyourbase.example.com",
      "finished": 1,
      "finishedbool": true,
      "success": 1,
      "successbool": true
    }
  }
}`)

var failure = heredoc.Doc(`{
  "result": {
    "allyourbase.example.com": {
      "ci": "allyourbase.example.com",
      "finished": 1,
      "finishedbool": true,
      "success": 0,
      "successbool": false
    }
  }
}`)

var successtrunc = heredoc.Doc(`{
  "result": {
    "allyourbase.example.com": {
      "ci": "allyourbase.example.com",
      "finished": 1,
      "finishedbool": true,
      "success": 1,
      "successbool": true
    }
  }
}`)
