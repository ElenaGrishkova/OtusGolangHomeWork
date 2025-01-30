#!/usr/bin/env bash

echo -e "BAR is (${BAR})
arg1 is (${arg1})
arg2 is (${arg2})
UNSET is (${UNSET})
arguments are $*"
exit ${arg1}