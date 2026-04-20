#!/bin/bash

set -e 

create_user(){
    bash scripts/create_user.sh

    current_user="$SUDO_USER"

    usermod -a -G portfolio $current_user
}

configure_systemd_service(){
	cp ./init/portfolio.service /etc/systemd/system/
	systemctl daemon-reload
	systemctl enable portfolio.service
	systemctl restart portfolio.service
	systemctl status portfolio.service
}

set_file_permissions(){
    chown portfolio:portfolio /usr/local/bin/portfolio
    chmod 770 /usr/local/bin/portfolio*

    chown -R portfolio:portfolio /var/lib/portfolio
    chmod -R 770 /var/lib/portfolio/posts
    chmod -R 770 /var/lib/portfolio/web
    chmod 660 /var/lib/portfolio/config.ini
}

main(){
    create_user
    set_file_permissions
    configure_systemd_service
}

main
