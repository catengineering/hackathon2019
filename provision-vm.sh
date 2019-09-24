#!/bin/bash

# install stuff
#snap install --classic go

# download the script to run the app from central storage account
/usr/bin/wget -O /opt/get-app.sh https://gobibeareaststorage.blob.core.windows.net/staging/get-app.sh
chmod +x /opt/get-app.sh

# configure the script to auto-run on reboot
echo "#!/bin/sh -e" >/etc/rc.local
echo /opt/get-app.sh >>/etc/rc.local
echo exit 0 >>/etc/rc.local
chmod +x /etc/rc.local

# start the app and leave running
nohup /opt/get-app.sh &
