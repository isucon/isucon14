FROM mysql:8

COPY my.cnf /etc/mysql/conf.d/my.cnf
RUN chmod 644 /etc/mysql/conf.d/my.cnf
RUN touch /var/log/mysql_error.log
RUN chmod 777 /var/log/mysql_error.log
RUN chmod 777 /var/log
RUN chmod 777 /var
