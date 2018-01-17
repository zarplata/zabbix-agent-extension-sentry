# zabbix-agent-extension-sentrt

zabbix-agent-extension-sentry - this extension for monitoring [sentry](https://github.com/getsentry/sentry).

### Supported features

- Discovering projects.
- Discovering queues.
- Fetch event per minutes statistic for each projects and organization too.
- Fetch event in queue for each queue.

### Installation

#### Manual build

```sh
# Building
git clone https://github.com/zarplata/zabbix-agent-extension-sentry.git
cd zabbix-agent-extension-sentry
make

#Installing
make install

# By default, binary installs into /usr/bin/ and zabbix config in /etc/zabbix/zabbix_agentd.conf.d/ but,
# you may manually copy binary to your executable path and zabbix config to specific include directory
```

### Dependencies

zabbix-agent-extension-sentry requires [zabbix-agent](http://www.zabbix.com/download) v3.4+ to run.

`WARNING:` You must define macro with name - `{$ZABBIX_SERVER_IP}` in global or local (template) scope with IP address of  zabbix server.
Also define next macro in host with sentry:
- {$SENTRY_ORG} = sentry organization name (sentry).
- {$SENTRY_PORT} = wsgi port (9000)
- {$SENTRY_TOKEN} = api token with `event:read, org:read, project:read`
- {$SENTRY_URL} = url sentry (http://sentry.io)

Set up ENV `SENTRY_CONF=/etc/sentry` in file `/etc/environment`.
