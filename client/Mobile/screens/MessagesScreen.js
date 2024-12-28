import React, { useState, useEffect } from 'react';
import { View, FlatList, Text, TouchableOpacity, StyleSheet, Alert, Image, TouchableWithoutFeedback } from 'react-native';
import { getChats } from '../utils/api';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { Ionicons } from '@expo/vector-icons';
import { jwtDecode } from "jwt-decode";

export default function MessagesScreen({ navigation }) {
  const [chats, setChats] = useState([]);

  useEffect(() => {
    const fetchChats = async () => {
      const token = await AsyncStorage.getItem('token');
    
      if (!token) {
        Alert.alert('Error', 'JWT token not found');
        return;
      }
    
      const decodedToken = jwtDecode(token);
    
      const userId = decodedToken.user_id;
      if (!userId) {
        Alert.alert('Error', 'User ID not found in token');
        return;
      }
    
      const response = await getChats(userId);
    
      if (response?.data) {
        setChats(response.data);
      } else {
        Alert.alert('Error', 'Invalid response from server');
      }
    };    

    fetchChats();
  }, []);

  const renderItem = ({ item }) => (
    <TouchableOpacity
      style={styles.chatItem}
      onPress={() => navigation.navigate('Chat', { chatId: item.id, chatName: item.title })}
    >
      <View style={styles.avatarContainer}>
        <Image
          source={{ uri: item.avatarUrl || 'https://i.pinimg.com/736x/97/29/79/972979190503cbd54ac183b47872ae80.jpg' }}
          style={styles.avatar}
        />
      </View>
      <View style={styles.chatDetails}>
        <Text style={styles.chatName}>{item.title}</Text>
        <Text style={styles.lastMessage} numberOfLines={1}>
          Participants: {item.participants_ids?.length || 0}
        </Text>
      </View>
      <View style={styles.chatMeta}>
        <Text style={styles.time}>{item.type}</Text>
      </View>
    </TouchableOpacity>
  );

  const openProfile = () => {
    navigation.navigate('Profile');
  };

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <TouchableWithoutFeedback onPress={openProfile}>
          <Ionicons name="person-circle-outline" size={30} color="#1e90ff" />
        </TouchableWithoutFeedback>
        <Text style={styles.headerTitle}>Messages</Text>
        <TouchableWithoutFeedback onPress={() => Alert.alert('Search clicked')}>
          <Ionicons name="search-outline" size={30} color="#1e90ff" />
        </TouchableWithoutFeedback>
      </View>

      <FlatList
  data={chats}
  keyExtractor={(item, index) => item.id?.toString() || index.toString()} // Ensure a valid key
  renderItem={renderItem}
  ListEmptyComponent={() => (
    <View style={{ padding: 20, alignItems: 'center' }}>
      <Text>No chats available</Text>
    </View>
  )}
/>

    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f9f9f9',
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingHorizontal: 16,
    paddingVertical: 10,
    backgroundColor: '#fff',
    elevation: 5,
    borderBottomWidth: 1,
    borderBottomColor: '#eee',
  },
  headerTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    color: '#333',
  },
  chatItem: {
    flexDirection: 'row',
    padding: 16,
    borderBottomWidth: 1,
    borderBottomColor: '#eee',
    backgroundColor: '#fff',
    alignItems: 'center',
  },
  avatarContainer: {
    marginRight: 15,
  },
  avatar: {
    width: 50,
    height: 50,
    borderRadius: 25,
    backgroundColor: '#ccc',
  },
  chatDetails: {
    flex: 1,
  },
  chatName: {
    fontSize: 18,
    fontWeight: '600',
    color: '#333',
  },
  lastMessage: {
    fontSize: 14,
    color: '#777',
    marginTop: 4,
  },
  chatMeta: {
    alignItems: 'flex-end',
  },
  time: {
    fontSize: 12,
    color: '#aaa',
  },
});
