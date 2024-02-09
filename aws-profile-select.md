# AWS Profile Select Tool

## A painless way to select an AWS profile

This script scans your aws configuration for profile names, and allows you to choose them by number, because messing with environment variables repeatedly is toil. Toil sucks.

It's also a handy way to see the currently selected profile, as it is given in an informational message above the selection menu. You may press crtl+c at any time to exit the tool.

---

### Prerequistites

##### Shell compatibility

aws-switch.zsh has been tested to be compatible with the following shell versions:

- Bash: v4 and newer
- ZSH: tested with v5.8 (with [Oh My Zsh!](https://github.com/ohmyzsh/ohmyzsh/wiki) installed, but should work equally well without it)

It may work with shells sharing compatibility with the above, but it's certainly not guaranteed.

##### AWS CLI

You will need to have your AWS CLI configured with profiles before this tool does anything useful. For more information on configuring AWS named profiles, see the documentation here: https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-profiles.html . This tool works with the named profiles in the `~.aws/config` file, as shown in the above link's second example.

For configuring different AWS profiles using a single IAM account and MFA, see the documentation here: https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-role.html

<details>
<summary>Example config file</summary>

```bash
[company-userauth]
region = eu-central-1
source_profile=company-userauth

[profile company-playground]
region = eu-central-1
role_arn = arn:aws:iam::123456789012:role/landing_zone_devops_administrator
source_profile = company-userauth

[profile company-playground-mgt]
region = eu-central-1
role_arn = arn:aws:iam::12345678912:role/landing_zone_devops_administrator
source_profile = company-userauth
```

</details>

<details>
<summary>Example credential file</summary>

```bash
[company-userauth-long-term]
aws_access_key_id = <AWS-CREDENTIALS-KEY-ID>
aws_secret_access_key = <AWS-CREDENTIALS-SECRET-ACCESS-KEY>
aws_mfa_device = arn:aws:iam::123456789012:mfa/<MFA-DEVICE-ALIAS>

[company-userauth]
aws_access_key_id = ASIAXSZQFYVIG374RTHB
aws_secret_access_key = 9Kdfk8SUICbA+5izT/oKZx9LODSQ7DmYLXiu/Z3U
assumed_role = False
aws_security_token =
aws_session_token =
expiration =

```

</details>

---

### To install/upgrade:

_Note_: the procedures might be a little different, depending on your personal configuration

1. Download or copy the script to "/usr/local/bin", add it to your PATH, and make it executable with `chmod +x /usr/local/bin/aws-profile-select.sh`:
2. Add aliases:

   For Ubuntu users (if alias not work use ~/.profile instead of ~/.bashrc):
   ```
   echo -e "\nalias aws-switch='source /usr/local/bin/aws-profile-select.sh'" >> ~/.bashrc
   echo -e "\nalias aws-switch='source /usr/local/bin/aws-profile-select.sh'" >> ~/.zshrc
   ```

   For oh-my-zsh users:

   ```
   echo -e "\nalias aws-switch='source /usr/local/bin/aws-profile-select.sh'" >> ~/.oh-my-zsh/custom/aliases.zsh
   ```

   Reload the shell
   ```
   source ~/.bashrc
   source ~/.zshrc
   ```

Adding an alias to both config files is advised. Even if you only use one of the above shells, this will ensure that aws-switch.zsh works the same in either, should the need arise.

### To use:

Using a new terminal window or tab (required for the new alias settings to take effect if you haven't reloaded the shell), simply run the script using the alias created above `aws-switch`:

```
aws-superstar@hackstation-[~]: aws-switch

------------- AWS Profile Select-O-Matic -------------
No profile set yet

Type the number of the profile you want to use from the list below, and press enter

-: Unset Profile
0: default
1: personal
2: company-main

Selection:
```

Typing `2` and pressing enter will make this terminal use the selected profile until you re-run this command and select another profile.

```
Selection: 2
Activating profile 2: company-main
aws-superstar@hackstation-[~]: (company-main):
```

Not necessary to run, just proof this thing does what it says:

```
aws-superstar@hackstation[~]: (company-main): echo $AWS_PROFILE
company-main
```

##### aws-mfa

To use `aws-mfa`, you need to:

- install the tool manually.
- set variable `aws_mfa`

**General installation**:

See [https://github.com/broamski/aws-mfa](https://github.com/broamski/aws-mfa)

**Homebrew installation**:

```
 brew install aws-iam-authenticator
```

###### Wouter van der Toorren. Forked from: Jesse Price
