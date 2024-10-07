# telegram_helpdesk_bot_go
MVP of telegram bot for helpdesk on Go

Bot can create tickets and send them to helpdesk of departments

### How to use :
    * Write to ./deploy/start.sh  your BOTTOKEN
    * Launch Redis in Docker with : docker-compose -f docker-compose.yml up -d --remove-orphans
    * ./start.sh
    * You have fill admins of departments : just write password of department to Bot. You can find passwords in start.sh
    * Then you can make tickets