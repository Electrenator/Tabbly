#!/bin/bash

# Add to cron using 'cd [project path]; bash autostart.sh'
python -m venv venv &&
source venv/bin/activate &&
(
	until pip install -r requirements.txt;
	do # until successful execution / dependency download
		echo -e "\nRetrying dependency install in 1 minute...\n";
		sleep 60;
	done;

	python src/main.py --verbose
) &>> log/crontab-$(hostname).log;

deactivate
