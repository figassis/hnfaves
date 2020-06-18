#!/bin/bash
sed 's/.*"\(.*\)".*/\1/' <<< "`grep "	version =" pkg/utl/config/config.go`"
