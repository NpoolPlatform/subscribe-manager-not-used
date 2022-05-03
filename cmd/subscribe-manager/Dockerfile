FROM jrei/systemd-centos:7

COPY SubscribeManager.viper.yaml /etc/SubscribeManager/SubscribeManager.viper.yaml
COPY subscribe-manager /usr/local/bin/subscribe-manager

CMD ["subscribe-manager", "run"]
