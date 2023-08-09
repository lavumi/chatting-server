'use strict'

let source = null;

const roomUri ="/api/chat/rooms";
const enterUri = "/api/chat/123/enter";
const msgUri = "/api/chat/123/message";
const userUri = "/api/chat/123/user";
// const enterUri = "";

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
    fetch(msgUri, {
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
    fetch(userUri)
        .then(response => {
            // console.log("response : " + JSON.stringify(response));
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
    // console.log('Chat closed');
    checkLoginStatus();
    resetChatting();
    resetUserList();
    toggleRoom();
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
const enterChat = (roomId) => {

    // document.cookie = `user=${document.getElementById("userName").innerHTML}`;

    let username = `${document.getElementById("name-input-field").value}`;
    if (username.length === 0) {
        return;
    }

    source = new EventSource(`/api/chat/${roomId}/enter`, {
        headers: {
            'Authorization': 'Bearer ' + "mytoken",
            'UserName': username
        }
    });

    source.onerror = (e) => {
        console.log("EventSource failed", e);
    };

    source.addEventListener("info", (e) => {
        // console.log('sse info', e.data)
        let chatData = JSON.parse(e.data);
        let myCard = chatData.sender === username;


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
    toggleRoom();
}

const getRoomInfo = ()=>{
    fetch(roomUri)
        .then(response => {
            // console.log("response : " + JSON.stringify(response.json()));
            return response.json()
        })
        .then(json => {
            for (let i = 0; i < json.roomList.length; i++) {
                // const user = document.createElement("li")
                // user.classList.add("collection-item");
                // user.innerHTML = json.userList[i];
                // userList.appendChild(user);
                createCard(json.roomList[i] , "1/5");
            }
            checkLoginStatus();
        });
}

function createCard(roomName, memberInfo) {
    const cardContainer = document.getElementById('room-div');

    const cardCol = document.createElement('div');
    cardCol.classList.add('col', 's3');

    const card = document.createElement('div');
    card.classList.add('card', 'blue', 'lighten-2', 'hoverable');

    const cardContent = document.createElement('div');
    cardContent.classList.add('card-content', 'white-text', 'unselectable');

    const cardTitle = document.createElement('span');
    cardTitle.classList.add('card-title');
    cardTitle.textContent = roomName;

    const memberInfoPara = document.createElement('p');
    memberInfoPara.textContent = memberInfo;

    cardContent.appendChild(cardTitle);
    cardContent.appendChild(memberInfoPara);
    card.appendChild(cardContent);


    card.addEventListener("click", () => {
        // console.log(roomName + "clicked");
        enterChat(roomName);
    })

    cardCol.appendChild(card);
    cardContainer.appendChild(cardCol);
}

function toggleRoom() {
    var topDiv = document.getElementById('room-div');
    var bottomDiv = document.getElementById('chat-div');

    topDiv.classList.toggle('expanded');
    topDiv.classList.toggle('collapsed');
    bottomDiv.classList.toggle('expanded');
    bottomDiv.classList.toggle('collapsed');
}