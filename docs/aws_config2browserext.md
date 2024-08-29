# AWS config converter

## A painless way to configure the aws config file for browser extention aws-extend-switch-roles

This script converts the `.aws/config` file configured for `aws-cli` to the config file that is used by the browser extention [aws-extend-switch-roles](https://github.com/tilfinltd/aws-extend-switch-roles)

It uses the config file in the home-directory of the current user:

`${HOME}/.aws/config`

---
### Usage

```bash
 ./aws_config2browserext.sh
```

<details>
<summary>Example output</summary>

```bash
/aws_config2browserext.sh                                                                                                                                                                                               tn-finops-
[profile customer1]
role_arn = arn:aws:iam::123456789012:role/landing_zone_devops_administrator
region = eu-central-1

[profile customer2]
role_arn = arn:aws:iam::223456789012:role/landing_zone_devops_administrator
region = eu-central-1

[profile customer1]
role_arn = arn:aws:iam::323456789012:role/landing_zone_devops_administrator
region = us-west-2
```
</details>

###### Wouter van der Toorren.