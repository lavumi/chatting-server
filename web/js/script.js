'use strict'

console.log('Script loaded');
let source = null;

window.onload = () => {
    document.getElementById("userName").innerHTML = `user-${Math.random().toString(36).substring(7)}`;

    if (source) loadUsers();
}

const sendMsg = () => {
    const input = document.getElementById('msg');
    fetch('/api/chat/message', {
        method: "POST",
        credentials: 'include',
        headers: {"content-type": "application/json"},
        body: JSON.stringify(
            {
                sender: `${document.getElementById("userName").innerHTML}`,
                msg: input.value
            })
    })
        .catch(err => console.log(err))
    document.getElementById('msg').value = '';
}

const loadUsers = () => {
    fetch('/api/chat/users')
        .then(response => {
            console.log("response : " + JSON.stringify(response));
            return response.json()
        })
        .then(json => {
            // console.log("json ---" + JSON.stringify(json))
            // console.log(document);
            // console.log(document.getElementById("user-list"));
            const userList = document.getElementById("user-list");
            for (let i = 0; i < json.userList.length; i++) {

                const user = document.createElement("li")
                user.classList.add("list-group-item");
                user.innerHTML = json.userList[i];
                userList.appendChild(user);
            }


            // document.getElementById("userList").innerHTML = json.userList[0] ? json.userList.join('<br>') : 'No user';
        });
}

const quitChat = () => {
    source.close();
    console.log('Chat closed');
    document.cookie = `user=`;
    document.getElementById('userList').innerHTML = '';
    document.getElementById('chat').innerHTML = '';
}

const addChatBoxBootstrap = (username, chatData) => {
    const chatCard = document.createElement("div");
    chatCard.className = "row m-3";
    const nameCard = document.createElement("div")
    nameCard.className = "col-sm-4 p-3 text-right";
    nameCard.innerHTML = chatData.sender;


    const messageCard = document.createElement("div")
    messageCard.className = "col-sm-8 card";
    const msgBody = document.createElement("div");
    msgBody.className = "card-body";
    msgBody.innerHTML = chatData.msg;
    messageCard.appendChild(msgBody);


    if (chatData.sender === username) {
        // nameCard.className += " text-right";
        nameCard.style.textAlign = "right";
        chatCard.appendChild(messageCard);
        chatCard.appendChild(nameCard);
    } else {
        chatCard.appendChild(nameCard);
        chatCard.appendChild(messageCard);
    }

    return chatCard;
}

const addChatBoxMaterialize = (username, chatData) => {
    const chatCard = document.createElement("div");
    chatCard.className = "row m-3";
    const nameCard = document.createElement("div")
    nameCard.className = "col s4 p-3";
    nameCard.innerHTML = chatData.sender;


    const messageCard = document.createElement("div")
    messageCard.className = "col s8 card";
    const msgBody = document.createElement("div");
    msgBody.className = "card-content";
    msgBody.innerHTML = chatData.msg;
    messageCard.appendChild(msgBody);


    if (chatData.sender === username) {
        // nameCard.className += " text-right";
        nameCard.className += " right-align";
        chatCard.appendChild(messageCard);
        chatCard.appendChild(nameCard);
    } else {
        chatCard.appendChild(nameCard);
        chatCard.appendChild(messageCard);
    }

    return chatCard;
}

const enterChat = () => {
    console.log('start sse')

    // document.cookie = `user=${document.getElementById("userName").innerHTML}`;

    let username = `${document.getElementById("userName").innerHTML}`;

    source = new EventSource("/api/chat/register?user=" + username);

    source.onerror = (e) => {
        console.log("EventSource failed", e);
    };

    source.addEventListener("info", (e) => {
        const chat = document.getElementById('chat-list');
        // console.log('sse info', e.data)
        let chatData = JSON.parse(e.data);

        let chatCard = addChatBoxMaterialize(username, chatData);



        chat.appendChild(chatCard);

        // chat.innerHTML += chatData.sender + "----" + chatData.msg  + '<br>';

    }, false);

    source.addEventListener("oper", (e) => {
        loadUsers();
        console.log('sse oper', e.data)
    }, false);

    loadUsers();
}