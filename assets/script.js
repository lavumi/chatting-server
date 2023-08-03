'use strict'

console.log('Script loaded');
let source = null;

window.onload = () => {
    const chatUser = `user-${Math.random().toString(36).substring(7)}`;
    document.getElementById("userName").innerHTML = chatUser;

    if (source) loadUsers();
}

const sendMsg = () => {
    const input = document.getElementById('msg');
    fetch('/api/chat/message', { method: "POST", credentials: 'include', headers: { "content-type": "application/json" }, body : JSON.stringify({ msg: input.value }) })
        .catch(err => console.log(err))
    document.getElementById('msg').value = '';
}

const loadUsers = () => {
    fetch('/api/chat/users')
        .then(response => response.json())
        .then(json => {
            document.getElementById("userList").innerHTML = json.userList[0] ? json.userList.join('<br>') : 'No user';
        });
}

const quitChat = () => {
    source.close();
    console.log('Chat closed');
    document.cookie = `user=`;
    document.getElementById('userList').innerHTML = '';
    document.getElementById('chat').innerHTML = '';
}

const enterChat = () => {
    console.log('start sse')

    // document.cookie = `user=${document.getElementById("userName").innerHTML}`;

    let username = `${document.getElementById("userName").innerHTML}`;

    source = new EventSource("/api/chat/register?user="+username);

    source.onerror = (e) => {
        console.log("EventSource failed", e);
    };

    source.addEventListener("info", (e) => {
        const chat = document.getElementById('chat');
        chat.innerHTML += e.data + '<br>';
        console.log('sse info', e.data)
    }, false);

    source.addEventListener("oper", (e) => {
        loadUsers();
        console.log('sse oper', e.data)
    }, false);

}