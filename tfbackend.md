# TF TOOLS

## TFBACKEND SWITCH
The program searches for all `*.tfbackend` files in a directory and let you choose which backend you want to configure

---

## To install/upgrade:

1. Download, optionally inspect, and copy the script to an appropriate folder
2. Add aliases:

**bash**:

```
echo -e "\nalias tfbackend='source /usr/local/bin/tfbackend.sh'" >> ~/.bash_profile
```
    
**zsh**:

```
echo -e "\nalias tfbackend='source /usr/local/bin/tfbackend.sh'" >> ~/.zshrc
```

***oh-my-zsh***:

```
echo -e "\nalias aws-switch='source /usr/local/bin/tfbackend.sh'" >> ~/.oh-my-zsh/custom/aliases.zsh
```

### Usage:

Using a new terminal window or tab (required for the new alias settings to take effect), simply run the script using the alias created above (tfbackend):

``` tfbackend ```

<details>
<summary>Example output `tfbackend` : </summary>

```
tfbackend                                                                                                                                                                                       ------------- Select Backend -------------
Type the number of the backend you want to use from the list below, and press enter

-: Unset backend
0: nonprod.tfbackend
1: prod.tfbackend

Selection: 0
Activating backend 0: nonprod.tfbackend

Initializing the backend...

Successfully configured the backend "s3"! Terraform will automatically
use this backend unless the backend configuration changes.
Initializing modules...

Initializing provider plugins...
- terraform.io/builtin/terraform is built in to Terraform
- Reusing previous version of hashicorp/aws from the dependency lock file
- Using previously-installed hashicorp/aws v5.33.0

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
``` 
</details>

# TFPLAN / TFAPPLY
These scripts look for the `.tfvars` files for the configured backend. The name of the `.tfvars` file need to be the same as the `.tfbackend` (eg: `nonprod.tfbackend` and `nonprod.tfvars`)

Additional parameters can be parsed. For example: `-target=` or `-auto-approve` or `-out <file.tfplan>`

If no `backend` is configured, no `vars-file` is being used, it's behaviour is like a normal `terraform plan` or `terraform apply` command.

# TFDESTROY
These scripts look for the `.tfvars` files for the configured backend. The name of the `.tfvars` file need to be the same as the `.tfbackend` (eg: `nonprod.tfbackend` and `nonprod.tfvars`).

It will __**prevent**__ the destroy of items in the `terraform-state file` with the names:

- backend
- dynamodb
- kms

This is very useful for projects that are using the wearetechnative-backend module.
