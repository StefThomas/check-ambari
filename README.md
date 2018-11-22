# check-ambari
Monitor HDP/HDF cluster on Nagios from Ambari API.

The plugin’s purpose is to monitor HDP/HDF plateform throught Ambari API. It’s intended to be used with external monitoring tools like Nagios, Centreon, Shinken, Sensu, etc…

The following program must be called by your monitoring tool. It returns the status (Nagios status normalization) of node/service with a human-readable messages and sometimes perfdata.

This program makes calls to the Ambari API to get the state of your HDP/HDF plateform. The availability of this API is then required to monitor the whole set of services of your plateform. You should probably set a parenthood relation between the monitoring of the API itself and all the other Ambari services.

You can use it to monitor the following section of your cluster:
- *Node state*: It checks that there is no alert for your host.
- *Service state*: It checks that there are no alert for your service.

## Usage

### Global parameters

You need to set the Ambari API informations for all checks.

```sh
./check-ambari --ambari-url https://ambari.company.com --ambari-login admin --ambari-password admin …
```

You need to specify the following parameters:
- **ambari-url**: It's your Ambari URL. Alternatively you can use the environment variable `AMBARI_URL`.
- **ambari-login**: It's the Ambari login to use when it call the API. Alternatively you can use the environment variable `AMBARI_LOGIN`.
- **ambari-password**: It's the password associated with the login. Alternatively you can use the environment variable `AMBARI_PASSWORD`.

You can also set this parameters in a YAML file(s) and use the `--config` parameter with the path of your main (or unique) YAML file.
```yaml
---
ambari-url: https://ambari.company.com
ambari-login: admin
ambari-password: admin
```

### Check the state of a node

You need to run the following command:

```sh
./check-ambari --ambari-url https://ambari.company.com --ambari-login admin --ambari-password admin check-node --cluster-name test --node-name worker01.company.com
```

You need to specify the following parameters:
- **--cluster-name**: The cluster’s name where to check node state.
- **--node-name**: The node’s name to check (usally its FQDN name)
- **--include-alerts**: Check only the alerts in this list, separated by coma.
- **--exclude-alerts**: Don’t check the alerts in this list, separated by coma.

This check follows this logic:
1. `OK` when there is no Ambari alert for node
2. `WARNING` when there is one or more warning Ambari alert for node
3. `CRITICAL` when there is one or more critical Ambari alert for node

> All Ambari nodes which have a problem are displayed on the output.
> Alerts with an `UNKNOWN` state are ignored (like in Ambari UI).

It returns the following perfdata:
- **nbAlert**: the number of current Ambari alerts


### Check the state service

You need to run the following command:

```sh
./check-ambari --ambari-url https://ambari.company.com --ambari-login admin --ambari-password admin check-node --cluster-name test --service-name hdfs
```

You need to specify the following parameters:
- **--cluster-name**: The cluster’s name where you want to check service state.
- **--service-name**: The service’s name to check.
- **--exclude-node-alerts** (optionnal): look only alert in the service scope, so it excludes all node alerts.
- **--include-alerts**: Check only the alerts in this list, separated by coma.
- **--exclude-alerts**: Don’t check the alerts in this list, separated by coma.

This check follows this logic:
1. `OK` when there is no Ambari alert in service
2. `WARNING` when there is one or more warning Ambari alert for service
3. `CRITICAL` when there is one or more critical Ambari alert for service

> All Ambari service which have a problem are displayed on the output.
> Alerts with an `UNKNOWN` state are ignored (like in Ambari UI).

It returns the following perfdata:
- **nbAlert**: the number of current Ambari alerts
