# Configuring SmarterSmarterChild With systemd

This document details the configuration of SmarterSmarterChild to run as an unprivileged user with `systemd` managing it
as a production service.

1. ** Download SmarterSmarterChild**

   Grab the latest Linux release from the [releases page](https://github.com/mk6i/smarter-smarter-child/releases)

2. ** Create the ssc user and group **

   Run the following commands:

   ```shell
   $ sudo useradd ssc
   $ sudo mkdir -p /opt/ssc
   $ sudo mkdir -p /var/ssc
   ```

3. ** Extract the archive **

   Extract the archive using the usual `tar` invocation, and move the extracted contents into `/opt/ssc`

4. ** Set Ownership and Permissions **

   ```shell
   $ sudo chown -R ssc:ssc /opt/ssc
   $ sudo chmod -R o-rx /opt/ssc
   ```

5. ** Copy the systemd service **

   Place the `ssc.service` file in `/etc/systemd/system`

6. ** Reload systemd **

   ```shell
   $ sudo systemctl daemon-reload
   ```

7. ** Enable and start the service **

  ```shell
  $ sudo systemctl enable --now ssc.service
  ```

8. ** Make sure the service is running **

   ```shell
   $ sudo systemctl status ssc.service
   $ sudo journalctl -xeu ssc.service
   ```

Note that the `systemd` service defines the configuration for SmarterSmarterChild directly, bypassing the usual `run.sh`
script and `settings.env`. Customizations may be performed in `/etc/systemd/system/ssc.service`.
