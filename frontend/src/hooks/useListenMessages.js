import { useEffect } from "react";

import { useSocketContext } from "../context/SocketContext";
import useConversation from "../zustand/useConversation";

import notificationSound from "../assets/sounds/notification.mp3";

const useListenMessages = () => {
	const { socket } = useSocketContext();
	const { messages, setMessages } = useConversation();

	useEffect(() => {
		socket.onmessage = function (event) {
			event.shouldShake = true;
			const sound = new Audio(notificationSound);
			sound.play();
			setMessages([...messages, JSON.parse(event.data)]);
		};


		// return () => socket?.off("newMessage");
	}, [socket, setMessages, messages]);
};
export default useListenMessages;
