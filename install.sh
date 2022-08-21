#!/bin/sh

# install Chia harvester
basepath=$(cd `dirname $0`; pwd)
user=${USER}
myhome=${HOME}
userbin=$myhome/.local/bin
echo "Start..."

sudo -S install -Dm777 $basepath/modifyRDconfig $userbin/ &&
sed -i 's#testpath#'$basepath'#g' $basepath/modifyRDconfig.service &&
sudo -S install -Dm644 $basepath/modifyRDconfig.service /usr/lib/systemd/system/modifyRDconfig.service &&
sudo -S systemctl enable modifyRDconfig &&
sudo -S systemctl daemon-reload &&
sudo -S systemctl start modifyRDconfig &&
echo "install complete"
