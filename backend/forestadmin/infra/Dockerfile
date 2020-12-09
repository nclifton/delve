FROM keymetrics/pm2:latest-stretch
RUN mkdir -p /forestadmin
ADD backend/forestadmin/package.json /forestadmin/package.json
ADD backend/forestadmin/package-lock.json /forestadmin/package-lock.json
RUN apt-get install python g++ gcc make tzdata
WORKDIR /forestadmin
RUN npm install --production --loglevel verbose
EXPOSE 3310

ADD backend/forestadmin/ /forestadmin

ENTRYPOINT ["pm2-runtime", "server.js"]
