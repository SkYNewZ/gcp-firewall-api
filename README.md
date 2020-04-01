# GCP Firewall API

This repository provides an API to create and manage Firewall rules in a GCP host project using API for an application.

## Getting started

THis project is deployed on 2 environements:

- Production: https://api.cloudservices.tech.adeo.cloud
- Stagging: https://gcp-firewall-api-2q3jhrmuuq-ew.a.run.app

## Create a rule

Rules are based on Google compute API [rest/v1/firewalls](https://cloud.google.com/compute/docs/reference/rest/v1/firewalls)
So, create a Google Rule, for example:

```json
{
  "network": "global/networks/lh-network",
  "allowed": [
    {
      "IPProtocol": "tcp",
      "ports": ["443"]
    }
  ]
}
```

And `POST` it to `/project/<LH>/service_project/<LZV2>/application/<APP>/firewall_rule/<NAME>`.

- `<LH>` Landing Hub project ID which host your Landing Zone v2
- `<LZV2>` your Landing Zone v2 project ID
- `<APP>` an arbitrary application name
- `<NAME>` your wanted firewall rule name

The final firewall rule name will be `<LZV2>-<APP>-<NAME>`. **It will be the same for the target tag.**

It will return the given [schema](#schema)

## List your application rules

`GET /project/<LH>/service_project/<LZV2>/application/<APP>/`

It will return the given [schema](#schema)

## Get a specific rule

`GET /project/<LH>/service_project/<LZV2>/application/<APP>/firewall_rule/<NAME>`

It will return the given [schema](#schema)

## Delete a specific rule

`DELETE /project/<LH>/service_project/<LZV2>/application/<APP>/firewall_rule/<NAME>`

It will return the given [schema](#schema)

## Schema

```json
{
  "application": "<APP>",
  "data": [
    {
      "custom_name": "<NAME>",
      "item": "*GoogleRule"
    }
  ],
  "project": "<LH>",
  "service_project": "<LZV2>"
}
```

## Deployement

Deployements are made by GitLabCI with service accounts.

Deployer service account have roles:

- `roles/run.admin` to deploy a new Cloud Run revision
- `roles/storage.admin` to store built Docker image
- `roles/iam.serviceAccountUser` https://cloud.google.com/run/docs/reference/iam/roles#additional-configuration

Runtime service account, **on each environement**, have roles:

- `roles/viewer` to view Compute resources
- `roles/compute.securityAdmin` to create network resources (of course to create firewall rules)

All theses credentials are stored in Vault on path `secret/gcp-firewall-api/*`
