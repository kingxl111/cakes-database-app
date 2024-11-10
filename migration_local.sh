#!/bin/bash
source local.env

sleep 5 && migrate -path "./migrations" -database "${MIGR_DSN}" up