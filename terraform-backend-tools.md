# TF TOOLS

- ## TFBACKEND
This script searches for all `*.tfbackend` files in a directory and let you choose which backend you want to configure

- ## TFPLAN / TFAPPLY
These scripts look for the `.tfvars` files for the configured backend.
The name of the `.tfvars` file need to be the same as the `.tfbackend` (eg: `nonprod.tfbackend` and `nonprod.tfvars`)

Additional parameters can be parsed. For example: `-target=` or `-auto-approve` or `-out <file.tfplan>`

If no `backend` is configured, no `vars-file` is being used, it's behaviour is like a normal `terraform plan` or `terraform apply` command.

- ## TFDESTROY
This script look for the `.tfvars` files for the configured backend.
The name of the `.tfvars` file need to be the same as the `.tfbackend` (eg: `nonprod.tfbackend` and `nonprod.tfvars`).

It will __**prevent**__ the destroy of items in the `terraform-state file` with the names:

- backend
- dynamodb
- kms

This is very useful for projects that are using the wearetechnative-backend module.

---

### To install/upgrade:

1. Download or copy the script(s) to "/usr/local/bin", add it to your PATH, and make it executable with `chmod +x /usr/local/bin/<script_name>.sh` replace <script_name> with the name of the script you want to use.
2. Add aliases:

**bash**:

For Ubuntu users (if alias not work use ~/.profile instead of ~/.bashrc):
```
echo -e "\nalias tfapply='source /usr/local/bin/tfapply.sh'" >> ~/.bashrc
echo -e "\nalias tfbackend='source /usr/local/bin/tfbackend.sh'" >> ~/.bashrc
echo -e "\nalias tfdestroy='source /usr/local/bin/tfdestroy.sh'" >> ~/.bashrc
echo -e "\nalias tfplan='source /usr/local/bin/tfplan.sh'" >> ~/.bashrc
```
    
**zsh**:

```
echo -e "\nalias tfapply='source /usr/local/bin/tfapply.sh'" >> ~/.zshrc
echo -e "\nalias tfbackend='source /usr/local/bin/tfbackend.sh'" >> ~/.zshrc
echo -e "\nalias tfdestroy='source /usr/local/bin/tfdestroy.sh'" >> ~/.zshrc
echo -e "\nalias tfplan='source /usr/local/bin/tfplan.sh'" >> ~/.zshrc
```

***oh-my-zsh***:

```
echo -e "\nalias aws-switch='source /usr/local/bin/tfbackend.sh'" >> ~/.oh-my-zsh/custom/aliases.zsh
etc ...
```

Reload the shell
```
source ~/.bashrc
source ~/.zshrc
```

### Usage:

Using a new terminal window or tab (required for the new alias settings to take effect if you haven't reloaded the shell), simply run the script using any of the following aliases:

``` bash
tfapply
tfbackend
tfdestroy
tfplan
```

Example usage of the `tfbackend` script:


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

