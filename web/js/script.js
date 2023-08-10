'use strict'

document.addEventListener('DOMContentLoaded', function () {

    let elems = document.querySelectorAll('.sidenav');
    M.Sidenav.init(elems);
    var nameInputField = document.getElementById('name-input-field');
    M.CharacterCounter.init(nameInputField);

    updateLoginStatus();
});


const Login = () => {
    let name = document.getElementById("name-input-field").value;
    Chat.Register(name);
    (async () => {
        let roomList = await Chat.List();
        for (const roomInfo of roomList) {
            let name = roomInfo.name;
            let uuid = roomInfo.uuid;
            let desc = roomInfo.desc;
            createCard(name, uuid, desc);
        }
        updateLoginStatus();
    })();
}
const QuitChat = () => {
    Chat.Exit();
    updateLoginStatus();
    resetChatting();
    resetUserList();
    toggleRoom(true);
}
const SendMessage = () => {
    const input = document.getElementById('message');
    Chat.Send(input.value);
    input.value = '';
}
const HandleEnter = (event) => {
    if (event.keyCode === 13) {
        SendMessage();
    }
}


let prev_sender = "";

function enterChat(roomId) {

    // document.cookie = `user=${document.getElementById("userName").innerHTML}`;

    let username = `${document.getElementById("name-input-field").value}`;
    if (username.length === 0) {
        return;
    }


    Chat.Enter(roomId, (type, event) => {
        switch (type) {
            case 'msg':
                messageHandler(event.data);
                break;
            case 'event':
                eventHandler(event.data);
        }
    })

    function messageHandler(data) {
        let chatData = JSON.parse(data);
        let myCard = chatData.sender === username;


        if (myCard === false && prev_sender !== chatData.sender) {
            let namePanel = document.createElement("div");
            namePanel.className = "chat-user";
            namePanel.innerHTML = chatData.sender;
            prev_sender = chatData.sender;
            document.getElementById('chat-list').appendChild(namePanel);
        }
        let cardPanel = document.createElement("div");
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
    }

    function eventHandler(event) {
        loadUsers();
        console.log('sse event', event.data)
    }

    updateLoginStatus();
    loadUsers();
    toggleRoom(false);


    (async () => {
        let chatRoomInfo = await Chat.Info();
        console.log("chatRoomInfo", chatRoomInfo);

        let messageData = chatRoomInfo.messages;

        for (const messageDatum of messageData) {
            console.log(messageDatum);
            messageHandler(messageDatum);
        }

    })();

}

function loadUsers() {
    // fetch(userUri)
    //     .then(response => {
    //         // console.log("response : " + JSON.stringify(response));
    //         return response.json()
    //     })
    //     .then(json => {
    //         const userList = document.getElementById("user-container");
    //         for (let i = 0; i < json.userList.length; i++) {
    //             const user = document.createElement("li")
    //             user.classList.add("collection-item");
    //             user.innerHTML = json.userList[i];
    //             userList.appendChild(user);
    //         }
    //     });
}


function createCard(roomName, uuid, desc) {
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
    memberInfoPara.textContent = "Empty";

    cardContent.appendChild(cardTitle);
    cardContent.appendChild(memberInfoPara);
    card.appendChild(cardContent);

    M.Tooltip.init(card, {html: desc});
    card.addEventListener("click", () => {
        // console.log(roomName + "clicked");
        enterChat(uuid);
    })

    cardCol.appendChild(card);
    cardContainer.appendChild(cardCol);
}

function resetChatting() {
    const myNode = document.getElementById("chat-list");
    while (myNode.firstChild) {
        myNode.removeChild(myNode.lastChild);
    }
}

function resetUserList() {
    const myNode = document.getElementById("user-container");
    while (myNode.firstChild) {
        myNode.removeChild(myNode.lastChild);
    }
}

function updateLoginStatus() {
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

function toggleRoom(isRoom) {
    let room = document.getElementById('room-div');
    let chat = document.getElementById('chat-div');


    if (isRoom === true) {
        room.classList.add('expanded');
        room.classList.remove('collapsed');
        chat.classList.add('collapsed');
        chat.classList.remove('expanded');
    } else {
        room.classList.remove('expanded');
        room.classList.add('collapsed');
        chat.classList.remove('collapsed');
        chat.classList.add('expanded');
    }
}