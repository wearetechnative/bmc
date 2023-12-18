# TF BACKEND SWITCH
The program searches for all `*.tfbackend` files in a directory and let you choose which backend you want to configure

---

## To install/upgrade:

1. Download, optionally inspect, and copy the script to an appropriate folder
2. Add aliases:
    ```
    echo -e "\nalias aws-switch='source /usr/local/bin/tfbackend.sh'" >> ~/.bash_profile
    echo -e "\nalias aws-switch='source /usr/local/bin/tfbackend.sh'" >> ~/.zshrc
    ```

    For oh-my-zsh users:
    echo -e "\nalias aws-switch='source /usr/local/bin/tfbackend.sh'" >> ~/.oh-my-zsh/custom/aliases.zsh

Adding an alias to both config files is advised. Even if you only use one of the above shells, this will ensure that aws-switch.zsh works the same in either, should the need arise.

### To use:

Using a new terminal window or tab (required for the new alias settings to take effect), simply run the script using the alias created above (aws-switch.zsh):

```
aws-superstar@hackstation-[~]: tfbackend
