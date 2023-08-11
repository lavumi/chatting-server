'use strict'

document.addEventListener('DOMContentLoaded', function () {

    let elems = document.querySelectorAll('.sidenav');
    M.Sidenav.init(elems);
    var nameInputField = document.getElementById('name-input-field');
    M.CharacterCounter.init(nameInputField);

    updateLoginStatus(false);
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
            createRoomCard(name, uuid, desc);
        }
        updateLoginStatus(true);
    })();
}
const QuitChat = () => {
    Chat.Exit();
    updateLoginStatus(false);
    resetRoomList();
    resetChatting();
    resetUserList();
    // toggleRoom(true);
}
const SendMessage = () => {
    const input = document.getElementById('message');

    if (input.value.trim().length !== 0)
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
        let continueTalking = prev_sender === chatData.sender;
        prev_sender = chatData.sender;


        if (continueTalking === true) {
            let lastChatPanel = document.getElementById('chat-list').lastElementChild;
            let lastChatContextCard = lastChatPanel.lastElementChild;
            lastChatContextCard.innerText += "\n" + chatData.msg;
            document.getElementById('chat-list').scrollIntoView(false);
            return;
        }

        let chatPanel = document.createElement("div");
        let chatContextCard = document.createElement("div");


        chatPanel.style.marginTop = "10px";
        if (myCard === true) {
            chatPanel.className = "col s8 offset-s4 my-chat-panel";
            chatContextCard.className = "card-panel small left-align teal lighten-2 chat-box";
        } else {
            let namePanel = document.createElement("div");
            namePanel.className = "chat-user";
            namePanel.innerHTML = chatData.sender;
            chatPanel.appendChild(namePanel);
            chatPanel.className = "col s8 left-align";
            chatContextCard.className = "card-panel small left-align orange lighten-4 chat-box";
        }

        chatContextCard.innerText = chatData.msg;

        chatPanel.appendChild(chatContextCard);

        let chatList = document.getElementById('chat-list');
        chatList.appendChild(chatPanel);

        document.getElementById('chat-list').scrollIntoView(false);
    }

    function eventHandler(event) {
        loadUsers();
        console.log('sse event', event.data)
    }

    resetChatting();


    (async () => {
        let chatRoomInfo = await Chat.Info();
        let messageData = chatRoomInfo.messages;
        for (const messageDatum of messageData) {
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

function createRoomCard(roomName, uuid, desc) {
    const cardContainer = document.getElementById('room-div');

    const cardCol = document.createElement('div');
    // cardCol.classList.add('col', 's3');

    const card = document.createElement('div');
    card.classList.add('card-panel', 'blue', 'lighten-2', 'hoverable');
    const cardTitle = document.createElement('span');
    cardTitle.classList.add('white-text');
    cardTitle.textContent = roomName;
    card.appendChild(cardTitle);

    // const cardContent = document.createElement('div');
    // cardContent.classList.add('card-content', 'white-text', 'unselectable');
    //
    // const cardTitle = document.createElement('span');
    // cardTitle.classList.add('card-title');
    // cardTitle.textContent = roomName;

    // const memberInfoPara = document.createElement('p');
    // memberInfoPara.textContent = "Empty";

    // cardContent.appendChild(cardTitle);
    // cardContent.appendChild(memberInfoPara);
    // card.appendChild(cardContent);

    M.Tooltip.init(card, {html: desc, position: "right"});
    card.addEventListener("click", () => {
        // console.log(roomName + "clicked");
        enterChat(uuid);
    })

    cardCol.appendChild(card);
    cardContainer.appendChild(cardCol);
}


function resetRoomList() {
    const myNode = document.getElementById('room-div');
    while (myNode.firstChild) {
        myNode.removeChild(myNode.lastChild);
    }
}

function resetChatting() {
    const myNode = document.getElementById("chat-list");
    while (myNode.firstChild) {
        myNode.removeChild(myNode.lastChild);
    }
}

function resetUserList() {
    // const myNode = document.getElementById("user-container");
    // while (myNode.firstChild) {
    //     myNode.removeChild(myNode.lastChild);
    // }
}

function updateLoginStatus(loggedIn) {
    let enterBtn = document.getElementById("enter-chat-btn");
    let quitBtn = document.getElementById("quit-chat-btn");
    let chatInput = document.getElementById("chat-input");
    let chatSend = document.getElementById("send-chat-btn");
    let nameField = document.getElementById("name-input-field");
    let chatField = document.getElementById("chat-input");
    if (loggedIn) {
        enterBtn.classList.add("disabled")
        nameField.classList.add("disabled")

        quitBtn.classList.remove("disabled")
        chatInput.classList.remove("disabled")
        chatSend.classList.remove("disabled")

        chatField.classList.remove("disabled")
    } else {
        quitBtn.classList.add("disabled")
        chatInput.classList.add("disabled")
        chatSend.classList.add("disabled")

        chatField.classList.add("disabled")


        enterBtn.classList.remove("disabled")
        nameField.classList.remove("disabled")

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