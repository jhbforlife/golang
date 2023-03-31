package main

import "github.com/robfig/cron/v3"

func startCron() {
	c := cron.New()

	c.AddFunc("0 0 * * *", updateLanguagesTable)
	c.AddFunc("0 1 * * 0", resetTranslationsTable)

	c.Start()
}
