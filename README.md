# Chage Exporter
Chage exporter is a simple tool that exports how long a password has before it requires to be changed. By default it exports this data for all users, but can be configured to only send for a select few.

## Usage
chage_exporter is a simple binary that offers a couple of inputs. Here is the `--help` menu:
```
usage: chage_exporter --config=CONFIG [<flags>]

Flags:
      --help           Show context-sensitive help (also try --help-long and
                       --help-man).
  -p, --port=9200      Port for chage_exporter to listen on.
  -c, --config=CONFIG  Path to chage_exporter config.
```
chage_exporter requires you to specify a configuration path.

You can and should run chage_exporter with a service file (if you're rolling with systemd). Here's an example of a simple service file:
```
[Unit]
Description=a Prometheus exporter to monitor the age of UNIX user passwords
After=network-online.target

[Service]
Type=simple
ExecStart=/path/to/binary --config=/path/to/config


[After]
WantedBy=multi-user.target
```

## Configuration
The configuration of chage_exporter is very simple, here's an example config:
```yaml
users:
  - "user1"
  - "user2"
```

## Consider this...
It is important to note that the chage command requires elevated privileges to run, for good reason. If possible, I recommend hiding this exporter behind a reverse proxy with SSL / Password Protection if your server is open to the public.