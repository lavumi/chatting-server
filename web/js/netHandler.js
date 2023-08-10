'use strict'

let source = null;
let userInfo = null;
let header = null;
let current_room_uuid = '';

function _setHeader(userName) {
    userInfo = {
        Authorization: 'Bearer ' + "mytoken",
        UserName: userName
    }
    header = new Headers({
        'Authorization': 'Bearer ' + "mytoken",
        'UserName': userName
    });
}

function _enterChatRoom(uuid, msgListener) {
    source = new EventSource(`/api/chat/${uuid}/enter`, userInfo);
    source.onerror = (e) => {
        console.log("EventSource failed", e);
    };
    source.addEventListener('msg', e => msgListener('msg', e), false);
    source.addEventListener('event', e => msgListener('event', e), false);

    current_room_uuid = uuid;
}

function _sendMessage(message) {
    if (source === null || current_room_uuid === '') {
        return;
    }

    fetch(`/api/chat/${current_room_uuid}/message`,
        {
            method: "POST",
            headers: header,
            body: message
        })
        .catch(err => console.log(err))
}

async function _getRoomInfo() {
    if (source === null || current_room_uuid === '') {
        return;
    }
    return await fetch(`/api/chat/${current_room_uuid}/info`,
        {
            method: "GET",
            headers: header,
        })
        .then(response => {

            return response.json()
        })
        .then(json => {
            console.log(json);
            return json.roomInfo
        })
        .catch(err => console.log(err))
}

function _exitChatRoom() {
    source.close();
    source = null;
}

async function _getRoomList() {
    return await fetch("/api/chat/rooms",
        {
            method: "GET",
            headers: header,
        })
        .then(response => {

            return response.json()
        })
        .then(json => json.roomList)
        .catch(err => console.log(err))
}

let Chat = {
    Register: _setHeader,
    List: _getRoomList,
    Enter: _enterChatRoom,
    Exit: _exitChatRoom,
    Send: _sendMessage,
    Info: _getRoomInfo,
}

