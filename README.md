Installation instruction:

1. Step: run following command with root rules
$ sudo sh ./install.sh

2. Step: Download and install Node.js
https://nodejs.org/en/download/package-manager/

3. Step: Install pm2.
$ npm init
$ npm install pm2 -g

4. Step: Configure config.toml file

5. Step: run app with pm2
$ pm2 start ./app

Examples of config parameter
uploads_dir="http://54.201.86.220/api/uploads"

Env vars:
SUB_PATH="/api" 
CONFIG_PATH="*"
UPLOADS_PATH="http://54.201.86.220/api/uploads"
LOAD_FIXTURES=1
PORT=0000

CONFIG_PATH: Path to config file. If it is set as "*" than config file should be placed in same directory as app. Or you can specify any different path fully.
UPLOADS_PATH: URL that will be atteached to notification images path