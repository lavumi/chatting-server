'use strict'

let source = null;

document.addEventListener('DOMContentLoaded', function () {
    checkLoginStatus();
    let elems = document.querySelectorAll('.sidenav');
    M.Sidenav.init(elems);
    var nameInputField = document.getElementById('name-input-field');
    M.CharacterCounter.init(nameInputField);
});


const checkLoginStatus = () => {
    let enterBtn = document.getElementById("enter-chat-btn");
    let quitBtn = document.getElementById("quit-chat-btn");
    let chatInput = document.getElementById("chat-input");
    let chatSend = document.getElementById("send-chat-btn");
    if (source !== null) {
        enterBtn.classList.add("disabled")
        quitBtn.classList.remove("disabled")
        chatInput.classList.remove("disabled")
        chatSend.classList.remove("disabled")
    } else {
        quitBtn.classList.add("disabled")
        chatInput.classList.add("disabled")
        chatSend.classList.add("disabled")
        enterBtn.classList.remove("disabled")
    }
}
const handleEnter = (event) => {
    if (event.keyCode === 13) {
        sendMsg();
    }
}
const sendMsg = () => {
    const input = document.getElementById('message');
    fetch('/api/chat/message', {
        method: "POST",
        credentials: 'include',
        headers: {"content-type": "application/json"},
        body: JSON.stringify(
            {
                sender: `${document.getElementById("name-input-field").value}`,
                msg: input.value
            })
    })
        .catch(err => console.log(err))
    input.value = '';
}
const loadUsers = () => {
    fetch('/api/chat/users')
        .then(response => {
            console.log("response : " + JSON.stringify(response));
            return response.json()
        })
        .then(json => {
            const userList = document.getElementById("user-container");
            for (let i = 0; i < json.userList.length; i++) {
                const user = document.createElement("li")
                user.classList.add("collection-item");
                user.innerHTML = json.userList[i];
                userList.appendChild(user);
            }
        });
}
const quitChat = () => {
    source.close();
    source = null;
    console.log('Chat closed');
    checkLoginStatus();
    resetChatting();
    resetUserList();
}
const resetChatting = () => {
    const myNode = document.getElementById("chat-list");
    while (myNode.firstChild) {
        myNode.removeChild(myNode.lastChild);
    }
}
const resetUserList = () => {
    const myNode = document.getElementById("user-container");
    while (myNode.firstChild) {
        myNode.removeChild(myNode.lastChild);
    }
}

let prev_sender = "";
const enterChat = () => {
    console.log('start sse')

    // document.cookie = `user=${document.getElementById("userName").innerHTML}`;

    let username = `${document.getElementById("name-input-field").value}`;
    if (username.length === 0) {
        return;
    }

    source = new EventSource("/api/chat/register?user=" + username);

    source.onerror = (e) => {
        console.log("EventSource failed", e);
    };

    source.addEventListener("info", (e) => {
        // console.log('sse info', e.data)
        let chatData = JSON.parse(e.data);

        console.log(chatData);
        console.log("username" + username);
        console.log("prev_sender" + prev_sender);
        let myCard = chatData.sender === username;

        console.log(myCard);
        console.log(prev_sender !== chatData.sender);

        if (myCard === false && prev_sender !== chatData.sender) {
            let namePanel = document.createElement("div");
            namePanel.className = "chat-user";
            namePanel.innerHTML = chatData.sender;
            prev_sender = chatData.sender;
            document.getElementById('chat-list').appendChild(namePanel);
        }
        var cardPanel = document.createElement("div");
        if (myCard === true) {
            cardPanel.className = "card-panel col s8 offset-s4 right-align teal lighten-2";
        } else {
            cardPanel.className = "card-panel col s8 left-align orange lighten-4";
        }

        cardPanel.style.borderRadius = "15px";
        cardPanel.innerText = chatData.msg;
        cardPanel.style.padding = "10px";
        cardPanel.style.marginBottom = "2px";


        document.getElementById('chat-list').appendChild(cardPanel);

        // chat.innerHTML += chatData.sender + "----" + chatData.msg  + '<br>';

    }, false);

    source.addEventListener("oper", (e) => {
        loadUsers();
        console.log('sse oper', e.data)
    }, false);

    checkLoginStatus();
    loadUsers();
}