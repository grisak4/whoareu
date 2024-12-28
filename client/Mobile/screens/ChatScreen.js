import React, { useState, useEffect, useRef, memo } from 'react';
import { View, FlatList, Text, TextInput, TouchableOpacity, StyleSheet, Alert, Image, KeyboardAvoidingView, Platform, ActivityIndicator } from 'react-native';
import AsyncStorage from '@react-native-async-storage/async-storage';
import moment from 'moment';
import { Ionicons } from '@expo/vector-icons';
import { jwtDecode } from "jwt-decode";

const ChatScreen = ({ route, navigation }) => {
  const { chatId, chatName, avatarUrl } = route.params;
  const [messages, setMessages] = useState([]);
  const [newMessage, setNewMessage] = useState('');
  const [userSenderId, setUserSenderId] = useState(null);
  const [userNickname, setUserNickname] = useState(null);
  const [loading, setLoading] = useState(true);
  const flatListRef = useRef(null);
  const ws = useRef(null);
  const reconnectAttempts = useRef(0);

  const scrollToBottom = () => {
    setTimeout(() => {
      flatListRef.current?.scrollToEnd({ animated: true });
    }, 100);
  };

  useEffect(() => {
    let isMounted = true;
  
    const fetchMessages = async () => {
      try {
        const response = await fetch(`https://beagle-mighty-terribly.ngrok-free.app/api/v1/chats/getmessages/${chatId}`);
        const data = await response.json();
        if (isMounted) {
          setMessages(data.data || []);
          scrollToBottom();
        }
      } catch (error) {
        Alert.alert('Ошибка', 'Не удалось загрузить сообщения');
      } finally {
        if (isMounted) {
          setLoading(false);
        }
      }
    };
  
    const fetchUserId = async () => {
      try {
        const token = await AsyncStorage.getItem('token');
        if (!token) {
          Alert.alert('Error', 'JWT token not found');
          return;
        }
  
        const decodedToken = jwtDecode(token);
        const userId = decodedToken.user_id;
        const nickname = decodedToken?.UserConf?.info?.nickname;
  
        if (!userId || !nickname) {
          Alert.alert('Error', 'Required user information not found in token');
          return;
        }
  
        if (isMounted) {
          setUserSenderId(userId);
          setUserNickname(nickname);
        }
  
        return { userId, nickname };
      } catch (error) {
        console.error('Ошибка получения ID пользователя:', error);
      }
    };
  
    const setupWebSocket = async () => {
      const { userId } = await fetchUserId();
      if (userId && isMounted) {
        ws.current = new WebSocket(`wss://beagle-mighty-terribly.ngrok-free.app/api/v1/ws/startchat/${chatId}/${userId}`);
  
        ws.current.onopen = () => {
          console.log('WebSocket соединение открыто');
          reconnectAttempts.current = 0;
        };
  
        ws.current.onmessage = (event) => {
          try {
            const message = JSON.parse(event.data);
            if (isMounted) {
              setMessages((prevMessages) => [...prevMessages, message]);
              scrollToBottom();
            }
          } catch (error) {
            console.error('Ошибка разбора сообщения:', error);
          }
        };
  
        ws.current.onclose = () => {
          console.log('WebSocket соединение закрыто');
          if (isMounted) attemptReconnect();
        };
  
        ws.current.onerror = (error) => {
          console.error('WebSocket ошибка:', error.message || 'Неизвестная ошибка');
        };
      }
    };
  
    const attemptReconnect = () => {
      if (reconnectAttempts.current < 5 && isMounted) {
        reconnectAttempts.current += 1;
        setTimeout(() => {
          console.log(`Попытка переподключения ${reconnectAttempts.current}`);
          setupWebSocket();
        }, 2000 * reconnectAttempts.current);
      } else if (!isMounted) {
        console.log('Компонент размонтирован, переподключение остановлено.');
      } else {
        Alert.alert('Ошибка', 'Не удалось подключиться к WebSocket после нескольких попыток. Пожалуйста, перезагрузите чат.');
      }
    };
  
    fetchMessages();
    setupWebSocket();
  
    return () => {
      isMounted = false;
      if (ws.current) {
        ws.current.onopen = null;
        ws.current.onmessage = null;
        ws.current.onclose = null;
        ws.current.onerror = null;
        ws.current.close();
        console.log('WebSocket закрыт при размонтировании компонента.');
      }
    };
  }, [chatId]);  

  const handleSendMessage = () => {
    if (newMessage.trim() === '' || !userSenderId || !userNickname) return;

    const message = {
      chat_id: chatId,
      user_id: userSenderId,
      message_content: newMessage,
      sender_name: userNickname,
    };

    if (ws.current && ws.current.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify(message));

      setMessages((prevMessages) => [
        ...prevMessages,
        { ...message, message_id: Date.now(), time_sent: new Date() },
      ]);

      setNewMessage('');
      scrollToBottom();
    } else {
      console.error('WebSocket не открыт. Текущее состояние:', ws.current?.readyState);
    }
  };

  const MessageItem = memo(({ item }) => (
    <View style={[styles.messageItem, userSenderId === item.user_id ? styles.ownMessage : styles.otherMessage]}>
      <Text style={styles.sender}>
        {item.sender_name ? `${item.sender_name}` : `Неизвестный Пользователь`}
      </Text>
      <Text style={styles.message}>{item.message_content}</Text>
      <Text style={styles.time}>{moment(item.time_sent).format('HH:mm DD/MM/YYYY')}</Text>
    </View>
  ));

  return (
    <KeyboardAvoidingView
      style={styles.container}
      behavior={Platform.OS === 'ios' ? 'padding' : undefined}
    >
      <View style={styles.chatHeader}>
        <TouchableOpacity onPress={() => navigation.goBack()}>
          <Ionicons name="arrow-back" size={24} color="#1e90ff" />
        </TouchableOpacity>
        <Image source={{ uri: avatarUrl || 'https://i.pinimg.com/736x/97/29/79/972979190503cbd54ac183b47872ae80.jpg' }} style={styles.avatar} />
        <Text style={styles.chatTitle}>{chatName}</Text>
      </View>

      {loading ? (
        <ActivityIndicator size="large" color="#1e90ff" />
      ) : (
        <FlatList
          ref={flatListRef}
          data={messages}
          keyExtractor={(item, index) => (item.message_id ? item.message_id.toString() : index.toString())}
          renderItem={({ item }) => <MessageItem item={item} />}
          onContentSizeChange={scrollToBottom}
          onLayout={scrollToBottom}
        />
      )}

      <View style={styles.inputContainer}>
        <TextInput
          style={styles.input}
          value={newMessage}
          onChangeText={setNewMessage}
          placeholder="Введите сообщение"
        />
        <TouchableOpacity style={styles.sendButton} onPress={handleSendMessage}>
          <Ionicons name="send-outline" size={24} color="#fff" />
        </TouchableOpacity>
      </View>
    </KeyboardAvoidingView>
  );
};

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: '#f9f9f9' },
  chatHeader: { flexDirection: 'row', alignItems: 'center', padding: 10, backgroundColor: '#fff', borderBottomWidth: 1, borderBottomColor: '#eee' },
  avatar: { width: 40, height: 40, borderRadius: 20, marginLeft: 10 },
  chatTitle: { fontSize: 18, fontWeight: 'bold', marginLeft: 10, color: '#333' },
  messageItem: { padding: 10, marginVertical: 5, borderRadius: 8, maxWidth: '80%', alignSelf: 'flex-start' },
  ownMessage: { backgroundColor: '#DCF8C6', alignSelf: 'flex-end' },
  otherMessage: { backgroundColor: '#F1F0F0' },
  sender: { fontWeight: 'bold', marginBottom: 4, color: '#1e90ff' },
  message: { fontSize: 16, marginBottom: 4 },
  time: { fontSize: 12, color: '#aaa', textAlign: 'right' },
  inputContainer: { flexDirection: 'row', alignItems: 'center', padding: 10, borderTopWidth: 1, borderTopColor: '#eee', backgroundColor: '#fff' },
  input: { flex: 1, height: 40, borderColor: 'gray', borderWidth: 1, borderRadius: 20, paddingHorizontal: 15, marginRight: 10 },
  sendButton: { backgroundColor: '#1e90ff', borderRadius: 20, padding: 10 },
});

export default ChatScreen;
