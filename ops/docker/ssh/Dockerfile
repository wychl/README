FROM centos:latest

RUN yum -y update
RUN yum -y install openssh-server
RUN yum -y install openssh-clients

RUN mkdir /var/run/sshd
RUN ssh-keygen -t rsa -f /etc/ssh/ssh_host_rsa_key -N '' 

RUN echo 'root:password' | chpasswd

EXPOSE 22
ENTRYPOINT ["/usr/sbin/sshd", "-D"]