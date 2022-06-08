# task_pace

## creating the chat server using golang 
task done till now :-- 

1. creating the api for chat server using socket.io liberary
2. any user can join the and having an unique id 
3. user should have mention ther username ... 
4. any user can create a chat room with other member or without other member
5. a particular user also joined the chat room wit using group id 
6. maximum 5 member including admin allowed to joining the chat room 
7. group id is created using the UUID no.

## working principal ..
1. if any user joined our chat server all user get datils of new user using  event listerner name <b>"allconnectedusers"</b>
2. if any user create a group with particular member the the mentioned member should also get a notification using event linner <b>"notification"</b>
3. if any user joined the group with group id all user in that group also get a notification about new member using event linner <b>"notification"</b>
4. if no of user geting max then 5 in any group a waring message should get by the user useing event linner <b>"warning"</b>

## .env
create a .env file and add this ::--

PORT = :5000
