package cron

import (
	"errors"

	"github.com/mileusna/crontab"
)

var CronObject *crontab.Crontab

func NewCron() (error, bool) {
	ctab := crontab.New()
	if ctab != nil {
		CronObject = ctab
	} else {
		return errors.New("Cron instance null"), false
	}

	return nil, true

}

func GetCronInstance() *crontab.Crontab {
	return CronObject
}
