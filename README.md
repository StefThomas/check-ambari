# check-ambari
Monitor HDP/HDF cluster on Nagios from Ambari API.

It's a general purpose to monitore HDP/HDF plateform throught Ambari with external monitoring tools like Nagios, Centreon, Shinken, Sensu, etc.

The following program must be called by your monitoring tools and it return the status (nagios status normalization) with human messages and some times perfdatas.
This program called Ambari API to compute the state of your HDP/HDF plateform.

You can use it to monitore the following section of your cluster:
- *Node state*: It's check that there are no alert in your host.
- *Service state*: It's check that there are no alert in your service.

## Usage

### Global parameters

You need to set the Ambari API informations for all checks.

```sh
./check-ambari --ambari-url https://ambari.company.com --ambari-login admin --ambari-password admin ... 
```

You need to specify the following parameters:
- **ambari-url**: It's your Ambari URL. Alternatively you can use environment variable `AMBARI_URL`.
- **ambari-login**: It's the Ambari login to use when it call the API. Alternatively you can use environment variable `AMBARI_LOGIN`.
- **ambari-password**: It's the password associated with the login. Alternatively you can use environment variable `AMBARI_PASSWORD`.

You can set also this parameters on yaml file (one or all) and use the parameters `--config` with the path of your Yaml file.
```yaml
---
ambari-url: https://ambari.company.com
ambari-login: admin
ambari-password: admin
```

### Check the state node

You need to lauch the following command:

```sh
./check-ambari --ambari-url https://ambari.company.com --ambari-login admin --ambari-password admin check-node --cluster-name test --node-name worker01.company.com
```

You need to specify the following parameters:
- **--cluster-name**: The cluster name where you should to check the node state.
- **--node-name**: The node name that you should to check (normally it's fqdn name)
- **--include-alerts**: Check only the alerts in this list, separated by coma.
- **--exclude-alerts**: Don't check the alerts in this list, separated by coma.

This check follow this logic:
1. `OK` when there are no Ambari alert in node
2. `WARNING` when there are one or more warning Ambari alert in node
3. `CRITICAL` when there are one or more critical Ambari alert in node

> All Ambari that have problem is displayed on the outpout.
> The alert with `UNKNOWN` state is ignored (like in Ambari UI)

It's return the following perfdata:
- **nbAlert**: the number of current Ambari alert


### Check the state service

You need to lauch the following command:

```sh
./check-ambari --ambari-url https://ambari.company.com --ambari-login admin --ambari-password admin check-node --cluster-name test --service-name hdfs
```

You need to specify the following parameters:
- **--cluster-name**: The cluster name where you should to check the node state.
- **--service-name**: The service name that you should to check.
- **--exclude-node-alerts** (optionnal): look only alert in service scope, so it's exclude all node alerts.
- **--include-alerts**: Check only the alerts in this list, separated by coma.
- **--exclude-alerts**: Don't check the alerts in this list, separated by coma.

This check follow this logic:
1. `OK` when there are no Ambari alert in service
2. `WARNING` when there are one or more warning Ambari alert in service
3. `CRITICAL` when there are one or more critical Ambari alert in service

> All Ambari that have problem is displayed on the outpout.
> The alert with `UNKNOWN` state is ignored (like in Ambari UI)

It's return the following perfdata:
- **nbAlert**: the number of current Ambari alert
