# Chage Exporter
Chage exporter is a simple tool that exports how long a password has before it requires to be changed. By default it exports this data for all users, but can be configured to only send for a select few.

## Consider this...
It is important to note that the chage command requires elevated privileges to run, for good reason. If possible, I recommend hiding this exporter behind a reverse proxy with SSL / Password Protection.