package main

import "github.com/MakeNowJust/heredoc"

var (
	finishedsuccessfulresultoneci = heredoc.Doc(`
	{
		"result": {
			"allyourbase.example.com": {
				"changeciid":"b9152d04-90de-4512-9dc5-8416b9bcae5d",
				"ci":"allyourbase.example.com",
				"status":"not running",
				"module":"",
				"starttime":"2018-11-16 21:00:57.968 +0000 UTC",
				"finished":1,
				"finishedbool":true,
				"success":1,
				"successbool":true,
				"error":"",
				"guid":"ee94e156-f103-45aa-a50a-0605618f2bbb",
				"report":"13ec10608a9519ec1f62b7ef75037739b7d6a865",
				"retry_reason":"",
				"retry_count":0,
				"max_retries_reached":"0",
				"is_max_retries_reached":false,
				"orchestrator_job":"12345",
				"aborted":false
			}
		}
	}
`)

	unsuccessfulresultoneci = heredoc.Doc(`
	{
		"result": {
			"allyourbase.example.com": {
				"changeciid":"b9152d04-90de-4512-9dc5-8416b9bcae5d",
				"ci":"allyourbase.example.com",
				"status":"not running",
                "module":"",
                "starttime":"2018-11-16 21:00:57.968 +0000 UTC",
                "finished":0,
                "finishedbool":false,
                "success":0,
				"successbool":false,
				"error":"",
				"guid":"ee94e156-f103-45aa-a50a-0605618f2bbb",
				"report":"13ec10608a9519ec1f62b7ef75037739b7d6a865",
				"retry_reason":"",
				"retry_count":0,
				"max_retries_reached":"0",
				"is_max_retries_reached":false,
				"orchestrator_job":"12345",
				"aborted":false
			}
		}
	}
`)

	finishedunsuccessfulresultoneci = heredoc.Doc(`
	{
		"result": {
			"allyourbase.example.com": {
				"changeciid":"b9152d04-90de-4512-9dc5-8416b9bcae5d",
				"ci":"allyourbase.example.com",
				"status":"not running",
                "module":"",
                "starttime":"2018-11-16 21:00:57.968 +0000 UTC",
                "finished":1,
                "finishedbool":true,
                "success":0,
				"successbool":false,
				"error":"",
				"guid":"ee94e156-f103-45aa-a50a-0605618f2bbb",
				"report":"13ec10608a9519ec1f62b7ef75037739b7d6a865",
				"retry_reason":"",
				"retry_count":0,
				"max_retries_reached":"0",
				"is_max_retries_reached":false,
				"orchestrator_job":"12345",
				"aborted":false
			}
		}
	}
`)
)
